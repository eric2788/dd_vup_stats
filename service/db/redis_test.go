package db

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Errorf("Error while loading environment file: %v", err)
	}
	if os.Getenv("REDIS_ADDR") == "" {
		return
	}
	InitRedis()
}
