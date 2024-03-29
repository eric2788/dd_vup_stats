package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// getMySQLDataSource returns a gorm.Dialector for MySQL.
// Deprecated: Use getPgSQLDataSource instead.
/*
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
*/

// getPgSQLDataSource returns a gorm.Dialector for PostgreSQL.
func getPgSQLDataSource() gorm.Dialector {
	pgUser := os.Getenv("PG_USER")
	pgPass := os.Getenv("PG_PASS")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDB := os.Getenv("PG_DB")
	pgSSL := os.Getenv("PG_SSL")

	if pgSSL == "" {
		pgSSL = "disable"
	}

	return postgres.Open(fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		pgHost,
		pgPort,
		pgUser,
		pgDB,
		pgPass,
		pgSSL,
	))
}
