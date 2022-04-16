package vup

import (
	"fmt"
	"math"
	"vup_dd_stats/service/db"
)

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

func GetTopGuestRecords(uid int64, limit int) ([]db.Behaviour, error) {

	var behaviours []db.Behaviour

	err := db.Database.
		Where("target_uid = ? and uid != target_uid", uid).
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

func GetGlobalRecords(search string, page, pageSize int, showSelf bool) (*ListResp[db.Behaviour], error) {

	// ensure page is valid
	page = int(math.Max(1, float64(page)))

	//ensure pageSize is valid
	pageSize = int(math.Max(1, float64(pageSize)))

	var behaviours []db.Behaviour

	r := db.Database.Order("created_at desc")

	if showSelf {
		r = r.Where("display like ?", fmt.Sprintf("%%%s%%", search))
	} else {
		r = r.Where("display like ? and uid != target_uid", fmt.Sprintf("%%%s%%", search))
	}

	err := r.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&behaviours).
		Error

	if err != nil {
		return nil, err
	}

	var totalSearchCount int64

	err = db.Database.
		Model(&db.Behaviour{}).
		Where("display like ?", fmt.Sprintf("%%%s%%", search)).
		Count(&totalSearchCount).
		Error

	if err != nil {
		return nil, err
	}

	return &ListResp[db.Behaviour]{
		Total:   totalSearchCount,
		List:    behaviours,
		Page:    page,
		Size:    pageSize,
		MaxPage: int64(math.Ceil(float64(totalSearchCount) / float64(pageSize))),
	}, nil
}
