package vup

import (
	"fmt"
	"vup_dd_stats/service/db"
)

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
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.uid, vups.uid").
		Order("spent desc").
		Limit(limit).
		Find(&mostSpentPricedVups).
		Error

	return mostSpentPricedVups, err
}

func GetMostBehaviourVupsByCommand(limit int, command string, price bool) ([]AnalysisUserInfo, error) {
	var mostDDBehaviourVups []AnalysisUserInfo

	orderBy := "count"
	if price {
		orderBy = "price"
	}

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
		Where("behaviours.target_uid != behaviours.uid and behaviours.command = ?", command).
		Group("behaviours.uid, vups.uid").
		Order(fmt.Sprintf("%s desc", orderBy)).
		Limit(limit).
		Find(&mostDDBehaviourVups).
		Error

	return mostDDBehaviourVups, err
}
