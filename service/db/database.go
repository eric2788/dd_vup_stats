package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strings"
)

var (
	log      = logrus.WithField("service", "db")
	Database *gorm.DB
)

func getMySQLDataSource() string {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASS")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDB := os.Getenv("MYSQL_DB")

	return fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=true",
		mysqlUser,
		mysqlPass,
		mysqlHost,
		mysqlPort,
		mysqlDB,
	)
}

func getPgSQLDataSource() string {
	pgUser := os.Getenv("PG_USER")
	pgPass := os.Getenv("PG_PASS")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDB := os.Getenv("PG_DB")

	return fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		pgHost,
		pgPort,
		pgUser,
		pgDB,
		pgPass,
	)
}

func InitDB() {

	log.Info("正在連接資料庫...")

	dbType := strings.ToLower(os.Getenv("DB_TYPE"))

	var dataSource string

	switch dbType {
	case "mysql":
		dataSource = getMySQLDataSource()
	case "postgres":
		dataSource = getPgSQLDataSource()
	default:
		log.Fatalf("不支持的資料庫類型: %v", dbType)
	}

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
