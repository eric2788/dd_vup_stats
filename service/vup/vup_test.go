package vup

import (
	"encoding/json"
	"testing"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/stats"
	"vup_dd_stats/utils/set"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

func aTestGetVup(t *testing.T) {
	vup, err := GetVup(1042854135)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, vup)

	vup, err = GetVup(123456789)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(vup)
	}
}

func aTestSearchVups(t *testing.T) {
	vup, err := SearchVups("", 1, 5, "last_listened_at", true)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, vup)
}

func aTestDeleteNonVups(t *testing.T) {
	vtbs, err := stats.GetVtbListVtbMoe()
	if err != nil {
		t.Fatal(err)
	}
	vtbUids := maps.Keys(vtbs)
	result := db.Database.Delete(&db.Vup{}, "uid NOT IN ?", vtbUids)

	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Logf("已成功刪除 %v 列非虛擬主播。", result.RowsAffected)
}

func ainit() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("Error while loading environment file: %v", err)
	}
	info, err := stats.GetListening()
	if err != nil {
		logrus.Fatal(err)
	}
	stats.Listening = set.FromArray(info.Rooms)
	db.InitDB()
	db.InitRedis()
}

func jsonPrettyPrint(t *testing.T, v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	println(string(b))
}
