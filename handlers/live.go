package handlers

import (
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"
)

func onLive(data *blive.LiveData) error {

	logger.Infof("%v 開播，正在刷新其用戶資訊...", data.LiveInfo.Name)
	info, err := stats.GetLiveInfo(data.LiveInfo.RoomId)

	if err != nil {
		logger.Warnf("刷新 %v 的用戶資訊時出現錯誤: %v, 已略過更新。", data.LiveInfo.Name, err)
		return nil
	}

	vup := &db.Vup{
		Uid:    info.UID,
		Name:   info.Name,
		Face:   info.UserFace,
		RoomId: info.RoomId,
		Sign:   info.UserDescription,
	}

	result := db.Database.Updates(vup)

	if result.Error != nil {
		logger.Warnf("更新 %v 的用戶資訊到資料庫時出現錯誤: %v", data.LiveInfo.Name, result.Error)
	} else {
		logger.Infof("已更新 %v 筆資料到資料庫。", result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.Live, onLive)
}
