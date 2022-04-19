package statistics

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"gorm.io/gorm/clause"
	"time"
	"vup_dd_stats/service/db"
	"vup_dd_stats/utils/set"
)

var (
	logger                = logrus.WithField("service", "statistics")
	Listening    *[]int64 = &[]int64{}
	cooldownList          = set.New[int64]()
	allowRoles            = set.FromArray[int]([]int{1, 2, 3})
)

func StartListenStats(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()
	go fetchListeningInfo()
	go fetchVupListToRedis()
	for {
		select {
		case <-ticker.C:
			go fetchListeningInfo()
			go fetchVupListToRedis()
			go removeUnusedVupListFromRedis()
		case <-ctx.Done():
			return
		}
	}
}

func removeUnusedVupListFromRedis() {

	cache, err := db.SetGet(db.VupListKey)

	if err != nil {
		logger.Errorf("從 redis 獲取 虛擬主播列表 緩存時出現錯誤: %v", err)
		return
	}

	var vupList []string

	err = db.Database.
		Model(&db.Vup{}).
		Where("uid IN (?)", cache).
		Pluck("uid", &vupList).
		Error

	if err != nil {
		logger.Errorf("從資料庫獲取 虛擬主播列表時出現錯誤: %v", err)
		return
	}

	cacheSet, vupSet := set.FromArray(cache), set.FromArray(vupList)

	i := 0
	for _, v := range cacheSet.Difference(vupSet).ToArray() {
		err = db.SetRemove(db.VupListKey, v)
		if err != nil {
			logger.Errorf("從 redis 移除 %v 出虛擬主播列表時出現錯誤: %v", v, err)
			return
		} else {
			i++
		}
	}

	if i > 0 {
		logger.Infof("已成功把 %v 個不再為虛擬主播的用戶移除出 redis 緩存。", i)
	} else {
		logger.Debugf("已完成, 沒有需要移除出 redis 緩存的用戶。")
	}
}

func fetchVupListToRedis() {

	cache, err := db.SetGet(db.VupListKey)

	if err != nil {
		logger.Errorf("從 redis 獲取 虛擬主播列表 緩存時出現錯誤: %v", err)
		cache = make([]string, 0)
	}

	logger.Debugf("從 redis 獲取到的虛擬主播列表數量: %v", len(cache))

	var vupList []int64

	re := db.Database.Model(&db.Vup{})

	if len(cache) > 0 {
		re = re.Where("uid NOT IN (?)", cache)
	}

	re = re.Pluck("uid", &vupList)

	if re.Error != nil {
		logger.Errorf("從資料庫獲取vup列表錯誤: %v", re.Error)
		return
	}

	logger.Debugf("已成功提取 %v 列虛擬主播，即將加入到 redis 緩存...", re.RowsAffected)

	for _, vup := range vupList {
		err = db.SetAdd(db.VupListKey, fmt.Sprintf("%d", vup))
		if err != nil {
			logger.Errorf("添加 vup %v 到 redis 緩存失敗: %v", vup, err)
		} else {
			logger.Debugf("成功添加 vup %v 到 redis 緩存。", vup)
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

	result := db.Database.
		Model(&db.Vup{}).
		Where("room_id IN (?)", stats.Rooms).
		Select("room_id").
		Find(&roomIds)

	if result.Error != nil {
		logger.Errorf("從資料庫請求數據時出現錯誤: %v", result.Error)
		return
	}

	logger.Debugf("已成功從 %v 個監聽虛擬主播的房間中提取 %v 個房間在資料庫的監聽數據。", len(stats.Rooms), len(roomIds))

	vtbList, err := GetVtbListVtbMoe()

	if err != nil {
		logger.Errorf("請求vtb數據列表時出現錯誤: %v", err)
		vtbList = make([]VtbsMoeResp, 0)
	}

	roomSet := set.FromArray(roomIds)

	toBeInsert := make(map[int64]*db.Vup)

	userExist, userNotVtb := 0, 0

	// 只新增未有記錄的vup
	for _, room := range stats.Rooms {

		exist := roomSet.Has(room)

		if exist {
			userExist += 1
			//logger.Debugf("用戶已存在: %d", room)
			continue
		}

		liveInfo, err := GetLiveInfo(room)

		if err != nil {
			logger.Errorf("刷取房間 %v 的直播資訊時出現錯誤: %v", room, err)
			continue
		}

		if exist, err := db.SetContain(db.VupListKey, fmt.Sprintf("%v", liveInfo.UID)); err == nil && exist {
			userExist += 1
			//logger.Debugf("用戶已存在: %d", room)
			continue
		} else if err != nil {
			logger.Warnf("從 redis 緩存查找用戶時出現錯誤: %v", err)
		}

		found := false

		if !cooldownList.Has(liveInfo.UID) {

			user, err := GetUserInfo(liveInfo.UID)

			if err != nil {
				logger.Errorf("刷取房間 %v 的用戶資訊 (%v) 時出現錯誤: %v", room, liveInfo.Name, err)
				continue
			}

			// 請求頻繁
			if user.Code == -412 {
				logger.Warnf("用戶 %v(%v) 請求頻繁，已添加冷卻列表(十分鐘後)。", liveInfo.Name, liveInfo.UID)
				cooldownList.Add(liveInfo.UID)
				go func() {
					<-time.After(time.Minute * 10)
					cooldownList.Remove(liveInfo.UID)
				}()
				continue
			}

			// 有閃電的主播 + 不要机构认证
			if user.Code == 0 && allowRoles.Has(user.Data.Official.Role) {
				found = true
			}

		}

		// 如果先前已發現是有閃電主播，則無需再做過濾
		if !found {
			// 否則檢查是否在 vtb list 內
			for _, resp := range vtbList {
				if resp.Mid == liveInfo.UID {
					found = true
					break
				}
			}
		}

		// 不是 vtb
		if !found {
			userNotVtb += 1
			logger.Debugf("用戶不是vtb或高能主播: %d (%d)", room, liveInfo.UID)
			continue
		}

		vup := &db.Vup{
			Uid:           liveInfo.UID,
			Name:          liveInfo.Name,
			Face:          liveInfo.UserFace,
			FirstListenAt: time.Now(),
			RoomId:        liveInfo.RoomId,
			Sign:          liveInfo.UserDescription,
		}

		if err := db.SetAdd(db.VupListKey, fmt.Sprintf("%d", liveInfo.UID)); err != nil {
			logger.Errorf("儲存緩存到 redis 時出現錯誤: %v", err)
		}
		toBeInsert[liveInfo.UID] = vup
	}

	if len(toBeInsert) == 0 {
		logger.Infof("資料索取完畢，沒有需要新增的用戶資訊。")
		return
	}

	logger.Debugf("在 %v 個正在監聽的房間中，有 %v 個為資料庫已存在，有 %v 個不是vtb, 有 %v 需要被加入到資料庫", len(stats.Rooms), userExist, userNotVtb, len(toBeInsert))

	result = db.Database.
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(maps.Values(toBeInsert), len(toBeInsert))

	if result.Error != nil {
		logger.Errorf("插入數據到資料庫時出現錯誤: %v", result.Error)
		return
	} else if result.RowsAffected > 0 {
		logger.Infof("已成功插入 %v 筆用戶資訊到資料庫, %v 筆資料被忽略。", result.RowsAffected, int64(len(toBeInsert))-result.RowsAffected)
	} else {
		logger.Debugf("有 %v 筆資料因為資料相同被忽略。", int64(len(toBeInsert))-result.RowsAffected)
	}
}
