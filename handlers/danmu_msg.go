package handlers

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
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

	// 先檢查 DD的人 是否在 vup 資料表中，如果是就記錄
	exist, err := vup.IsVup(uid)

	if err != nil {
		return err
	}

	// 不知名用戶
	if !exist {
		return nil
	}

	// 再檢查被DD的人 是否在 vup 資料表中，如果是就記錄
	exist, err = vup.IsVup(targetUid)

	if err != nil {
		return err
	}

	// 不知名用戶
	if !exist {
		return nil
	}

	var imageUrl = ""
	display := fmt.Sprintf("%s 在 %s 的直播间发送了一则消息: %s", uname, data.LiveInfo.Name, danmu)

	// 是表情包弹幕
	if obj, ok := base[13].(map[string]interface{}); ok {
		imageUrl = obj["url"].(string)
		display = fmt.Sprintf("%s 在 %s 的直播间发送了一则表情包:", uname, data.LiveInfo.Name)
	}

	logger.Info(display)

	behaviour := &db.Behaviour{
		Uid:       uid,
		CreatedAt: time.Now(),
		TargetUid: targetUid,
		Command:   blive.DanmuMsg,
		Display:   display,
		Image: sql.NullString{
			String: strings.Replace(imageUrl, "http://", "https://", -1),
			Valid:  imageUrl != "",
		},
	}

	result := db.Database.Create(behaviour)

	if result.Error != nil {
		logger.Warnf("記錄彈幕訊息行為到資料庫失敗: %v", result.Error)
	} else {
		logger.Infof("記錄彈幕訊息行為到資料庫成功。(%v 筆資料)", result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.DanmuMsg, danmuMsg)
}
