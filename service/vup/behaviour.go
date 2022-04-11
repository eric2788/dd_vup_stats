package vup

import "vup_dd_stats/service/db"

func GetTotalBehaviourCount() (int64, error) {

	var recordCount int64

	err := db.Database.
		Model(&db.Behaviour{}).
		Count(&recordCount).
		Error

	return recordCount, err
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
		}).
		Joins("left join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDBehaviourVups).
		Error

	return mostDDBehaviourVups, err
}

func GetMostBehaviourVupsByCommand(limit int, command string) []AnalysisUserInfo {
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
		}).
		Joins("left join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid != behaviours.uid and behaviours.command = ?", command).
		Group("behaviours.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDBehaviourVups).
		Error

	if err != nil {
		logger.Errorf("獲取在 %v 的DD行為數量最多的vup時出現錯誤: %v", command, err)
		mostDDBehaviourVups = make([]AnalysisUserInfo, 0)
	}

	return mostDDBehaviourVups
}
