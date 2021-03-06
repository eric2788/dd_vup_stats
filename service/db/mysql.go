package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var (
	log      = logrus.WithField("service", "db")
	Database *gorm.DB
)

func InitDB() {

	log.Info("正在連接資料庫...")

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASS")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDB := os.Getenv("MYSQL_DB")

	dataSource := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=true",
		mysqlUser,
		mysqlPass,
		mysqlHost,
		mysqlPort,
		mysqlDB,
	)

	var logLevel logger.LogLevel

	if os.Getenv("GIN_MODE") != "release" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("啟動資料庫時出現錯誤: %v", err)
	}

	log.Info("資料庫連接成功")

	if err = db.
		Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(&Vup{}, &Behaviour{}, &LastListen{}); err != nil {
		log.Fatalf("Error while auto migrating tables: %v", err)
	}

	Database = db
}
