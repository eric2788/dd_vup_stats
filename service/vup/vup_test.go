package vup

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"testing"
	"vup_dd_stats/service/db"
)

func ATestGetVups(t *testing.T) {
	vups, err := GetVups(1, 3, true, "dd_count")
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, vups)
}

func ATestGetVup(t *testing.T) {
	vup, err := GetVup(1042854135)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, vup)
}

func aTestSearchVups(t *testing.T) {
	vup, err := SearchVups("Official", 1, 3, "dd_count", true)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, vup)
}

func ainit() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Fatalf("Error while loading environment file: %v", err)
	}
	db.InitDB()
}

func jsonPrettyPrint(t *testing.T, v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	println(string(b))
}
