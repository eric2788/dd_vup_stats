package db

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	log      = logrus.WithField("service", "db")
	Database *gorm.DB
)

const (
	CountStatement = "SELECT cast(reltuples as bigint) AS count FROM pg_class where relname = ?"
)

func InitDB() {

	log.Info("正在連接資料庫...")

	var logLevel logger.LogLevel

	if os.Getenv("GIN_MODE") != "release" {
		logLevel = logger.Warn
	} else {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(getPgSQLDataSource(), &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		log.Fatalf("啟動資料庫時出現錯誤: %v", err)
	}

	pool, err := db.DB()
	if err == nil {
		pool.SetMaxIdleConns(5)
		pool.SetMaxOpenConns(500)
		pool.SetConnMaxLifetime(time.Minute * 15)
		pool.SetConnMaxIdleTime(time.Minute * 2)
	} else {
		log.Warnf("設定資料庫連接池時出現錯誤: %v", err)
	}

	log.Info("資料庫連接成功")

	if err = db.
		AutoMigrate(&Vup{}, &Behaviour{}, &LastListen{}, &UserAnalysis{}, &SearchAnalysis{}, &WatcherBehaviour{}); err != nil {
		log.Errorf("Error while auto migrating tables: %v", err)
	}

	go createMaterializedViews(db)

	Database = db
}
