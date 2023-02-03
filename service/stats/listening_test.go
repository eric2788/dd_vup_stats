package stats

import (
	"testing"
	"vup_dd_stats/service/db"
	"vup_dd_stats/utils/set"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func aTestFetchVupToRedis(t *testing.T) {
	fetchVupListToRedis()
}

func aTestRemoveUnsedVups(t *testing.T) {
	removeUnusedVupListFromRedis()
}

func aTestFetchListeningInfo(t *testing.T) {
	fetchListeningInfo()
}

func ainit() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("Error while loading environment file: %v", err)
	}
	info, err := GetListening()
	if err != nil {
		logrus.Fatal(err)
	}
	Listening = set.FromArray(info.Rooms)
	db.InitDB()
	db.InitRedis()
}
