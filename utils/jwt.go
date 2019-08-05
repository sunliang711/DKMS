package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sunliang711/DKMS/types"
	"time"
)

func GenToken(phone string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iss":   "dkms",
		"phone": phone,
	})
	tokenStr, err := t.SignedString([]byte(types.Key))
	if err != nil {
		return "", err
	}
	return tokenStr, nil

}
