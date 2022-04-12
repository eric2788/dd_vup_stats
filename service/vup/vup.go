package vup

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gorm.io/gorm/clause"
	"math"
	"time"
	"vup_dd_stats/service/blive"
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

func GetVup(uid int64) (*UserDetailResp, error) {

	listeningRooms := *(statistics.Listening)

	var vup UserInfo
	err := db.Database.
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

	listening := slices.Contains(listeningRooms, vup.RoomId)
	lastListenAt := GetLastListen(&vup, listening)

	return &UserDetailResp{
		UserResp: UserResp{
			UserInfo:       vup,
			Listening:      listening,
			LastListenedAt: lastListenAt,
		},
		BehavioursCount: map[string]int64{
			blive.DanmuMsg:         GetTotalCountByCommand(uid, blive.DanmuMsg),
			blive.InteractWord:     GetTotalCountByCommand(uid, blive.InteractWord),
			blive.SuperChatMessage: GetTotalCountByCommand(uid, blive.SuperChatMessage),
		},
	}, nil

}

func GetLastListen(vup *UserInfo, listening bool) time.Time {

	lastListenAt := time.Now()

	if !listening {

		var lastListen db.LastListen

		err := db.Database.
			Clauses(clause.OnConflict{DoNothing: true}).
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

func SearchVups(name string, page, pageSize int, orderBy string, desc bool) (*ListResp, error) {

	var vups []UserInfo

	order := "desc"

	if !desc {
		order = "asc"
	}

	err := db.Database.
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
		Joins("left join behaviours on behaviours.uid = vups.uid and behaviours.uid != behaviours.target_uid").
		Where("vups.name like ?", fmt.Sprintf("%%%s%%", name)).
		Group("vups.uid").
		Order(fmt.Sprintf("%s %s", orderBy, order)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&vups).Error

	if err != nil {
		return nil, err
	}

	var totalSearchCount int64

	err = db.Database.
		Model(&db.Vup{}).
		Select("count(*)").
		Where("name like ?", fmt.Sprintf("%%%s%%", name)).
		Find(&totalSearchCount).Error

	if err != nil {
		return nil, err
	}

	var userResps []*UserResp

	for _, vup := range vups {

		listeningRooms := *(statistics.Listening)

		var userResp UserResp

		userResp.UserInfo = vup

		listening := slices.Contains(listeningRooms, vup.RoomId)
		lastListenAt := GetLastListen(&vup, listening)

		userResp.Listening = listening
		userResp.LastListenedAt = lastListenAt

		userResps = append(userResps, &userResp)
	}

	return &ListResp{
		Page:    page,
		Size:    pageSize,
		MaxPage: int64(math.Ceil(float64(totalSearchCount)/float64(pageSize))) + 1,
		Total:   totalSearchCount,
		List:    userResps,
	}, nil
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

	listeningRooms := *(statistics.Listening)

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

		listening := slices.Contains(listeningRooms, info.RoomId)
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
		Limit(limit).
		Find(&mostDDVups).
		Error

	return mostDDVups, err
}
