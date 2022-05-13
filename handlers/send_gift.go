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

	// 对礼物进行筛选，如小心心等不应记录到数据库中
	if !filterGift(gift) {
		return nil
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
	number := gift.Num
	price := gift.Price

	display := fmt.Sprintf("在 %s 的直播间收到来自 %s 的 %v 个 %s (%v元)", data.LiveInfo.Name, gift.Uname, number, giftName, price)
	logger.Infof(display)

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
		logger.Warnf("記錄送礼行為到資料庫失敗: %v", result.Error)
	} else {
		logger.Infof("記錄送礼行为到資料庫成功。(%v 筆資料)", result.RowsAffected)
	}

	return nil
}

func filterGift(gift *blive.SendGiftData) bool {
	giftName := gift.GiftName
	filter_gift_list := []string{"小心心", "辣条", "小花花"}
	if !in(giftName, filter_gift_list) {
		return true
	}
	return false

}

func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

func init() {
	blive.RegisterHandler(blive.SendGift, giftMsg)
}
