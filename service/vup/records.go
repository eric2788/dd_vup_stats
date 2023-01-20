package vup

import (
	"errors"
	"fmt"
	"math"
	"vup_dd_stats/service/db"

	"gorm.io/gorm"
)

func GetTopDDRecords(uid int64, limit int) ([]db.Behaviour, error) {

	var behaviours []db.Behaviour

	err := db.Database.
		Where("uid = ? and uid != target_uid", uid).
		Order("created_at desc").
		Limit(limit).
		Find(&behaviours).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
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

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
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

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return behaviours, nil
}

func GetGlobalRecords(search, cmd string, page, pageSize int, showSelf bool) (*ListResp[RecordResp], error) {

	// ensure page is valid
	page = int(math.Max(1, float64(page)))

	//ensure pageSize is valid
	pageSize = int(math.Max(1, float64(pageSize)))

	var records []RecordResp

	r := db.Database.Model(&db.Behaviour{})

	r = buildWhereStatement(r, search, cmd, showSelf)

	err := r.
		Select("behaviours.*, vups.face as vup_face").
		Joins("left join vups on vups.uid = behaviours.uid").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("behaviours.created_at desc").
		Find(&records).
		Error

	if err != nil {
		return nil, err
	}

	var totalSearchCount int64

	r = db.Database.Model(&db.Behaviour{})

	err = buildWhereStatement(r, search, cmd, showSelf).Count(&totalSearchCount).Error

	if err != nil {
		return nil, err
	}

	go logger.Info("记录行为搜索: %s", search)

	return &ListResp[RecordResp]{
		Total:   totalSearchCount,
		List:    records,
		Page:    page,
		Size:    pageSize,
		MaxPage: int64(math.Ceil(float64(totalSearchCount) / float64(pageSize))),
	}, nil
}

func buildWhereStatement(r *gorm.DB, search, cmd string, showSelf bool) *gorm.DB {

	r = r.Where("behaviours.display like ?", fmt.Sprintf("%%%s%%", search))

	if !showSelf {
		r = r.Where("behaviours.uid != behaviours.target_uid")
	}

	if cmd != "" {
		r = r.Where("behaviours.command = ?", cmd)
	}

	return r
}
