package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"

	_ "vup_dd_stats/handlers"
)

func init() {
	if os.Getenv("GIN_MODE") != "release" {
		// debug mode
		logrus.SetLevel(logrus.DebugLevel)
	}

	if _, err := os.Open(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			logrus.Errorf("Error while loading environment file: %v", err)
		}
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}

	db.InitDB()

	go blive.StartWebSocket(ctx, wg)
	go stats.StartListenStats(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

	cancel()
	wg.Wait()
}
