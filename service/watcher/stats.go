package watcher

import (
	"fmt"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"
)

func GetTotalCommandStats(uid int64) ([]stats.TotalStats, error) {

	var totalStatus []stats.TotalStats

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"command",
			"COUNT(*) as count",
			"SUM(price) as price",
		}).
		Where("uid = ?", uid).
		Group("command").
		Find(&totalStatus).
		Error

	if err != nil {
		logger.Errorf("獲取 %v 的行為時出現錯誤: %v", uid, err)
		return nil, err
	}

	return totalStatus, nil
}

func GetStatsCommand(uid int64, limit int, command string, price bool) (*Analysis, error) {

	var mostDDVups []AnalysisVupInfo

	orderBy := "count"

	if price {
		orderBy = "SUM(price)"
	}

	// D 最多
	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Where("watcher_behaviours.uid = ? AND watcher_behaviours.command = ?", uid, command).
		Joins("left join vups on vups.uid = watcher_behaviours.target_uid").
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.face",
			"COUNT(*) as count",
			"SUM(watcher_behaviours.price) as price",
		}).
		Group("watcher_behaviours.target_uid, vups.uid").
		Order(fmt.Sprintf("%s desc", orderBy)).
		Limit(limit).
		Find(&mostDDVups).
		Error

	if err != nil {
		return nil, err
	}

	// 花費最多
	var mostSpentVups []PricedVupInfo
	err = db.Database.
		Model(&db.WatcherBehaviour{}).
		Joins("left join vups on vups.uid = watcher_behaviours.target_uid").
		Where("watcher_behaviours.uid = ? and watcher_behaviours.command = ?", uid, command).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.face",
			"SUM(watcher_behaviours.price) as spent",
		}).
		Group("watcher_behaviours.target_uid, vups.uid").
		Order("spent desc").
		Limit(limit).
		Find(&mostSpentVups).
		Error

	if err != nil {
		return nil, err
	}

	return &Analysis{
		TopDDVups:    mostDDVups,
		TopSpentVups: mostSpentVups,
	}, nil
}

func GetStats(uid int64, limit int) (*Analysis, error) {

	var mostDDVups []AnalysisVupInfo

	// D 最多
	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Joins("left join vups on vups.uid = watcher_behaviours.target_uid").
		Where("watcher_behaviours.uid = ?", uid).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.face",
			"COUNT(*) as count",
			"SUM(watcher_behaviours.price) as price",
		}).
		Group("watcher_behaviours.target_uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDVups).
		Error

	if err != nil {
		return nil, err
	}

	var mostSpentVups []PricedVupInfo

	// 花費最多
	err = db.Database.
		Model(&db.WatcherBehaviour{}).
		Joins("left join vups on vups.uid = watcher_behaviours.target_uid").
		Where("watcher_behaviours.uid = ?", uid).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.face",
			"SUM(watcher_behaviours.price) as spent",
		}).
		Group("watcher_behaviours.target_uid, vups.uid").
		Order("spent desc").
		Limit(limit).
		Find(&mostSpentVups).
		Error

	if err != nil {
		return nil, err
	}

	return &Analysis{
		TopDDVups:    mostDDVups,
		TopSpentVups: mostSpentVups,
	}, nil
}
