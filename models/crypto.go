package models

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/sunliang711/DKMS/types"
)

// GetCryptoConfig TOOD
// 2019/08/17 12:07:43
func GetCryptoConfig(pid string) (*types.CryptoConfig, error) {
	sql := "select * from `cryptoConfig` where pid = ?"
	rows, err := db.Query(sql, pid)
	if err != nil {
		msg := fmt.Sprintf("run sql : %v error: %v", sql, err)
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	defer rows.Close()
	var cc types.CryptoConfig
	if rows.Next() {
		rows.Scan(&cc.PID, &cc.EccType, &cc.EncAlgorithm, &cc.GenManner, &cc.T, &cc.N, &cc.API)
	}
	cc.PID = pid
	return &cc, nil
}

// UpdateAPI TODO
// 2019/08/17 12:13:14
func UpdateAPI(pid, api string) error {
	cc, err := GetCryptoConfig(pid)
	if err != nil {
		msg := fmt.Sprintf("Query api with pid: %v error: %v", api, err)
		log.Error(msg)
		return fmt.Errorf(msg)
	}
	if cc.PID == pid {
		sql := "update `cryptoConfig` set api = ? where pid = ?"
		_, err = db.Exec(sql, api, pid)
		if err != nil {
			msg := fmt.Sprintf("Run sql: %v error: %v", sql, err)
			log.Error(msg)
			return fmt.Errorf(msg)
		}
		return nil
	} else {
		msg := fmt.Sprintf("No config with pid: %v", pid)
		log.Error(msg)
		return fmt.Errorf(msg)
	}

}
