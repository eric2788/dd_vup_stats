package handlers

import (
	"fmt"
	"time"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/vup"
	"vup_dd_stats/service/watcher"
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

	// 先检查 收到舰长的人 是否在 vup 资料表中，如果是就记录
	isVup, err := vup.IsVup(targetUid)

	if err != nil {
		return err
	}

	// 不知名用户
	if !isVup {
		return nil
	}

	// 再检查 送舰长的人 是否在 vup 资料表中，如果是就记录
	isVup, err = vup.IsVup(uid)

	if err != nil {
		return err
	}

	var log = logger.Infof
	if !isVup {
		log = logger.Debugf
	}

	guardName := guardBuy.RoleName

	display := fmt.Sprintf("在 %s 的直播间收到来自 %s 的 %s", data.LiveInfo.Name, guardBuy.Username, guardName)

	log(display)

	behaviour := &db.Behaviour{
		Uid:       uid,
		CreatedAt: time.Now().UTC(),
		TargetUid: targetUid,
		Command:   data.Command,
		Display:   display,
		Price:     float64(guardBuy.Price / 1000),
	}

	if isVup {
		go vup.InsertBehaviour(behaviour)
	} else {
		go watcher.SaveWatcherBehaviour(behaviour.ToWatcherBehaviour(guardBuy.Username))
	}

	return nil
}

func init() {
	blive.RegisterHandler(blive.GuardBuyToast, guardBuyMsg)
}
