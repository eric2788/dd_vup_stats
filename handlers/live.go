package handlers

import (
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/statistics"
	"vup_dd_stats/service/vup"
)

func updateInfo(data *blive.LiveData) error {

	// filter duplicate
	if _, ok := data.Content["live_time"]; !ok {
		return nil
	}

	if exist, err := vup.IsVup(data.LiveInfo.UID); err != nil {
		return err
	} else if !exist {
		return nil
	}

	info, err := statistics.GetListeningInfo(data.LiveInfo.RoomId)

	if err != nil {
		logger.Warnf("刷新 %v 的用戶資訊時出現錯誤: %v, 已略過更新。", data.LiveInfo.Name, err)
		return nil
	}

	v := &db.Vup{
		Uid:    info.UID,
		Name:   info.Name,
		Face:   info.UserFace,
		RoomId: info.RoomId,
		Sign:   info.UserDescription,
	}

	result := db.Database.Updates(v)

	if result.Error != nil {
		logger.Warnf("更新 %v 的用戶資訊到資料庫時出現錯誤: %v", data.LiveInfo.Name, result.Error)
	} else if result.RowsAffected > 0 {
		logger.Infof("已更新 %s 的用戶資訊到資料庫。(%v 筆資料)", info.Name, result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.Live, updateInfo)
	blive.RegisterHandler(blive.Preparing, updateInfo)
}
