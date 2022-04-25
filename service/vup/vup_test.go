package vup

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"testing"
	"vup_dd_stats/service/db"
	"vup_dd_stats/service/statistics"
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
	vtbs, err := statistics.GetVtbListVtbMoe()
	if err != nil {
		t.Fatal(err)
	}
	vtbUids := make([]int64, len(vtbs))
	for i, vtb := range vtbs {
		vtbUids[i] = vtb.Mid
	}
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
	info, err := statistics.GetListening()
	if err != nil {
		logrus.Fatal(err)
	}
	statistics.Listening = &info.Rooms
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
