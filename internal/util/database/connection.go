package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DataSourceFormat is "user:password@tcp(host:port)/dbname?parseTime=true".
const DataSourceFormat = "%s:%s@tcp(%s:%s)/%s?parseTime=true"

func Connect(dataSourceName string) *sql.DB {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	err = db.Ping()
	if err != nil {
		db.Close()
		panic(err)
	}

	return db
}
