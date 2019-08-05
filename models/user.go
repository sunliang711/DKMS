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

//UpdateExpired update or insert current timestamp to expired table
func UpdateExpired(pid string, phone string,last int) error {
	exist, err := ExistExpired(pid, phone)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	var sql string
	if exist {
		sql = "update `expired` set last = ? where phone = ? and pid = ?"
		_, err = db.Exec(sql, last, phone, pid)
		if err != nil {
			msg := fmt.Sprintf("Update expired failed: %v", err)
			logrus.Error(msg)
			return errors.New(msg)
		}
		logrus.Debugf("Update expired with phone: %v pid: %v successful", phone, pid)
	} else {
		sql = "insert into `expired` values(?,?,?);"
		_, err = db.Exec(sql, pid, phone, last)
		if err != nil {
			msg := fmt.Sprintf("Insert into expired failed: %v", err)
			logrus.Error(msg)
			return errors.New(msg)
		}
		logrus.Debugf("Insert to expired with phone: %v pid: %v successful", phone, pid)
	}
	return nil
}

//ExistExpired checks existence of a phone
func ExistExpired(pid string, phone string) (bool, error) {
	var count int
	sql := "select count(*) from `expired` where phone = ? and pid = ?"
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

//GetExpired gets expired timestamp with phone
func GetExpired(pid string, phone string) (int, error) {
	sql := "select last from `expired` where phone = ? and pid = ?"
	rows, err := db.Query(sql, phone, pid)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var expire int
	if rows.Next() {
		err = rows.Scan(&expire)
		if err != nil {
			return 0, err
		}
		return expire, nil
	}
	return 0, fmt.Errorf("No such expired timestamp with phone: %v pid: %v", phone, pid)
}
