package vup

import (
	"fmt"
	"time"
	"vup_dd_stats/service/db"
)

func IsVup(uid int64) (bool, error) {
	var exist bool

	err := db.Database.
		Model(&db.Vup{}).
		Where("uid = ?", uid).
		Select("count(*) > 0").
		Find(&exist).
		Error

	if err != nil {
		return false, err
	}

	return exist, nil
}

func GetTotalVupCount() (int64, error) {
	var count int64
	err := db.Database.Model(&db.Vup{}).Count(&count).Error
	return count, err
}

func GetVups(page, size int, desc bool) (*ListResp, error) {

	total, err := GetTotalVupCount()

	if err != nil {
		return nil, err
	}

	var infos []UserInfo

	order := "desc"
	if !desc {
		order = "asc"
	}

	err = db.Database.
		Model(&db.Vup{}).
		Limit(size).
		Offset((page - 1) * size).
		Order(fmt.Sprintf("first_listen_at %v", order)).
		Select([]string{"uid", "name", "face", "first_listen_at", "room_id", "sign", "count(behaviour.uid) AS DD_Count"}).
		Joins("left join behaviour on behaviour.uid = vup.uid").
		Find(&infos).
		Error

	if err != nil {
		return nil, err
	}

	user := make([]*UserResp, len(infos))

	for i, info := range infos {
		user[i] = &UserResp{
			UserInfo:        info,
			Listening:       false,
			LastListenedAt:  time.Time{},
			LastBehaviourAt: time.Time{},
		}
	}

	return &ListResp{
		Page:    page,
		Size:    size,
		MaxPage: total/int64(size) + 1,
		List:    user,
	}, nil
}
