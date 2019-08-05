package models

import (
	"database/sql"
	// register mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var (
	db *sql.DB
)

//InitMysql open mysql with dsn
func InitMysql(dsn string) {
	var err error
	if len(dsn) == 0 {
		logrus.Fatal("Mysql DSN is empty.")
	}
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatalf("Open mysql error: %v", err)
	}
	err = db.Ping()
	if err != nil {
		logrus.Fatalf("Ping mysql error: %v", err)
	}

	logrus.Infoln("Connected to mysql.")
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
}

//CloseMysql close mysql connection
func CloseMysql() {
	logrus.Infoln("Close mysql.")
	db.Close()
}
