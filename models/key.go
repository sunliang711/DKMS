package models

import "fmt"
import "github.com/sirupsen/logrus"

// GetKeyPair TODO
func GetKeyPair(pid string, phone string, keyType int) (string, string, error) {
	sql := "select pk,sk from `key` where pid = ? and phone = ? and keyType = ?"
	rows, err := db.Query(sql, pid, phone, keyType)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	var (
		pk string
		sk string
	)
	if rows.Next() {
		err = rows.Scan(&pk, &sk)
		if err != nil {
			return "", "", err
		}
		return pk, sk, nil
	}

	return "", "", fmt.Errorf(fmt.Sprintf("No such key pair with pid: %v phone: %v keyType: %v", pid, phone, keyType))
}

// AddKeyPair TODO
func AddKeyPair(pid, phone, pk, sk string, keyType int) error {
	curPK, _, _ := GetKeyPair(pid, phone, keyType)
	if curPK != "" {
		return fmt.Errorf("AddKeyPair: already exists")
	}
	sql := "insert into `key` values(?,?,?,?,?)"
	_, err := db.Exec(sql, pid, phone, pk, sk, keyType)
	if err != nil {
		msg := fmt.Sprintf("Exec sql: %v error: %v", sql, err)
		logrus.Error(msg)
		return fmt.Errorf(msg)
	}
	logrus.Info("AddKeyPair successful.")
	return nil
}

// UpdateKeyPair TODO
func UpdateKeyPair(pid, phone, pk, sk string, keyType int) error {
	curPK, _, _ := GetKeyPair(pid, phone, keyType)
	if curPK == "" {
		return fmt.Errorf("UpdateKeyPair: no old key pair")
	}
	sql := "update `key` set pk = ?,sk = ? where pid = ? and phone = ? and keyType = ?"
	_, err := db.Exec(sql)
	if err != nil {
		msg := fmt.Sprintf("Exec sql: %v error: %v", sql, err)
		logrus.Error(msg)
		return fmt.Errorf(msg)
	}
	logrus.Info("UpdateKeyPair successful.")
	return nil
}

// GetKeyFrag TODO
func GetKeyFrag(pid, phone, receiver string) ([]string, error) {
	sql := "select n,segment from `keyfrag` where pid = ? and phone = ? and receiver = ?"
	rows, err := db.Query(sql, pid, phone, receiver)
	if err != nil {
		msg := fmt.Sprintf("Exec sql: %v error: %v", sql, err)
		logrus.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	defer rows.Close()

	var (
		n    int
		seg  string
		segs []string
	)
	for rows.Next() {
		rows.Scan(&n, &seg)
		segs = append(segs, seg)
	}
	if n != len(segs) {
		msg := fmt.Sprintf("Internal db error: n != segment count")
		logrus.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	return segs, nil
}

// AddKeyFrag TODO
func AddKeyFrag(pid, phone, receiver string, t, n int, segs []string) error {
	sql := "select count(*) from `keyfrag` where pid = ? and phone = ? and receiver = ?"
	rows, err := db.Query(sql, pid, phone, receiver)
	if err != nil {
		msg := fmt.Sprintf("Exec sql: %v error: %v", sql, err)
		logrus.Error(msg)
		return fmt.Errorf(msg)
	}
	var count int
	if rows.Next() {
		rows.Scan(&count)
	}
	rows.Close()
	if count > 0 {
		msg := fmt.Sprintf("Already exists segments with pid: %v,phone: %v,receiver: %v", pid, phone, receiver)
		logrus.Error(msg)
		return fmt.Errorf(msg)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	sql = "insert into `keyfrag` values(?,?,?,?,?,?)"
	for _, s := range segs {
		tx.Exec(sql, pid, phone, receiver, t, n, s)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// DeleteKeyFrag TODO
func DeleteKeyFrag(pid, phone, receiver string) error {
	sql := "delete from `keyfrag` where pid = ? and phone = ? and receiver = ?"
	_, err := db.Exec(sql, pid, phone, receiver)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAuthKey TODO
func UpdateAuthKey(pid, phone, authKey string) error {
	exist, err := ExistUser(pid, phone)
	if err != nil {
		return err
	}
	if !exist {
		msg := fmt.Sprintf("不存在此用户")
		logrus.Error(msg)
		return fmt.Errorf(msg)
	}
	sql := "update `user` set authKey = ? where pid = ? and phone = ?"
	if _, err := db.Exec(sql, authKey, pid, phone); err != nil {
		return err
	}
	return nil
}
