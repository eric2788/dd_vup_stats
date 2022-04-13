package db

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"testing"
)

func aTestGetTakeIsVup(t *testing.T) {
	err := PutUserIsVup(690608693, true)
	if err != nil {
		t.Fatal(err)
	}
	isVup, ok := GetUserIsVup(690608693)
	if !ok {
		t.Fatal(err)
	}
	if isVup != true {
		t.Fatal("isVup is not true")
	}
}

func ainit() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("Error while loading environment file: %v", err)
	}
	InitRedis()
}
