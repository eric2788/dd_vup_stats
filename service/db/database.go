package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strings"
)

var (
	log      = logrus.WithField("service", "db")
	Database *gorm.DB
)

func InitDB() {

	log.Info("正在連接資料庫...")

	dbType := strings.ToLower(os.Getenv("DB_TYPE"))

	if dbType == "" {
		log.Fatal("未設定資料庫類型, 請在環境參數中設定 DB_TYPE")
	}

	getDataSource, exist := dialectMap[dbType]

	if !exist {
		log.Fatalf("不支持的資料庫類型: %v", dbType)
	}

	var logLevel logger.LogLevel

	if os.Getenv("GIN_MODE") != "release" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(getDataSource(), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("啟動資料庫時出現錯誤: %v", err)
	}

	log.Info("資料庫連接成功")

	if err = db.
		AutoMigrate(&Vup{}, &Behaviour{}, &LastListen{}); err != nil {
		log.Fatalf("Error while auto migrating tables: %v", err)
	}

	Database = db
}
