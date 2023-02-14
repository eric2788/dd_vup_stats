package watcher

import (
	"errors"
	"gorm.io/gorm"
	"vup_dd_stats/service/db"

	"github.com/sirupsen/logrus"
)

var logger = logrus.WithField("service", "watcher")

func GetWatcher(uid int64) (*WatcherResp, error) {
	var resp WatcherResp

	uName := "u_name"
	if db.DatabaseType == "postgres" {
		uName = "(array_agg(u_name order by created_at desc))[1] as u_name"
	}

	uNames := "GROUP_CONCAT(DISTINCT u_name SEPARATOR ',') as u_names"
	if db.DatabaseType == "postgres" {
		uNames = "array_to_string(array_agg(distinct u_name), ',') as u_names"
	}

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			uName,
			uNames,
			"COUNT(target_uid) as dd_count",
			"MAX(created_at) AS last_behaviour_at",
			"MIN(created_at) AS first_listen_at",
			"SUM(price) as total_spent",
		}).
		Where("uid = ?", uid).
		Group("uid").
		Take(&resp).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) || resp.Uid == 0 {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	resp.Behaviours, err = GetTotalCommandStats(uid)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
