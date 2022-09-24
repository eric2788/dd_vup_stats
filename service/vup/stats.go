package vup

import (
	"fmt"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/statistics"
)

func GetStats(uid int64, limit int) (*Analysis, error) {

	var mostDDVups []AnalysisUserInfo

	// D 最多
	err := db.Database.
		Model(&db.Behaviour{}).
		Joins("left join vups on vups.uid = behaviours.target_uid").
		Where("behaviours.uid = ? and behaviours.target_uid != behaviours.uid", uid).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(*) as count",
			"SUM(behaviours.price) as price",
		}).
		Group("behaviours.target_uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDVups).
		Error

	if err != nil {
		return nil, err
	}

	var mostGuestVups []AnalysisUserInfo

	// 被 D 最多
	err = db.Database.
		Model(&db.Behaviour{}).
		Joins("left join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid = ? and behaviours.target_uid != behaviours.uid", uid).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(*) as count",
			"SUM(behaviours.price) as price",
		}).
		Group("behaviours.uid, vups.uid").
		Limit(limit).
		Order("count desc").
		Find(&mostGuestVups).
		Error

	if err != nil {
		return nil, err
	}

	var mostSpentVups []PricedUserInfo

	// 花費最多
	err = db.Database.
		Model(&db.Behaviour{}).
		Joins("left join vups on vups.uid = behaviours.target_uid").
		Where("behaviours.uid = ? and behaviours.target_uid != behaviours.uid", uid).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"SUM(behaviours.price) as spent",
		}).
		Group("behaviours.target_uid, vups.uid").
		Order("price desc").
		Limit(limit).
		Find(&mostSpentVups).
		Error

	if err != nil {
		return nil, err
	}

	return &Analysis{
		TopDDVups:    mostDDVups,
		TopGuestVups: mostGuestVups,
		TopSpentVups: mostSpentVups,
	}, nil
}

func GetStatsCommand(uid int64, limit int, command string, price bool) (*Analysis, error) {

	var mostDDVups []AnalysisUserInfo

	orderBy := "count"

	if price {
		orderBy = "price"
	}

	// D 最多
	err := db.Database.
		Model(&db.Behaviour{}).
		Joins("left join vups on vups.uid = behaviours.target_uid").
		Where("behaviours.uid = ? and behaviours.target_uid != behaviours.uid and behaviours.command = ?", uid, command).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(*) as count",
			"SUM(behaviours.price) as price",
		}).
		Group("behaviours.target_uid, vups.uid").
		Order(fmt.Sprintf("%s desc", orderBy)).
		Limit(limit).
		Find(&mostDDVups).
		Error

	if err != nil {
		return nil, err
	}

	var mostGuestVups []AnalysisUserInfo

	// 被 D 最多
	err = db.Database.
		Model(&db.Behaviour{}).
		Joins("left join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid = ? and behaviours.target_uid != behaviours.uid and behaviours.command = ?", uid, command).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(*) as count",
			"SUM(behaviours.price) as price",
		}).
		Group("behaviours.uid, vups.uid").
		Limit(limit).
		Order(fmt.Sprintf("%s desc", orderBy)).
		Find(&mostGuestVups).
		Error

	if err != nil {
		return nil, err
	}

	return &Analysis{
		TopDDVups:    mostDDVups,
		TopGuestVups: mostGuestVups,
	}, nil
}

func GetGlobalStats(top int) (*GlobalStatistics, error) {

	s, err := statistics.GetListening()
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

	// get all registered commands
	registeredCommands := blive.GetRegisteredCommands()
	mostDDBehaviourVupCommands := make(map[string][]AnalysisUserInfo, len(registeredCommands))
	for _, command := range registeredCommands {
		mostDDBehaviourVupCommands[command] = GetMostBehaviourVupsByCommand(top, command)
	}

	return &GlobalStatistics{
		TotalVupRecorded:           recordCount,
		CurrentListeningCount:      listeningCount,
		MostDDBehaviourVupCommands: mostDDBehaviourVupCommands,
		MostDDBehaviourVups:        mostDDBehaviourVups,
		MostDDVups:                 mostDDVups,
		MostSpentVups:              mostSpentVups,
		TotalDDBehaviours:          behaviourCount,
	}, nil
}

func GetTotalStatusByCommand(uid int64, command string) TotalStats {

	var totalStatus TotalStats

	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"COUNT(*) as count",
			"SUM(price) as price",
		}).
		Where("uid = ? and command = ? and uid != target_uid", uid, command).
		Find(&totalStatus).
		Error

	if err != nil {
		logger.Errorf("獲取 %v 的 %v 行為時出現錯誤: %v", uid, command, err)
		return TotalStats{-1, -1}
	}

	return totalStatus

}
