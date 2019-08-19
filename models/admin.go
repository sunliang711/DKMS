package models

import (
	"fmt"

	"github.com/sunliang711/DKMS/types"
)

//ExistAdmin check existence with pid in table `admin`
func ExistAdmin(pid string, token string) (exist bool, err error) {
	var count int
	sql := "select count(*) from `admin` where pid=? and token = ?"
	rows, err := db.Query(sql, pid, token)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&count)
	}
	if count > 0 {
		exist = true
	}
	return
}

// GetAdminKeyPair
// 2019/08/14 17:46:47
func GetAdminKeyPair(pid string) (pk string, sk string, err error) {
	sql := "select phone from `admin` where pid = ?"
	rows, err := db.Query(sql, pid)
	if err != nil {
		return pk, sk, fmt.Errorf(fmt.Sprintf("get admin phone by pid error: %v", err))
	}
	defer rows.Close()
	var phone string
	if rows.Next() {
		rows.Scan(&phone)
	} else {
		return pk, sk, fmt.Errorf(fmt.Sprintf("no such phone with pid: %v", pid))
	}

	sql = "select "
	pk, sk, err = GetKeyPair(pid, phone, types.KeyTypeMerchant)
	return
}
