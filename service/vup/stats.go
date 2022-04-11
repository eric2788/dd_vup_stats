package vup

import "vup_dd_stats/service/db"

func GetStats(uid, top int64) (map[string]*Analysis, error) {

	err := db.Database.
		Model(&db.Behaviour{}).
		Where("uid = ? or target_uid = ? and target_uid != uid", uid, uid).
		Select("uid, target_uid, count(*) as count").
		Group("uid, target_uid").
		Order("count desc").
		Error

	if err != nil {
		return nil, err
	}

	//TODO

	return nil, nil
}

func GetTopDDRecords(uid, top int64) ([]db.Behaviour, error) {
	//TODO

	return nil, nil
}

func GetTopSelfRecords(uid, top int64) ([]db.Behaviour, error) {
	//TODO

	return nil, nil
}
