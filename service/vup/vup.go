package vup

import (
	"errors"
	"fmt"
	"math"
	"time"
	"vup_dd_stats/service/analysis"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
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

	listening := stats.Listening.Has(vup.RoomId)

	if !listening && vup.LastListenedAt.IsZero() {
		vup.LastListenedAt = time.Now().UTC()
	}

	if listening {
		go UpdateLastListens([]int64{uid}, []int64{})
	} else {
		go UpdateLastListens([]int64{}, []int64{uid})
	}

	behaviourCounts := make(map[string]stats.TotalStats)
	commandStats, err := GetTotalCommandStats(uid)
	if err != nil {
		logger.Errorf("嘗試獲取用戶 %v 的行为统计時出現錯誤: %v", uid, err)
	} else {
		for _, stat := range commandStats {
			behaviourCounts[stat.Command] = stat
		}

		// for those didn't appear from table
		for _, stat := range blive.GetRegisteredCommands() {
			if _, ok := behaviourCounts[stat]; !ok {
				behaviourCounts[stat] = stats.TotalStats{
					Command: stat,
					Count:   0,
					Price:   0,
				}
			}
		}
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

func UpdateLastListens(listening, unListening []int64) {

	if len(listening) > 0 {
		re := db.Database.
			Clauses(clause.OnConflict{DoNothing: true}).
			Where("uid IN ?", listening).
			Delete(&db.LastListen{})

		if re.Error != nil {
			logger.Errorf("嘗試刪除最後監聽訊息時出現錯誤: %v", re.Error)
		}

		if re.RowsAffected > 0 {
			logrus.Debugf("Successfully Delete %v record of %v because it is listening.", re.RowsAffected, listening)
		}

	}

	if len(unListening) > 0 {

		insert := make([]db.LastListen, len(unListening))
		for i, uid := range unListening {
			insert[i] = db.LastListen{
				Uid:          uid,
				LastListenAt: time.Now().UTC(),
			}
		}

		re := db.Database.
			Clauses(clause.OnConflict{DoNothing: true}).
			CreateInBatches(insert, len(unListening))

		if re.Error != nil {
			logger.Errorf("嘗試更新最後監聽訊息時出現錯誤: %v", re.Error)
		}

		if re.RowsAffected > 0 {
			logrus.Debugf("Successfully Update %v record of %v because it is unListening.", re.RowsAffected, unListening)
		}
	}

}

// SearchVups search vups by name
func SearchVups(name string, page, pageSize int, orderBy string, desc bool) (*stats.ListResp[UserResp], error) {

	// ensure page is valid
	page = int(math.Max(1, float64(page)))

	//ensure pageSize is valid
	pageSize = int(math.Max(1, float64(pageSize)))

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

	var vups []UserInfo
	var totalSearchCount int64

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

	err = db.Database.
		Model(&db.Vup{}).
		Select("count(*)").
		Where("name like ?", fmt.Sprintf("%%%s%%", name)).
		Find(&totalSearchCount).Error

	if err != nil {
		return nil, err
	}

	var userResps []UserResp

	var listening []int64
	var unListening []int64

	for _, vup := range vups {

		var userResp UserResp

		userResp.UserInfo = vup
		userResp.Listening = stats.Listening.Has(vup.RoomId)

		// if not listening and not find from last_listens table: use now as last_listened_at
		if !userResp.Listening && userResp.LastListenedAt.IsZero() {
			userResp.LastListenedAt = time.Now().UTC()
		}

		if userResp.Listening {
			listening = append(listening, vup.Uid)
		} else {
			unListening = append(unListening, vup.Uid)
		}

		userResps = append(userResps, userResp)
	}

	go UpdateLastListens(listening, unListening)

	// annoymous record
	go analysis.RecordSearchText(name, totalSearchCount)

	return &stats.ListResp[UserResp]{
		Page:    page,
		Size:    pageSize,
		MaxPage: int64(math.Ceil(float64(totalSearchCount) / float64(pageSize))),
		Total:   totalSearchCount,
		List:    userResps,
	}, nil
}
