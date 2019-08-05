package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var (
	DB *sql.DB
)

func InitMysql(dsn string) {
	var err error
	if len(dsn) == 0 {
		logrus.Fatal("Mysql DSN is empty.")
	}
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatalf("Open mysql error: %v", err)
	}
	err = DB.Ping()
	if err != nil {
		logrus.Fatalf("Ping mysql error: %v", err)
	}

	logrus.Infoln("Connected to mysql.")
	DB.SetMaxIdleConns(20)
	DB.SetMaxOpenConns(20)
}

func CloseMysql() {
	logrus.Infoln("Close mysql.")
	DB.Close()
}
