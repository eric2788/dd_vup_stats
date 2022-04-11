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
		Where("behaviours.target_uid = ? and behaviours.target_uid != behaviours.uid", uid).
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

func GetGlobalStats() (*GlobalStatistics, error) {

	s, err := statistics.GetListening()
	if err != nil {
		return nil, err
	}

	listeningCount := s.TotalListeningCount

	recordCount, err := GetTotalVupCount()
	if err != nil {
		logger.Errorf("獲取總vup人數出現錯誤: %v", err)
		recordCount = 0
	}

	behaviourCount, err := GetTotalBehaviourCount()
	if err != nil {
		logger.Errorf("獲取總dd行為數時出現錯誤: %v", err)
		behaviourCount = 0
	}

	mostDDBehaviourVups, err := GetMostBehaviourVups(3)
	if err != nil {
		logger.Errorf("獲取dd行為最多的vup時出現錯誤: %v", err)
		mostDDBehaviourVups = []AnalysisUserInfo{}
	}

	mostDDVups, err := GetMostDDVups(3)
	if err != nil {
		logger.Errorf("獲取dd最多的vup時出現錯誤: %v", err)
		mostDDVups = []AnalysisUserInfo{}
	}

	return &GlobalStatistics{
		TotalVupRecorded:      recordCount,
		CurrentListeningCount: listeningCount,
		MostDDBehaviourVupCommands: map[string][]AnalysisUserInfo{
			blive.DanmuMsg:         GetMostBehaviourVupsByCommand(3, blive.DanmuMsg),
			blive.InteractWord:     GetMostBehaviourVupsByCommand(3, blive.InteractWord),
			blive.SuperChatMessage: GetMostBehaviourVupsByCommand(3, blive.SuperChatMessage),
		},
		MostDDBehaviourVups: mostDDBehaviourVups,
		MostDDVups:          mostDDVups,
		TotalDDBehaviours:   behaviourCount,
	}, nil
}
