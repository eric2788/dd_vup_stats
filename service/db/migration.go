package db

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
