package handlers

import (
	"fmt"
	"time"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
)

func roomEnter(data *blive.LiveData) error {

	entered := data.Content["data"].(map[string]interface{})
	uname := entered["uname"].(string)
	// DD 的人
	uid := int64(entered["uid"].(float64))
	// 被 DD 的人
	targetUid := data.LiveInfo.UID

	// 先檢查 DD的人 是否在 vup 資料表中，如果是就記錄
	exist, err := vup.IsVup(uid)

	if err != nil {
		return err
	}

	// 不知名用戶
	if !exist {
		return nil
	}

	// 先檢查 被DD的人 是否在 vup 資料表中，如果是就記錄
	exist, err = vup.IsVup(targetUid)

	if err != nil {
		return err
	}

	// 不知名用戶
	if !exist {
		return nil
	}

	display := fmt.Sprintf("%s 进入了 %s 的直播间", uname, data.LiveInfo.Name)
	logger.Infof(display)

	// 將資料寫入資料庫

	behaviour := &db.Behaviour{
		Uid:       uid,
		CreatedAt: time.Now().UTC(),
		TargetUid: targetUid,
		Command:   data.Command,
		Display:   display,
	}

	result := db.Database.Create(behaviour)

	if result.Error != nil {
		logger.Warnf("記錄進入房間行為到資料庫失敗: %v", result.Error)
	} else {
		logger.Infof("記錄進入房間行為到資料庫成功。(%v 筆資料)", result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.InteractWord, roomEnter)
}
