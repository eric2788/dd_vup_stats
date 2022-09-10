package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var (
	dialectMap = map[string]func() gorm.Dialector{
		"mysql":    getMySQLDataSource,
		"postgres": getPgSQLDataSource,
	}
)

func getMySQLDataSource() gorm.Dialector {
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASS")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDB := os.Getenv("MYSQL_DB")

	return mysql.Open(fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=true",
		mysqlUser,
		mysqlPass,
		mysqlHost,
		mysqlPort,
		mysqlDB,
	))
}

func getPgSQLDataSource() gorm.Dialector {
	pgUser := os.Getenv("PG_USER")
	pgPass := os.Getenv("PG_PASS")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDB := os.Getenv("PG_DB")

	return postgres.Open(fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		pgHost,
		pgPort,
		pgUser,
		pgDB,
		pgPass,
	))
}
