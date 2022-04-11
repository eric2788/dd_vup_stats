package vup

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"testing"
	"vup_dd_stats/service/db"
)

func TestGetVups(t *testing.T) {
	vups, err := GetVups(1, 3, true)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(vups, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("Error while loading environment file: %v", err)
	}
	db.InitDB()
}
