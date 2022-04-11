package vup

import "vup_dd_stats/service/db"

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
