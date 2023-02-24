package watcher

import (
	"context"
	"time"
	"vup_dd_stats/service/db"
)

// create a queue of watcher_behaviour records to be saved to the database
// to avoid hitting the database too often and huge performance issues

var (
	watcherBehaviourQueue = make(chan *db.WatcherBehaviour, 2048)
)

func SaveWatcherBehaviour(wb *db.WatcherBehaviour) {
	watcherBehaviourQueue <- wb
}

// RunSaveTimer save the watcher_behaviour records to the database
// this is run in a goroutine
func RunSaveTimer(ctx context.Context) {
	timer := time.NewTicker(5 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			insertWatchers()
			return
		case <-timer.C:
			insertWatchers()
		}
	}
}

func insertWatchers() {
	inserts := make([]*db.WatcherBehaviour, 0)
	for watcher := range watcherBehaviourQueue {
		inserts = append(inserts, watcher)
		if len(watcherBehaviourQueue) == 0 {
			break
		}
	}
	if len(inserts) == 0 {
		logger.Debugf("没有可以插入的 watcher_behaviour 数据, 跳过")
		return
	}

	// when it reached the maximum number of inserts in a single query
	for len(inserts) >= 10000 {
		// split the inserts
		go insertRecords(inserts[:10000])
		<-time.After(time.Second * 3)
		inserts = inserts[10000:]
	}
}

func insertRecords(records []*db.WatcherBehaviour) {
	result := db.Database.CreateInBatches(records, len(records))
	if result.Error != nil {
		logger.Errorf("寫入 watcher_behaviour 記錄失敗: %v", result.Error)
	} else if result.RowsAffected > 0 {
		logger.Infof("成功寫入 %d 筆 watcher_behaviour 的記錄。", result.RowsAffected)
	}
}
