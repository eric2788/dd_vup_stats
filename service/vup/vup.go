package vup

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gorm.io/gorm/clause"
	"math"
	"time"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/statistics"
)

var logger = logrus.WithField("service", "vup")

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

func GetVup(uid int64) (*UserResp, error) {

	resp, err := statistics.GetListening()

	if err != nil {
		return nil, err
	}

	var vup UserInfo
	err = db.Database.
		Model(&db.Vup{}).
		Where("vups.uid = ?", uid).
		Select([]string{
			"vups.uid",
			"vups.name",
			"vups.face",
			"vups.first_listen_at",
			"vups.room_id",
			"vups.sign",
			"COUNT(behaviours.uid) AS dd_count",
			"MAX(behaviours.created_at) AS last_behaviour_at",
		}).
		Joins("left join behaviours on behaviours.uid = vups.uid and behaviours.uid != behaviours.target_uid").
		Where("vups.uid = ?", uid).
		Find(&vup).
		Error

	if err != nil {
		return nil, err
	}

	listening := slices.Contains(resp.Rooms, vup.RoomId)
	lastListenAt := GetLastListen(&vup, listening)

	return &UserResp{
		UserInfo:       vup,
		Listening:      listening,
		LastListenedAt: lastListenAt,
	}, nil

}

func GetLastListen(vup *UserInfo, listening bool) time.Time {

	lastListenAt := time.Now()

	if !listening {

		var lastListen db.LastListen

		err := db.Database.
			Where("uid = ?", vup.Uid).
			FirstOrCreate(&lastListen, db.LastListen{
				Uid:          vup.Uid,
				LastListenAt: time.Now(),
			}).Error

		if err != nil {
			logger.Errorf("嘗試插入最後監聽訊息時出現錯誤: %v", err)
			lastListenAt = vup.FirstListenAt
		} else {
			lastListenAt = lastListen.LastListenAt
		}

	} else {

		err := db.Database.
			Clauses(clause.OnConflict{DoNothing: true}).
			Delete(&db.LastListen{}, vup.Uid).Error

		if err != nil {
			logger.Errorf("嘗試刪除最後監聽訊息時出現錯誤: %v", err)
		}

	}

	return lastListenAt
}

func GetVups(page, size int, desc bool, orderBy string) (*ListResp, error) {

	total, err := GetTotalVupCount()

	if err != nil {
		return nil, err
	}

	var infos []UserInfo

	order := "desc"
	if !desc {
		order = "asc"
	}

	resp, err := statistics.GetListening()

	if err != nil {
		return nil, err
	}

	err = db.Database.
		Model(&db.Vup{}).
		Select([]string{
			"vups.uid",
			"vups.name",
			"vups.face",
			"vups.first_listen_at",
			"vups.room_id",
			"vups.sign",
			"COUNT(behaviours.uid) AS dd_count",
			"MAX(behaviours.created_at) AS last_behaviour_at",
		}).
		Limit(size).
		Offset((page - 1) * size).
		Order(fmt.Sprintf("%v %v", orderBy, order)).
		Joins("left join behaviours on behaviours.uid = vups.uid and behaviours.uid != behaviours.target_uid").
		Group("uid").
		Find(&infos).
		Error

	if err != nil {
		return nil, err
	}

	user := make([]*UserResp, len(infos))

	for i, info := range infos {

		listening := slices.Contains(resp.Rooms, info.RoomId)

		lastListenAt := GetLastListen(&info, listening)

		user[i] = &UserResp{
			UserInfo:       info,
			Listening:      listening,
			LastListenedAt: lastListenAt,
		}
	}

	return &ListResp{
		Page:    page,
		Size:    size,
		MaxPage: int64(math.Ceil(float64(total)/float64(size))) + 1,
		Total:   total,
		List:    user,
	}, nil
}

func GetMostDDVups(limit int) ([]AnalysisUserInfo, error) {

	var mostDDVups []AnalysisUserInfo

	err := db.Database.
		Model(&db.Behaviour{}).
		Select([]string{
			"vups.name",
			"vups.uid",
			"vups.room_id",
			"vups.face",
			"vups.sign",
			"COUNT(DISTINCT behaviours.target_uid) as count",
		}).
		Joins("left join vups on vups.uid = behaviours.uid").
		Where("behaviours.target_uid != behaviours.uid").
		Group("behaviours.uid").
		Order("count desc").
		Limit(3).
		Find(&mostDDVups).
		Error

	return mostDDVups, err
}
