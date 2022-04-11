package vup

import "vup_dd_stats/service/db"

func GetTopDDRecords(uid int64, limit int) ([]db.Behaviour, error) {

	var behaviours []db.Behaviour

	err := db.Database.
		Where("uid = ? and uid != target_uid", uid).
		Order("created_at desc").
		Limit(limit).
		Find(&behaviours).
		Error

	if err != nil {
		return nil, err
	}

	return behaviours, nil
}

func GetTopSelfRecords(uid int64, limit int) ([]db.Behaviour, error) {

	var behaviours []db.Behaviour

	err := db.Database.
		Where("uid = ? and uid = target_uid", uid).
		Order("created_at desc").
		Limit(limit).
		Find(&behaviours).
		Error

	if err != nil {
		return nil, err
	}

	return behaviours, nil
}
