package vup

import (
	"context"
	"time"
	"vup_dd_stats/service/db"
)

// because behaviours is not much frequent as watcher_behaviours, so we can use without lock
var behaviourQueue = make(chan *db.Behaviour, 1000)

func InsertBehaviour(b *db.Behaviour) {
	behaviourQueue <- b
}

func RunSaveTimer(ctx context.Context) {
	logger.Infof("開始運行 behaviour 記錄寫入程序...")
	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			insertBehaviours()
			return
		case <-timer.C:
			insertBehaviours()
		}
	}
}

func insertBehaviours() {
	if len(behaviourQueue) == 0 {
		logger.Debugf("没有可以插入的 behaviour 数据, 跳过")
		return
	}

	queueSize := len(behaviourQueue)
	inserts := make([]*db.Behaviour, 0)

	for behaviour := range behaviourQueue {
		inserts = append(inserts, behaviour)
		if len(behaviourQueue) == 0 {
			break
		}
	}

	logger.Debugf("即將寫入 %v -> %v 個 behaviours 記錄...", queueSize, len(inserts))
	insertRecords(inserts)
}

func insertRecords(records []*db.Behaviour) {
	result := db.Database.CreateInBatches(records, len(records))
	if result.Error != nil {
		logger.Errorf("寫入 behaviour 記錄失敗: %v, 10 分鐘後重試寫入", result.Error)
		go func() {
			<-time.After(10 * time.Minute)
			insertRecords(records)
		}()
	} else if result.RowsAffected > 0 {
		logger.Infof("成功寫入 %d 筆 behaviour 的記錄。", result.RowsAffected)
	}
}
