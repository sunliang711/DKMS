package models

import (
	"errors"
	"fmt"
	"time"

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
func AddUser(ui *types.RegisterObj) error {
	sql := "insert into `user` values(?,?,?,?,?,?,?,?,?);"
	// _, err := db.Exec(sql, ui.PID, ui.Username, ui.Password, ui.Phone, ui.AuthKey, ui.BeginTimestamp,ui.ExpiredTimestamp, ui.BeginTimestamp3rd,ui.ExpiredTimestamp3rd)
	_, err := db.Exec(sql, ui.PID, ui.Username, types.DefaultPassword, ui.Phone, ui.AuthKey, ui.BeginTimestamp, ui.ExpiredTimestamp, ui.BeginTimestamp3rd, ui.ExpiredTimestamp3rd)
	if err != nil {
		logrus.Errorf("insert into `user` error:", err)
		return err
	}
	return nil

}

// UpdateUserTimestamp TODO
// 只用到types.UserInfo里的PID Phone ExpiredTimestamp ExpiredTimestamp3rd
// 2019/08/08 17:27:37
func UpdateUserTimestamp(ui *types.UserInfo) error {
	exist, err := ExistUser(ui.PID, ui.Phone)
	if err != nil {
		msg := fmt.Sprintf("Error while check user: %v", err)
		logrus.Errorf(msg)
		return err
	}
	if exist {
		sql := "update `user` set beginTimestamp = ?,expiredTimestamp = ? ,beginTimestamp3rd = ? ,expiredTimestamp3rd = ? where pid = ? and phone = ?"
		now := time.Now().Unix()
		logrus.Debugf("update user timestamp with pid: %v phone: %v expired timestamp: %v expired timestamp 3rd: %v", ui.PID, ui.Phone, ui.ExpiredTimestamp, ui.ExpiredTimestamp3rd)
		_, err = db.Exec(sql, now, int64(ui.ExpiredTimestamp), now, int64(ui.ExpiredTimestamp3rd), ui.PID, ui.Phone)
		if err != nil {
			return err
		}
		return nil

	} else {
		msg := fmt.Sprintf("No such user with pid: %v phone: %v", ui.PID, ui.Phone)
		logrus.Error(msg)
		return err
	}
}

// GetUserTimeStamp TODO
// 2019/08/09 10:45:08
func GetUserTimestamp(ui *types.Identity) (int, int, error) {
	exist, err := ExistUser(ui.PID, ui.Phone)
	if err != nil {
		msg := fmt.Sprintf("Error while check user: %v", err)
		logrus.Errorf(msg)
		return 0, 0, err
	}
	if exist {
		sql := "select expiredTimeStamp,expiredTimeStamp3rd from `user` where pid = ? and phone = ?"
		rows, err := db.Query(sql, ui.PID, ui.Phone)
		if err != nil {
			return 0, 0, err
		}
		defer rows.Close()
		var (
			ts    int
			ts3rd int
		)
		if rows.Next() {
			rows.Scan(&ts, &ts3rd)
		}
		return ts, ts3rd, nil

	} else {
		msg := fmt.Sprintf("No such user with pid: %v phone: %v", ui.PID, ui.Phone)
		logrus.Error(msg)
		return 0, 0, err
	}
}

// getAllAdmins 获取所有admin的pid和phone
// 2019/08/13 17:12:00
func GetAllAdmins() ([]*types.PidPhone, error) {
	sql := "select pid,phone from `admin`;"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	var (
		pid   string
		phone string
		ret   []*types.PidPhone
	)
	for rows.Next() {
		rows.Scan(&pid, &phone)
		ret = append(ret, &types.PidPhone{pid, phone})
	}
	return ret, nil
}
