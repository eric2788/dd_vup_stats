package watcher

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"vup_dd_stats/service/db"
)

// create a queue of watcher_behaviour records to be saved to the database
// to avoid hitting the database too often and huge performance issues

var (
	watcherBehaviourQueue = make(chan *db.WatcherBehaviour, 4096)
	writing               atomic.Bool
	maxBuffer             int
)

func SaveWatcherBehaviour(wb *db.WatcherBehaviour) {
	for writing.Load() || len(watcherBehaviourQueue) > 4000 {
		<-time.After(time.Second)
	}
	watcherBehaviourQueue <- wb
}

// RunSaveTimer save the watcher_behaviour records to the database
// this is run in a goroutine
func RunSaveTimer(ctx context.Context) {
	logger.Infof("開始運行 watcher_behaviour 記錄寫入程序...")
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

	if len(watcherBehaviourQueue) == 0 {
		logger.Debugf("没有可以插入的 watcher_behaviour 数据, 跳过")
		return
	}

	defer writing.Store(false)
	writing.Store(true)

	queueSize := len(watcherBehaviourQueue)

	inserts := make([]*db.WatcherBehaviour, 0)
	for watcher := range watcherBehaviourQueue {

		// remove 0x00 byte null characters
		watcher.Display = strings.Replace(watcher.Display, "\x00", "", -1)

		inserts = append(inserts, watcher)
		if len(watcherBehaviourQueue) == 0 {
			break
		}
	}

	logger.Debugf("即將寫入 %v -> %v 個 watcher_behaviours 記錄...", queueSize, len(inserts))

	// when it reached the maximum number of inserts in a single query
	for len(inserts) > maxBuffer {
		// split the inserts
		insertRecords(inserts[:maxBuffer])
		<-time.After(time.Second * 5)
		inserts = inserts[maxBuffer:]
		logger.Debugf("剩余 %v 個 watcher_behaviours 記錄...", len(inserts))
	}

	insertRecords(inserts)
}

func insertRecords(records []*db.WatcherBehaviour) {
	result := db.Database.CreateInBatches(records, len(records))
	if result.Error != nil {
		logger.Errorf("寫入 watcher_behaviour 記錄失敗: %v, 10 分鐘後重試寫入", result.Error)
		go func() {
			<-time.After(10 * time.Minute)
			insertRecords(records)
		}()
	} else if result.RowsAffected > 0 {
		logger.Infof("成功寫入 %d 筆 watcher_behaviour 的記錄。", result.RowsAffected)
	}
}

func init() {
	b, err := strconv.Atoi(os.Getenv("MAX_BUFFER"))
	if err != nil {
		logger.Errorf("error parsing MAX_BUFFER: %v, will use default value 1000", err)
		b = 1000
	}
	maxBuffer = b
}
