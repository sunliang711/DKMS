package models

import "github.com/sirupsen/logrus"
import "fmt"

//GetExpired gets expired timestamp with phone and pid
func GetExpired(pid string, phone string) (int, int, error) {
	sql := "select last,last3rd from `expired` where phone = ? and pid = ?"
	rows, err := db.Query(sql, phone, pid)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	var (
		expire    int
		expire3rd int
	)
	if rows.Next() {
		err = rows.Scan(&expire, &expire3rd)
		if err != nil {
			return 0, 0, err
		}
		return expire, expire3rd, nil
	}
	return 0, 0, fmt.Errorf("No such expired timestamp with phone: %v pid: %v", phone, pid)
}

//UpdateExpired update or insert current timestamp to expired table
func UpdateExpired(pid string, phone string, last int, last3rd int) error {
	exist, err := ExistExpired(pid, phone)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	var sql string
	if exist {
		sql = "update `expired` set last = ? ,last3rd = ? where phone = ? and pid = ?"
		_, err = db.Exec(sql, last, last3rd, phone, pid)
		if err != nil {
			msg := fmt.Sprintf("Update expired failed: %v", err)
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}
		logrus.Debugf("Update expired with phone: %v pid: %v successful", phone, pid)
	} else {
		sql = "insert into `expired` values(?,?,?,?);"
		_, err = db.Exec(sql, pid, phone, last, last3rd)
		if err != nil {
			msg := fmt.Sprintf("Insert into expired failed: %v", err)
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}
		logrus.Debugf("Insert to expired with phone: %v pid: %v successful", phone, pid)
	}
	return nil
}

//ExistExpired checks existence of phone and pid
func ExistExpired(pid string, phone string) (bool, error) {
	var count int
	sql := "select count(*) from `expired` where phone = ? and pid = ?"
	rows, err := db.Query(sql, phone, pid)
	if err != nil {
		msg := fmt.Sprintf("Execute sql: %v error: %v", sql, err)
		logrus.Error(msg)
		return false, fmt.Errorf(msg)
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
