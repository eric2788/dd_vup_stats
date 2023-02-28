package crontab

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var logger = logrus.WithField("service", "crontab")

func StartRefreshMViewJob(db *gorm.DB) {
	c := cron.New()

	c.AddFunc("0 1 * * *", func() {
		logger.Infof("即将刷新 Materialized View `watchers`...")
		err := db.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY watchers`).Error
		if err != nil {
			logger.Errorf("刷新 Materialized View `watchers` 失败: %v", err)
		}
	})

	c.AddFunc("0 2 * * *", func() {
		logger.Infof("即将刷新 Materialized View `watchers_stats`...")
		err := db.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY watchers_stats`).Error
		if err != nil {
			logger.Errorf("刷新 Materialized View `watchers_stats` 失败: %v", err)
		}
	})

	c.AddFunc("0 3 * * *", func() {
		logger.Infof("即将刷新 Materialized View `vups_with_watcher_behaviours`...")
		err := db.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY vups_with_watcher_behaviours`).Error
		if err != nil {
			logger.Errorf("刷新 Materialized View `vups_with_watcher_behaviours` 失败: %v", err)
		}
	})

	c.AddFunc("0 4 * * *", func() {
		logger.Infof("即将刷新 Materialized View `watchers_fans`...")
		err := db.Exec(`REFRESH MATERIALIZED VIEW CONCURRENTLY watchers_fans`).Error
		if err != nil {
			logger.Errorf("刷新 Materialized View `watchers_fans` 失败: %v", err)
		}
	})

	c.AddFunc("30 5 * * 1", func() {
		logger.Infof("即将分析所有 tables...")
		err := db.Exec(`ANALYZE`).Error
		if err != nil {
			logger.Errorf("分析所有 tables 失败: %v", err)
		}
	})

	c.Start()
}
