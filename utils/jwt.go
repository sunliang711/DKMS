package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sunliang711/DKMS/types"
)
import "github.com/spf13/viper"

func GenToken(PID, phone string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour * time.Duration(viper.GetInt("server.jwtExp"))).Unix(),
		"iss":   "dkms",
		"phone": phone,
		"pid":   PID,
	})
	tokenStr, err := t.SignedString([]byte(types.Key))
	if err != nil {
		return "", err
	}
	return tokenStr, nil

}
