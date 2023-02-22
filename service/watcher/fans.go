package watcher

import (
	"fmt"
	"vup_dd_stats/service/db"
)

func GetFanStatsForVup(uid int64, limit int, t string) ([]AnalysisWatcherInfo, error) {
	switch t {
	case "count":
		return GetMostBehavioursByVup(uid, limit)
	case "spent":
		return GetMostSpentByVup(uid, limit)
	default:
		return nil, fmt.Errorf("不支持的类型: %v", t)
	}
}

// GetMostBehavioursByVup 返回该 vup 中最高互动的 watcher
func GetMostBehavioursByVup(uid int64, limit int) ([]AnalysisWatcherInfo, error) {

	var mostDDWatchers []AnalysisWatcherInfo

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			"(array_agg(u_name order by created_at desc))[1] as u_name",
			"COUNT(*) as count",
		}).
		Where("target_uid = ?", uid).
		Group("uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDWatchers).
		Error

	if err != nil {
		return nil, err
	}

	return mostDDWatchers, nil
}

// GetMostSpentByVup 返回该 vup 中花费最多的 watcher
func GetMostSpentByVup(uid int64, limit int) ([]AnalysisWatcherInfo, error) {

	var mostSpentWatchers []AnalysisWatcherInfo

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			"(array_agg(u_name order by created_at desc))[1] as u_name",
			"SUM(price) as price",
		}).
		Where("target_uid = ? and price > 0", uid).
		Group("uid").
		Order("price desc").
		Limit(limit).
		Find(&mostSpentWatchers).
		Error

	if err != nil {
		return nil, err
	}

	return mostSpentWatchers, nil
}
