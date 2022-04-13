package statistics

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"gorm.io/gorm/clause"
	"time"
	db2 "vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
	"vup_dd_stats/utils/set"
)

var (
	logger             = logrus.WithField("service", "statistics")
	Listening *[]int64 = &[]int64{}
)

func StartListenStats(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()
	go fetchListeningInfo()
	for {
		select {
		case <-ticker.C:
			go fetchListeningInfo()
		case <-ctx.Done():
			return
		}
	}
}

func fetchListeningInfo() {

	stats, err := GetListening()
	if err != nil {
		logger.Errorf("刷取監控訊息時出現錯誤: %v", err)
		return
	}

	var roomIds []int64

	Listening = &stats.Rooms

	result := db2.Database.
		Model(&db2.Vup{}).
		Where("room_id IN ?", stats.Rooms).
		Select("room_id").
		Find(&roomIds)

	if result.Error != nil {
		logger.Errorf("從資料庫請求數據時出現錯誤: %v", result.Error)
		return
	}

	vtbList, err := GetVtbListVtbMoe()

	if err != nil {
		logger.Errorf("請求vtb數據列表時出現錯誤: %v", err)
		vtbList = make([]VtbsMoeResp, 0)
	}

	roomSet := set.FromArray(roomIds)

	toBeInsert := make(map[int64]*db2.Vup)

	// 只新增未有記錄的vup
	for _, room := range stats.Rooms {

		exist := roomSet.Has(room)

		if exist {
			logger.Debugf("用戶已存在: %d", room)
			continue
		}

		liveInfo, err := GetLiveInfo(room)

		if err != nil {
			logger.Errorf("刷取房間 %v 的直播資訊時出現錯誤: %v", room, err)
			continue
		}

		found := false
		for _, resp := range vtbList {
			if resp.Mid == liveInfo.UID {
				found = true
				break
			}
		}

		// 不是 vtb
		if !found {
			logger.Debugf("用戶不是vtb: %d", room)
			vup.Caches.Store(liveInfo.UID, false)
			continue
		}

		v := &db2.Vup{
			Uid:           liveInfo.UID,
			Name:          liveInfo.Name,
			Face:          liveInfo.UserFace,
			FirstListenAt: time.Now(),
			RoomId:        liveInfo.RoomId,
			Sign:          liveInfo.UserDescription,
		}

		vup.Caches.Store(liveInfo.UID, true)
		toBeInsert[liveInfo.UID] = v
	}

	if len(toBeInsert) == 0 {
		logger.Infof("資料索取完畢，沒有需要新增的用戶資訊。")
		return
	}

	logger.Debugf("即將插入 %v 筆用戶資料到資料庫", len(toBeInsert))

	result = db2.Database.
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(maps.Values(toBeInsert), len(toBeInsert))

	if result.Error != nil {
		logger.Errorf("插入數據到資料庫時出現錯誤: %v", result.Error)
		return
	} else if result.RowsAffected > 0 {
		logger.Infof("已成功插入 %v 筆用戶資訊到資料庫, %v 筆資料被忽略。", result.RowsAffected, int64(len(toBeInsert))-result.RowsAffected)
	}
}
