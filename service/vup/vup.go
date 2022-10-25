package vup

import (
	"errors"
	"fmt"
	"math"
	"time"
	"vup_dd_stats/service/analysis"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/statistics"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	logger = logrus.WithField("service", "vup")
)

func IsVup(uid int64) (bool, error) {

	if re, err := db.SetContain(db.VupListKey, fmt.Sprintf("%d", uid)); err == nil {
		return re, nil
	} else if err != nil && err != redis.Nil {
		logger.Errorf("從 redis 提取緩存錯誤: %v, 將使用資料庫", err)
	}

	var exist bool

	err := db.Database.
		Model(&db.Vup{}).
		Select("count(*) > 0").
		Where("uid = ?", uid).
		Find(&exist).
		Error

	if err != nil {
		return false, err
	}

	if exist {
		if err := db.SetAdd(db.VupListKey, fmt.Sprintf("%d", uid)); err != nil {
			logger.Errorf("儲存用戶 %v 到 redis 時出現錯誤: %v", uid, err)
		} else {
			logger.Debugf("從 IsVup 新增了 %v 到 redis", uid)
		}
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
			"MAX(last_listens.last_listen_at) AS last_listened_at",
			"SUM(behaviours.price) AS total_spent",
		}).
		Joins("left join behaviours on behaviours.uid = vups.uid and behaviours.uid != behaviours.target_uid").
		Joins("left join last_listens on last_listens.uid = vups.uid").
		Group("behaviours.uid, vups.uid").
		Take(&vup).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) || vup.Uid == 0 {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	listening := slices.Contains(listeningRooms, vup.RoomId)
	lastListenAt := GetLastListen(&vup, listening)

	if !listening {
		vup.LastListenedAt = lastListenAt
	}

	registeredCommands := blive.GetRegisteredCommands()
	behaviourCounts := make(map[string]TotalStats, len(registeredCommands))
	for _, command := range registeredCommands {
		behaviourCounts[command] = GetTotalStatusByCommand(uid, command)
	}

	// annoymous record
	go analysis.RecordSearchUser(uid, vup.Name)

	return &UserDetailResp{
		UserResp: UserResp{
			UserInfo:  vup,
			Listening: listening,
		},
		BehavioursCount: behaviourCounts,
	}, nil

}

func GetLastListen(vup *UserInfo, listening bool) time.Time {

	lastListenAt := time.Now()

	if !listening {

		var lastListen = &db.LastListen{
			Uid:          vup.Uid,
			LastListenAt: lastListenAt,
		}

		// FirstOrCreate will throw duplicate entry error if the record already exists
		err := db.Database.Take(lastListen, vup.Uid).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Debugf("Record of %v not found, create new one", vup.Name)
			err = db.Database.Create(lastListen).Error
		}

		if err != nil {
			logger.Errorf("嘗試插入最後監聽訊息時出現錯誤: %v", err)
			logrus.Debugf("Record Insert Error, using first listen at")
			return vup.FirstListenAt
		} else {
			return lastListen.LastListenAt
		}

	} else {

		re := db.Database.
			Clauses(clause.OnConflict{DoNothing: true}).
			Delete(&db.LastListen{}, vup.Uid)

		if re.Error != nil {
			logger.Errorf("嘗試刪除最後監聽訊息時出現錯誤: %v", re.Error)
		}

		if re.RowsAffected > 0 {
			logrus.Debugf("Successfully Delete %v record of %v because it is listening.", re.RowsAffected, vup.Name)
		}

	}

	return lastListenAt
}

func SearchVups(name string, page, pageSize int, orderBy string, desc bool) (*ListResp[UserResp], error) {

	// ensure page is valid
	page = int(math.Max(1, float64(page)))

	//ensure pageSize is valid
	pageSize = int(math.Max(1, float64(pageSize)))

	var vups []UserInfo

	order := "desc"

	if !desc {
		order = "asc"
	}

	// ==============
	// postgres only

	var nullsLast = ""

	if db.DatabaseType == "postgres" {
		nullsLast = " NULLS LAST"
	}

	// ==============

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
			"MAX(last_listens.last_listen_at) AS last_listened_at",
			"SUM(behaviours.price) AS total_spent",
		}).
		Joins("left join behaviours on behaviours.uid = vups.uid and behaviours.uid != behaviours.target_uid").
		Joins("left join last_listens on last_listens.uid = vups.uid").
		Where("vups.name like ?", fmt.Sprintf("%%%s%%", name)).
		Group("behaviours.uid, vups.uid").
		Order(fmt.Sprintf("%s %s%s", orderBy, order, nullsLast)).
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

	var userResps []UserResp

	for _, vup := range vups {

		listeningRooms := *(statistics.Listening)

		var userResp UserResp

		userResp.UserInfo = vup

		listening := slices.Contains(listeningRooms, vup.RoomId)
		lastListenAt := GetLastListen(&vup, listening)

		userResp.Listening = listening
		userResp.LastListenedAt = lastListenAt

		userResps = append(userResps, userResp)
	}

	// annoymous record
	go analysis.RecordSearchText(name, totalSearchCount)

	return &ListResp[UserResp]{
		Page:    page,
		Size:    pageSize,
		MaxPage: int64(math.Ceil(float64(totalSearchCount) / float64(pageSize))),
		Total:   totalSearchCount,
		List:    userResps,
	}, nil
}

// GetMostDDVups 獲取進入最多不同直播間的 vups
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
		Group("behaviours.uid, vups.uid").
		Order("count desc").
		Limit(limit).
		Find(&mostDDVups).
		Error

	return mostDDVups, err
}
