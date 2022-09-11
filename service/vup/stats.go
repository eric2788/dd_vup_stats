package vup

import (
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
		Distinct([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
		}).
		Select("COUNT(*) as count").
		Group("behaviours.target_uid").
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
		Distinct([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
		}).
		Select("COUNT(*) as count").
		Group("behaviours.uid").
		Limit(limit).
		Order("count desc").
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

func GetStatsCommand(uid int64, limit int, command string) (*Analysis, error) {

	var mostDDVups []AnalysisUserInfo

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
		}).
		Group("behaviours.target_uid").
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
		Where("behaviours.target_uid = ? and behaviours.target_uid != behaviours.uid and behaviours.command = ?", uid, command).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(*) as count",
		}).
		Group("behaviours.uid").
		Limit(limit).
		Order("count desc").
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

	return &GlobalStatistics{
		TotalVupRecorded:      recordCount,
		CurrentListeningCount: listeningCount,
		MostDDBehaviourVupCommands: map[string][]AnalysisUserInfo{
			blive.DanmuMsg:         GetMostBehaviourVupsByCommand(top, blive.DanmuMsg),
			blive.InteractWord:     GetMostBehaviourVupsByCommand(top, blive.InteractWord),
			blive.SuperChatMessage: GetMostBehaviourVupsByCommand(top, blive.SuperChatMessage),
		},
		MostDDBehaviourVups: mostDDBehaviourVups,
		MostDDVups:          mostDDVups,
		TotalDDBehaviours:   behaviourCount,
	}, nil
}

func GetTotalCountByCommand(uid int64, command string) int64 {

	var totalCount int64
	err := db.Database.
		Model(&db.Behaviour{}).
		Where("uid = ? and command = ? and uid != target_uid", uid, command).
		Count(&totalCount).
		Error

	if err != nil {
		logger.Errorf("獲取 %v 的 %v 行為時出現錯誤: %v", uid, command, err)
		totalCount = -1
	}

	return totalCount

}
