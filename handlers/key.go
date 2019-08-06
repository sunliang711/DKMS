package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sunliang711/DKMS/models"
	"github.com/sunliang711/DKMS/types"
)

// GetPkSkPairs TODO
func GetPkSkPairs(c *gin.Context) {
	var id types.Identity
	if err := c.BindQuery(&id); err != nil {
		msg := fmt.Sprintf("Request data format error: %v", err)
		logrus.Errorf(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	if id.PID == "" || id.Phone == "" {
		msg := fmt.Sprintf("pid or phone is empty")
		logrus.Errorf(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	pkMer, skMer, err := models.GetKeyPair(id.PID, id.Phone, types.KeyTypeMerchant)
	if err != nil {
		msg := fmt.Sprintf("Get merchant key pair error: %v", err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	pk3rd, sk3rd, err := models.GetKeyPair(id.PID, id.Phone, types.KeyType3rd)
	if err != nil {
		msg := fmt.Sprintf("Get 3rd key pair error: %v", err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}

	keyPairs := struct {
		KeyMerchant types.KeyPair `json:"key_merchant"`
		Key3rd      types.KeyPair `json:"key_3rd"`
	}{KeyMerchant: types.KeyPair{pkMer, skMer},
		Key3rd: types.KeyPair{pk3rd, sk3rd},
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": keyPairs,
	})
}

// updateAuthKey TODO
func updateAuthKey(c *gin.Context) {
	var aKey types.AuthKey
	if err := c.BindJSON(&aKey); err != nil {
		msg := fmt.Sprintf("Request data format error: %v", err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	if err := models.UpdateAuthKey(aKey.PID, aKey.Phone, aKey.AuthKey); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "OK",
	})

}
