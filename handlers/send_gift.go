package handlers

import (
	"fmt"
	"time"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
)

func giftMsg(data *blive.LiveData) error {
	d := data.Content["data"]

	var gift = &blive.SendGiftData{}

	if dict, ok := d.(map[string]interface{}); ok {
		if err := gift.Parse(dict); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("解析 Gift 数据失败")
	}

	// 送礼物的人
	uid := gift.UID
	// 收到礼物的人
	targetUid := data.LiveInfo.UID

	// 先检查 送礼物的人 是否在 vup 资料表中，如果是就记录
	exist, err := vup.IsVup(uid)

	if err != nil {
		return err
	}

	// 不知名用户
	if !exist {
		return nil
	}

	// 再检查 收到礼物的人 是否在 vup 资料表中，如果是就记录
	exist, err = vup.IsVup(targetUid)

	if err != nil {
		return err
	}

	// 不知名用户
	if !exist {
		return nil
	}

	giftName := gift.GiftName
	price := gift.Price

	display := fmt.Sprintf("在 %s 的直播间收到来自 %s 的 %s (%v元)", data.LiveInfo.Name, gift.Uname, giftName, price)

	// 将送礼行为记录到数据库
	behaviour := &db.Behaviour{
		Uid:       uid,
		CreatedAt: time.Now(),
		TargetUid: targetUid,
		Command:   blive.SendGift,
		Display:   display,
	}

	result := db.Database.Create(behaviour)

	if result.Error != nil {
		logger.Warnf("記錄醒目留言行為到資料庫失敗: %v", result.Error)
	} else {
		logger.Infof("記錄醒目留言行為到資料庫成功。(%v 筆資料)", result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.SendGift, giftMsg)
}
