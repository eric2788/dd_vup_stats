package main

import (
	"context"
	"os"
	"sync"
	"vup_dd_stats/controller/records"
	"vup_dd_stats/controller/stats"
	"vup_dd_stats/controller/user"
	"vup_dd_stats/controller/watcher"
	"vup_dd_stats/service/blive"
	"vup_dd_stats/service/db"
	statistics "vup_dd_stats/service/stats"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "vup_dd_stats/handlers"
)

func init() {

	if _, err := os.Open(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			logrus.Errorf("Error while loading environment file: %v", err)
		} else {
			logrus.Infof("Environment File (.env) loaded.")
		}
	}

	if os.Getenv("GIN_MODE") != "release" {
		// debug mode
		logrus.SetLevel(logrus.DebugLevel)
	}

}

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}

	db.InitDB()
	db.InitRedis()

	go blive.StartWebSocket(ctx, wg)
	go statistics.StartListenStats(ctx)

	router := gin.New()

	router.Use(func(c *gin.Context) {
		if c.Query("log") == "false" {
			c.Next()
			return
		}
		gin.Logger()(c)
	})

	router.Use(CORS())
	router.Use(ErrorHandler())

	user.Register(router.Group("/user"))
	stats.Register(router.Group("/stats"))
	records.Register(router.Group("/records"))
	watcher.Register(router.Group("/watcher"))

	if err := router.Run(":8086"); err != nil {
		logrus.Errorf("Error while starting server: %v", err)
	}

	cancel()
	wg.Wait()
}

func CORS() gin.HandlerFunc {
	def := cors.DefaultConfig()
	return cors.New(cors.Config{
		AllowOrigins: []string{
			"https://ddstats.ericlamm.xyz",
			"https://ddstats.pages.dev",
			os.Getenv("DEV_HOST"),
		},
		AllowWebSockets: true,
		AllowMethods:    def.AllowMethods,
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Origin",
			"Content-Length",
		},
	})
}
