package models

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunliang711/DKMS/types"
)

//ExistUser checks if user exists with phone
func ExistUser(pid string, phone string) (bool, error) {
	var count int
	sql := "select count(*) from `user` where phone = ? and pid = ?"
	rows, err := db.Query(sql, phone, pid)
	if err != nil {
		msg := fmt.Sprintf("Execute sql: %v error: %v", sql, err)
		logrus.Error(msg)
		return false, errors.New(msg)
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&count)
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

//CheckUser checks if user exists with phone and password
func CheckUser(pid, phone, password string) (valid bool, err error) {
	var count int
	sql := "select count(*) from `user` where pid = ? and phone = ? and password = ?"
	rows, err := db.Query(sql, pid, phone, password)
	if err != nil {
		msg := fmt.Sprintf("Ececut sql: %v error: %v", sql, err)
		logrus.Error(msg)
		err = errors.New(msg)
		return
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&count)
	}
	if count > 0 {
		valid = true
		return
	}
	valid = false
	return
}

//AddUser adds user into db for register
func AddUser(ui *types.UserInfo) error {
	sql := "insert into `user` values(?,?,?,?,?,?,?);"
	_, err := db.Exec(sql, ui.PID, ui.Username, ui.Password, ui.Phone, ui.AuthKey, ui.ExpiredTimestamp, ui.ExpiredTimestamp3rd)
	if err != nil {
		logrus.Errorf("insert into `user` error:", err)
		return err
	}
	return nil

}
