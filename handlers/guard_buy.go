package handlers

import (
	"fmt"
	"time"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
)

func guardBuyMsg(data *blive.LiveData) error {
	d := data.Content["data"]

	var guardBuy = &blive.GuardBuyData{}

	if dict, ok := d.(map[string]interface{}); ok {
		if err := guardBuy.Parse(dict); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("解析 GuardBuy 数据失败")
	}

	// 上舰长的人
	uid := guardBuy.UID
	// 收到舰长的人
	targetUid := data.LiveInfo.UID

	// 先检查 送舰长的人 是否在 vup 资料表中，如果是就记录
	exist, err := vup.IsVup(uid)

	if err != nil {
		return err
	}

	// 不知名用户
	if !exist {
		return nil
	}

	// 再检查 收到舰长的人 是否在 vup 资料表中，如果是就记录
	exist, err = vup.IsVup(targetUid)

	if err != nil {
		return err
	}

	// 不知名用户
	if !exist {
		return nil
	}

	guardName := guardBuy.RoleName

	display := fmt.Sprintf("在 %s 的直播间收到来自 %s 的 %s", data.LiveInfo.Name, guardBuy.Username, guardName)

	behaviour := &db.Behaviour{
		Uid:       uid,
		CreatedAt: time.Now(),
		TargetUid: targetUid,
		Command:   blive.GuardBuyToast,
		Display:   display,
		Price:     float64(guardBuy.Price / 1000),
	}

	result := db.Database.Create(behaviour)

	if result.Error != nil {
		logger.Warnf("記錄上舰行為到資料庫失敗: %v", result.Error)
	} else {
		logger.Infof("記錄上舰行為到資料庫成功。(%v 筆資料)", result.RowsAffected)
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.GuardBuyToast, guardBuyMsg)
}
