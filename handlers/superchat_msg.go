package handlers

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
	"vup_dd_stats/service/watcher"
)

func superChatMsg(data *blive.LiveData) error {

	d := data.Content["data"]

	var superchat = &blive.SuperChatMessageData{}

	if dict, ok := d.(map[string]interface{}); ok {
		if err := superchat.Parse(dict); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("解析 SuperChat 數據失敗")
	}

	// DD 的人
	uid := superchat.UID
	// 被 DD 的人
	targetUid := data.LiveInfo.UID

	// 先檢查被DD的人 是否在 vup 資料表中，如果是就記錄
	isVup, err := vup.IsVup(targetUid)

	if err != nil {
		return err
	}

	// 不知名用戶
	if !isVup {
		return nil
	}

	// 再檢查 DD的人 是否在 vup 資料表中，如果是就記錄
	isVup, err = vup.IsVup(uid)

	if err != nil {
		return err
	}

	var log = logger.Infof
	if !isVup {
		log = logger.Debugf
	}

	message := superchat.Message
	price := superchat.Price

	display := fmt.Sprintf("在 %s 的直播间收到来自 %s 的 %v 元醒目留言: %s", data.LiveInfo.Name, superchat.UserInfo.UName, price, message)
	log(display)

	// 將訊息記錄到資料庫
	behaviour := &db.Behaviour{
		Uid:       uid,
		CreatedAt: time.Now().UTC(),
		TargetUid: targetUid,
		Command:   data.Command,
		Display:   display,
		Price:     float64(price),
	}

	var result *gorm.DB

	if isVup {
		result = db.Database.Create(behaviour)
	} else {
		watcher.SaveWatcher(uid)
		result = db.Database.Create(behaviour.ToWatcherBehaviour())
	}

	if result.Error != nil {
		logger.Warnf("記錄醒目留言行為到資料庫失敗: %v", result.Error)
	} else {
		log("記錄醒目留言行為到資料庫成功。(%v 筆資料)", result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.SuperChatMessage, superChatMsg)
}
