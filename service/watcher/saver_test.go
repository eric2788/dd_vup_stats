package watcher_test

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	queue = make(chan string, 10)
	wg    = &sync.WaitGroup{}
)

func TestSaveWatcherQueue(t *testing.T) {
	go func() {
		for i := 0; i < 100; i++ {
			wg.Wait()
			queue <- fmt.Sprintf("test-%d", i)
		}
	}()

	go func() {
		timer := time.NewTicker(2 * time.Second)
		defer timer.Stop()
		for {
			t.Log("wait")
			<-timer.C
			wg.Add(1)
			t.Log(len(queue))
			for a := range queue {

				t.Log(a, len(queue))
				if len(queue) == 0 {
					break
				}
			}
			wg.Done()
		}
	}()

	<-time.After(time.Second * 10)
}
