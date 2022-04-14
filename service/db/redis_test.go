package db

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func ainit() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("Error while loading environment file: %v", err)
	}
	InitRedis()
}
