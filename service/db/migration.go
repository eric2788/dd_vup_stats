package db

import "gorm.io/gorm"

// MigrateToVup will migrate specific watcher to vup
// it will first move the records of the watcher from watcher_behaviour to behaviour table
// then delete the records of that watcher in watcher_behaviour table
func MigrateToVup(uid int64) {

	var watcherBehaviours []WatcherBehaviour

	err := Database.
		Model(&WatcherBehaviour{}).
		Where("uid = ?", uid).
		Find(&watcherBehaviours).
		Error

	if err != nil {
		log.Errorf("從 watcher_behaviour 表中獲取 UID:%v 的行為記錄時出現錯誤: %v", uid, err)
		return
	}

	// if there is no record in watcher_behaviour table, then return
	if len(watcherBehaviours) == 0 {
		return
	}

	// insert records into behaviour table
	behaviours := make([]*Behaviour, len(watcherBehaviours))

	for i, b := range watcherBehaviours {
		behaviours[i] = b.ToBehaviour()
	}

	log.Debugf("即將遷移 %d 筆記錄到 behaviour 表", len(behaviours))

	result := Database.
		Model(&Behaviour{}).
		CreateInBatches(behaviours, len(behaviours))

	if result.Error != nil {
		log.Errorf("將 UID:%v 的行為記錄從 watcher_behaviour 表中移動到 behaviour 表時出現錯誤: %v", uid, err)
		return
	} else if result.RowsAffected > 0 {
		log.Infof("成功將 UID:%v 的 %d 筆行為記錄從 watcher_behaviour 表中移動到 behaviour 表", uid, result.RowsAffected)
	}

	log.Debugf("即將刪除 UID:%v 的 %d 筆記錄", uid, len(watcherBehaviours))

	// delete records from watcher_behaviour table
	result = Database.
		Where("uid = ?", uid).
		Delete(&WatcherBehaviour{})

	if result.Error != nil {
		log.Errorf("刪除 UID:%v 的行為記錄時出現錯誤: %v", uid, err)
	} else if result.RowsAffected > 0 {
		log.Infof("成功從 watcher_behaviour 刪除 UID:%v 的 %d 筆行為記錄", uid, result.RowsAffected)
	}

}

func createMaterializedViews(db *gorm.DB) {

	// create materialized views: watchers and watchers_stats
	go db.Transaction(func(tx *gorm.DB) error {

		// create materialized views
		err := db.Exec(`
			create materialized view if not exists watchers as 
			SELECT 
				watcher_behaviours.uid,
				(array_agg(watcher_behaviours.u_name ORDER BY watcher_behaviours.created_at DESC))[1] AS u_name
			FROM watcher_behaviours
			GROUP BY watcher_behaviours.uid;
		`).Error

		if err != nil {
			log.Errorf("创建 Materialized View `watchers` 失败: %v", err)
			return err
		}

		err = db.Exec(`
			create materialized view if not exists watchers_stats as
			SELECT 
				v.uid,
				v.u_name,
				count(*) AS count,
				sum(b.price) AS spent,
				count(DISTINCT b.target_uid) AS dd
			FROM (watchers v
			JOIN watcher_behaviours b ON ((b.uid = v.uid)))
			GROUP BY b.uid, v.uid, v.u_name;
		`).Error

		if err != nil {
			log.Errorf("创建 Materialized View `watchers_stats` 失败: %v", err)
			return err
		}

		// create indexes for materialized view

		err = db.Exec(`CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS watchers_uid_idx ON public.watchers USING btree (uid)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `watchers` 的索引 watchers_uid_idx 失败: %v", err)
		}

		err = db.Exec(`CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS watchers_stats_uid_idx ON public.watchers_stats USING btree (uid)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `watchers_stats` 的索引 watchers_stats_uid_idx 失败: %v", err)
		}

		return nil
	})

	// create materialized views: vups_with_watcher_behaviours
	go db.Transaction(func(tx *gorm.DB) error {

		err := db.Exec(`
			create materialized view if not exists vups_with_watcher_behaviours as
			SELECT 
				vups.name,
				vups.uid,
				count(DISTINCT watcher_behaviours.uid) AS famous,
				count(*) AS interacted,
				sum(watcher_behaviours.price) AS earned
			FROM (watcher_behaviours
			JOIN vups ON ((vups.uid = watcher_behaviours.target_uid)))
			GROUP BY watcher_behaviours.target_uid, vups.uid;
		`).Error

		if err != nil {
			log.Errorf("创建 Materialized View `vups_with_watcher_behaviours` 失败: %v", err)
		}

		err = db.Exec(`CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS vups_with_watcher_behaviours_uid_idx ON public.vups_with_watcher_behaviours USING btree (uid)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `vups_with_watcher_behaviours` 的索引 vups_with_watcher_behaviours_uid_idx 失败: %v", err)
			return err
		}

		return nil
	})

	// create materialized views: watchers_fans
	go db.Transaction(func(tx *gorm.DB) error {
		err := db.Exec(`
			create materialized view if not exists watchers_fans as 
			SELECT 
				concat(watcher_behaviours.uid, '-', watcher_behaviours.target_uid, '-', watcher_behaviours.command) AS id,
				watcher_behaviours.uid,
				(array_agg(watcher_behaviours.u_name ORDER BY watcher_behaviours.created_at DESC))[1] AS u_name,
				watcher_behaviours.target_uid,
				watcher_behaviours.command,
				count(*) AS count,
				sum(watcher_behaviours.price) AS price
			FROM watcher_behaviours
			GROUP BY watcher_behaviours.uid, watcher_behaviours.target_uid, watcher_behaviours.command;
		`).Error

		if err != nil {
			log.Errorf("创建 Materialized View `watchers_fans` 失败: %v", err)
			return err
		}

		err = db.Exec(`CREATE INDEX CONCURRENTLY IF NOT EXISTS watchers_fans_uid_target_uid_command_idx ON public.watchers_fans USING btree (uid, target_uid, command)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `watchers_fans` 的索引 watchers_fans_uid_target_uid_command_idx 失败: %v", err)
			return err
		}

		err = db.Exec(`CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS watchers_fans_id_idx ON public.watchers_fans USING btree (id)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `watchers_fans` 的索引 watchers_fans_id_idx 失败: %v", err)
			return err
		}

		err = db.Exec(`CREATE INDEX CONCURRENTLY IF NOT EXISTS watchers_fans_target_uid_idx ON public.watchers_fans USING btree (target_uid)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `watchers_fans` 的索引 watchers_fans_target_uid_idx 失败: %v", err)
			return err
		}

		err = db.Exec(`CREATE INDEX CONCURRENTLY IF NOT EXISTS watchers_fans_command_idx ON public.watchers_fans USING btree (command)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `watchers_fans` 的索引 watchers_fans_command_idx 失败: %v", err)
			return err
		}

		err = db.Exec(`CREATE INDEX CONCURRENTLY IF NOT EXISTS watchers_fans_uid_target_uid_idx ON public.watchers_fans USING btree (uid, target_uid)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `watchers_fans` 的索引 watchers_fans_uid_target_uid_idx 失败: %v", err)
			return err
		}

		err = db.Exec(`CREATE INDEX CONCURRENTLY IF NOT EXISTS watchers_fans_uid_u_name_idx ON public.watchers_fans USING btree (uid) INCLUDE (u_name)`).Error
		if err != nil {
			log.Errorf("创建 Materialized View `watchers_fans` 的索引 watchers_fans_uid_u_name_idx 失败: %v", err)
			return err
		}

		return nil
	})

}
