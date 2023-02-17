package watcher_test

import (
	"context"
	"fmt"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

var (
	queue = make(chan string, 10)
	ctx   context.Context
)

func TestSaveWatcherQueue(t *testing.T) {
	go func() {
		for i := 0; i < 100; i++ {
			<-time.After(time.Duration(rand.Int63n(1000)) * time.Millisecond)
			if ctx != nil {
				<-ctx.Done()
			}
			queue <- fmt.Sprintf("test-%d", i)
		}
	}()

	go func() {
		timer := time.NewTicker(2 * time.Second)
		defer timer.Stop()
		for {
			t.Log("wait")
			<-timer.C
			c, cancel := context.WithCancel(context.Background())
			ctx = c
			t.Log(len(queue))
			for a := range queue {
				t.Log(a, len(queue))
				if len(queue) == 0 {
					break
				}
			}
			cancel()
		}
	}()

	<-time.After(time.Second * 10)
}
