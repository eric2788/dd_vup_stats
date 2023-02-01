package db

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func TestParallel(t *testing.T) {
	var a []int

	go func() {
		<-time.After(time.Second * 2)
		_ = json.Unmarshal([]byte(`[1,2,3]`), &a)
	}()

	<-time.After(time.Second * 3)
	t.Log(a)
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	if err := godotenv.Load("./../../.env"); err != nil {
		logrus.Errorf("Error while loading environment file: %v", err)
	}
	if os.Getenv("DB_TYPE") == "" {
		return
	}
	InitDB()
}
