package vup

import (
	"fmt"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"
)

func GetCountStats() (*map[string]int64, error) {
	s, err := stats.GetListening()
	if err != nil {
		logger.Errorf("獲取總聆聽人數出現錯誤: %v", err)
		return nil, err
	}

	recordCount, err := GetTotalVupCount()
	if err != nil {
		logger.Errorf("獲取總vup人數出現錯誤: %v", err)
		return nil, err
	}

	behaviourCount, err := GetTotalBehaviourCount()
	if err != nil {
		logger.Errorf("獲取總dd行為數時出現錯誤: %v", err)
		return nil, err
	}

	return &map[string]int64{
		"total_vup_recorded":      recordCount,
		"current_listening_count": s.TotalListeningCount,
		"total_dd_behaviours":     behaviourCount,
	}, nil
}

func GetStatsByType(top int, t string) (interface{}, error) {
	switch t {
	case "count":
		return GetCountStats()
	case "dd":
		return GetMostDDVups(top)
	case "behaviours":
		return GetMostBehaviourVups(top)
	case "spent":
		return GetMostSpentPricedVups(top)
	case "famous":
		return GetMostFamousVups(top)
	case "interacted":
		return GetMostInteractedVups(top)
	case "earned":
		return GetMostEarnedVups(top)
	default:
		return nil, fmt.Errorf("不支持的类型: %q", t)
	}
}

// GetMostDDVups 獲取進入最多不同直播間的 vups
func GetMostDDVups(limit int) ([]AnalysisUserInfo, error) {

	var mostDDVups []AnalysisUserInfo

	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(DISTINCT behaviours.target_uid) as count",
		}).
		Joins("inner join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDVups).
		Error

	return mostDDVups, err
}

// GetMostEarnedVups 获取营收最多的vups (被最多的vup打赏)
func GetMostEarnedVups(limit int) ([]AnalysisUserInfo, error) {
	var mostEarnedVups []AnalysisUserInfo
	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"SUM(behaviours.price) as price",
		}).
		Joins("inner join vups on vups.uid = behaviours.target_uid").
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.target_uid, vups.uid").
		Order("price desc").
		Limit(limit).
		Find(&mostEarnedVups).
		Error

	return mostEarnedVups, err

}

// GetMostFamousVups 获取最受欢迎的vups (被最多的vup访问过)
func GetMostFamousVups(limit int) ([]AnalysisUserInfo, error) {
	var mostFamousVups []AnalysisUserInfo
	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(DISTINCT behaviours.uid) as count",
		}).
		Joins("inner join vups on vups.uid = behaviours.target_uid").
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.target_uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostFamousVups).
		Error
	return mostFamousVups, err
}

// GetMostInteractedVups 获取经常被互动的vups (被互动次数最多)
func GetMostInteractedVups(limit int) ([]AnalysisUserInfo, error) {
	var mostInteractedVups []AnalysisUserInfo
	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(*) as count",
		}).
		Joins("inner join vups on vups.uid = behaviours.target_uid").
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.target_uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostInteractedVups).
		Error
	return mostInteractedVups, err
}

func GetMostBehaviourVups(limit int) ([]AnalysisUserInfo, error) {
	var mostDDBehaviourVups []AnalysisUserInfo

	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(*) as count",
			"SUM(behaviours.price) as price",
		}).
		Joins("inner join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDBehaviourVups).
		Error

	return mostDDBehaviourVups, err
}

func GetMostSpentPricedVups(limit int) ([]PricedUserInfo, error) {
	var mostSpentPricedVups []PricedUserInfo

	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"SUM(behaviours.price) as spent",
		}).
		Joins("inner join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid != behaviours.uid and behaviours.price > 0").
		Group("behaviours.uid, vups.uid").
		Order("spent desc").
		Limit(limit).
		Find(&mostSpentPricedVups).
		Error

	return mostSpentPricedVups, err
}
