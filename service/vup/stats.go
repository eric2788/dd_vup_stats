package vup

import (
	"fmt"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"

	"gorm.io/gorm"
)


// GetStats get user stats
// TODO: make concurrent
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
		Where("behaviours.uid = ? and behaviours.target_uid != behaviours.uid and behaviours.price > 0", uid).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"SUM(behaviours.price) as spent",
		}).
		Group("behaviours.target_uid, vups.uid").
		Order("spent desc").
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

// GetStatsCommand get user stats by command
// TODO: make concurrent
func GetStatsCommand(uid int64, limit int, command string, price bool) (*Analysis, error) {

	var mostDDVups []AnalysisUserInfo

	r := db.Database.Model(&db.Behaviour{})

	orderBy := "count"

	if price {
		orderBy = "SUM(price)"
		r = r.Where("behaviours.price > 0")
	}

	r = r.Session(&gorm.Session{})

	// D 最多
	err := r.
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
	err = r.
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

// GetTotalBehaviourCount get total behaviour count by command
//
// Deprecated: use GetTotalCommandStats instead
func GetTotalStatusByCommand(uid int64, command string) stats.TotalStats {

	var totalStatus stats.TotalStats

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
		return stats.TotalStats{Command: command, Count: -1, Price: -1}
	}

	return totalStatus
}

func GetTotalCommandStats(uid int64) ([]stats.TotalStats, error) {
	var stats []stats.TotalStats

	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"command",
			"COUNT(*) as count",
			"SUM(price) as price",
		}).
		Where("uid = ? and uid != target_uid", uid).
		Group("command").
		Find(&stats).
		Error

	if err != nil {
		logger.Errorf("獲取 %v 的行為時出現錯誤: %v", uid, err)
		return nil, err
	}

	return stats, nil
}
