package vup

import (
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"
)

func GetCountStats() (map[string]int64, error) {
	s, err := stats.GetListening()
	if err == nil {
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

	return map[string]int64{
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
	default:
		return GetGlobalStats(top)
	}
}

// GetGlobalStats get global stats with all information
//
// Deprecated: high performance cost and slow response, recommend use GetStatsByType instead
func GetGlobalStats(top int) (*Globalstats, error) {
	s, err := stats.GetListening()
	if err != nil {
		return nil, err
	}

	listeningCount := s.TotalListeningCount

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

	mostDDBehaviourVups, err := GetMostBehaviourVups(top)
	if err != nil {
		logger.Errorf("獲取dd行為最多的vup時出現錯誤: %v", err)
		return nil, err
	}

	mostDDVups, err := GetMostDDVups(top)
	if err != nil {
		logger.Errorf("獲取dd最多的vup時出現錯誤: %v", err)
		return nil, err
	}

	mostSpentVups, err := GetMostSpentPricedVups(top)
	if err != nil {
		logger.Errorf("獲取花費最多的vup時出現錯誤: %v", err)
		return nil, err
	}

	return &Globalstats{
		TotalVupRecorded:      recordCount,
		CurrentListeningCount: listeningCount,
		MostDDBehaviourVups:   mostDDBehaviourVups,
		MostDDVups:            mostDDVups,
		MostSpentVups:         mostSpentVups,
		TotalDDBehaviours:     behaviourCount,
	}, nil
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
		Joins("left join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDVups).
		Error

	return mostDDVups, err
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
		Joins("left join vups on vups.uid = behaviours.uid").
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
		Joins("left join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid != behaviours.uid and behaviours.price > 0").
		Group("behaviours.uid, vups.uid").
		Order("spent desc").
		Limit(limit).
		Find(&mostSpentPricedVups).
		Error

	return mostSpentPricedVups, err
}