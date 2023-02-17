package handlers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
	"vup_dd_stats/service/watcher"

	"gorm.io/gorm"
)

func danmuMsg(data *blive.LiveData) error {

	info := data.Content["info"].([]interface{})

	base := info[0].([]interface{})

	if base[9].(float64) != 0 {
		// 抽獎/紅包彈幕
		return nil
	}

	userInfo := info[2].([]interface{})

	// 被DD的人
	targetUid := data.LiveInfo.UID

	danmu := info[1].(string)
	uname := userInfo[1].(string)

	// DD的人
	uid := int64(userInfo[0].(float64))

	// 先檢查被DD的人 是否在 vup 資料表中，如果是就記錄
	isVup, err := vup.IsVup(targetUid)

	if err != nil {
		return err
	}

	// 不知名用戶
	if !isVup {
		return nil
	}

	// 再檢查DD的人 是否在 vup 資料表中，如果是就記錄
	isVup, err = vup.IsVup(uid)

	if err != nil {
		return err
	}

	var log = logger.Infof
	if !isVup {
		log = logger.Debugf
	}

	var imageUrl = ""
	display := fmt.Sprintf("%s 在 %s 的直播间发送了一则消息: %s", uname, data.LiveInfo.Name, danmu)

	// 是表情包弹幕
	if obj, ok := base[13].(map[string]interface{}); ok {
		imageUrl = obj["url"].(string)
		display = fmt.Sprintf("%s 在 %s 的直播间发送了一则表情包:", uname, data.LiveInfo.Name)
		log("%s 在 %s 的直播间发送了一则表情包: [%s]", uname, data.LiveInfo.Name, danmu)
	} else {
		log(display)
	}

	behaviour := &db.Behaviour{
		Uid:       uid,
		CreatedAt: time.Now().UTC(),
		TargetUid: targetUid,
		Command:   data.Command,
		Display:   display,
		Image: sql.NullString{
			String: strings.Replace(imageUrl, "http://", "https://", -1),
			Valid:  imageUrl != "",
		},
	}

	var result *gorm.DB

	if isVup {
		result = db.Database.Create(behaviour)
	} else {
		go watcher.SaveWatcherBehaviour(behaviour.ToWatcherBehaviour(uname))
		return nil
	}

	if result.Error != nil {
		logger.Warnf("記錄彈幕訊息行為到資料庫失敗: %v", result.Error)
	} else {
		log("記錄彈幕訊息行為到資料庫成功。(%v 筆資料)", result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.DanmuMsg, danmuMsg)
}
