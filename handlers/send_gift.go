package handlers

import (
	"fmt"
	"strconv"
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
	// 对免费礼物进行筛选，如小心心等不应记录到数据库中
	if gift.CoinType == "silver" {
		logger.Debugf("%s 的禮物價值為銀瓜子類別, 已略過。", gift.GiftName)
		return nil
	}

	// 1000 coins / 100 = 10 電池
	batteries := gift.TotalCoin / 100

	// 10 電池 = 1 元
	price, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(batteries)/10), 64)

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

	display := fmt.Sprintf("在 %s 的直播间收到来自 %s 的 %v 个 %s (%v元)", data.LiveInfo.Name, gift.Uname, number, giftName, price)
	logger.Infof(display)

	// 将送礼行为记录到数据库
	behaviour := &db.Behaviour{
		Uid:       uid,
		CreatedAt: time.Now().UTC(),
		TargetUid: targetUid,
		Command:   data.Command,
		Display:   display,
		Price:     price,
	}

	result := db.Database.Create(behaviour)

	if result.Error != nil {
		logger.Warnf("記錄送礼行為到資料庫失敗: %v", result.Error)
	} else {
		logger.Infof("記錄送礼行为到資料庫成功。(%v 筆資料)", result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.SendGift, giftMsg)
}
