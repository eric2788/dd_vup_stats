package db

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"testing"
)

type Stats struct {
	Rooms []int64
}

func aTestInQuery(t *testing.T) {
	var vups []int64

	var rooms []int64

	r := Database.
		Model(&Vup{}).
		Pluck("room_id", &rooms)

	if r.Error != nil {
		logrus.Fatal(r.Error)
	}

	logrus.Info(len(rooms))

	r = Database.Model(&Vup{}).
		Where("room_id in ?", rooms).
		Pluck("uid", &vups)

	if r.Error != nil {
		logrus.Fatal(r.Error)
	}

	logrus.Info(len(vups))
	logrus.Debugf("affected %v", r.RowsAffected)
}

func binit() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("Error while loading environment file: %v", err)
	}
	InitDB()
}
