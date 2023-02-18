package watcher

import (
	"fmt"
	"vup_dd_stats/service/db"
)

func GetStatsByType(top int, t string) (interface{}, error) {
	switch t {
	case "count":
		return GetTotalCount()
	case "dd":
		return GetMostDDWatchers(top)
	case "behaviours":
		return GetMostBehaviourWatchers(top)
	case "spent":
		return GetMostSpentWatchers(top)
	case "famous":
		return GetMostFamousVups(top)
	case "interacted":
		return GetMostInteractedVups(top)
	default:
		return nil, fmt.Errorf("unknown type: %s", t)
	}
}

// GetMostDDVups 獲取最受普通B站用户欢迎的vups (最多人访问的vup)
func GetMostFamousVups(limit int) ([]AnalysisVupInfo, error) {
	var mostFamousVups []AnalysisVupInfo
	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"vups.uid",
			"vups.name",
			"vups.face",
			"COUNT(DISTINCT watcher_behaviours.uid) as count",
		}).
		Joins("inner join vups on vups.uid = watcher_behaviours.target_uid").
		Group("watcher_behaviours.target_uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostFamousVups).
		Error
	return mostFamousVups, err
}

// GetMostInteractedVups 获取经常被普通B站用户互动的vups (被普通B站用户互动次数最多)
func GetMostInteractedVups(limit int) ([]AnalysisVupInfo, error) {
	var mostInteractedVups []AnalysisVupInfo
	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"vups.uid",
			"vups.name",
			"vups.face",
			"COUNT(watcher_behaviours.uid) as count",
		}).
		Joins("inner join vups on vups.uid = watcher_behaviours.target_uid").
		Group("watcher_behaviours.target_uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostInteractedVups).
		Error
	return mostInteractedVups, err
}

// GetMostDDWatchers 獲取進入最多不同直播間的 dd
func GetMostDDWatchers(limit int) ([]AnalysisWatcherInfo, error) {

	var mostDDWatchers []AnalysisWatcherInfo

	u_name := "u_name"
	if db.DatabaseType == "postgres" {
		u_name = "(array_agg(u_name order by created_at desc))[1] as u_name"
	}

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			u_name,
			"COUNT(DISTINCT target_uid) as count",
		}).
		Group("uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDWatchers).
		Error

	return mostDDWatchers, err
}

func GetTotalCount() (int64, error) {
	var count int64
	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Count(&count).
		Error
	return count, err
}

// GetMostBehaviourWatchers 獲取最多行為的 dd
func GetMostBehaviourWatchers(limit int) ([]AnalysisWatcherInfo, error) {

	var mostBehaviourWatchers []AnalysisWatcherInfo

	u_name := "u_name"
	if db.DatabaseType == "postgres" {
		u_name = "(array_agg(u_name order by created_at desc))[1] as u_name"
	}

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			u_name,
			"COUNT(*) as count",
		}).
		Group("uid").
		Order("count desc").
		Limit(limit).
		Find(&mostBehaviourWatchers).
		Error

	return mostBehaviourWatchers, err
}

// GetMostSpentWatchers 獲取花費最多的 dd
func GetMostSpentWatchers(limit int) ([]AnalysisWatcherInfo, error) {

	var mostSpentWatchers []AnalysisWatcherInfo

	u_name := "u_name"
	if db.DatabaseType == "postgres" {
		u_name = "(array_agg(u_name order by created_at desc))[1] as u_name"
	}

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			u_name,
			"SUM(price) as price",
		}).
		Where("price > 0").
		Group("uid").
		Order("price desc").
		Limit(limit).
		Find(&mostSpentWatchers).
		Error

	return mostSpentWatchers, err
}

func GetMostBehaviourWatchersByCommand(limit int, command string, price bool) ([]AnalysisWatcherInfo, error) {
	var mostDDBehaviourVups []AnalysisWatcherInfo

	orderBy := "count"
	if price {
		orderBy = "price"
	}

	uName := "u_name"
	if db.DatabaseType == "postgres" {
		uName = "(array_agg(u_name order by created_at desc))[1] as u_name"
	}

	err := db.Database.
		Model(&db.WatcherBehaviour{}).
		Select([]string{
			"uid",
			uName,
			"COUNT(*) as count",
			"SUM(price) as price",
		}).
		Where("command = ?", command).
		Group("uid").
		Order(fmt.Sprintf("%s desc", orderBy)).
		Limit(limit).
		Find(&mostDDBehaviourVups).
		Error

	return mostDDBehaviourVups, err
}
