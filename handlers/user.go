package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sunliang711/DKMS/models"
	"github.com/sunliang711/DKMS/types"
	"github.com/sunliang711/DKMS/utils"
	"net/http"
)
import "github.com/spf13/viper"

func register(c *gin.Context) {
	var ui types.UserInfo
	if err := c.BindJSON(&ui); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "register data format error",
		})
		return
	}
	exist, err := models.ExistAdmin(ui.PID)
	if err != nil {
		msg := fmt.Sprintf("Query pid: %v failed: %v", ui.PID, err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	if !exist {
		msg := fmt.Sprintf("No such pid: %v", ui.PID)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	//query in db
	exist, err = models.ExistUser(ui.PID, ui.Phone)
	if err != nil {
		logrus.Errorf("Internal db error: %v", err)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "inter db error",
		})
		return
	}
	if exist {
		msg := "该用户已注册!"
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		logrus.Error(msg)
		return
	}
	//insert into db
	err = models.AddUser(&ui)
	if err != nil {
		msg := fmt.Sprintf("Internal db error: %v", err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}

	keyResults,err := keyGen()
	if err != nil{
		msg := fmt.Sprintf("Generate key failed: %v",err)
		logrus.Error(msg)
	}else{
		err = models.AddKeyPair(ui.PID,ui.Phone,keyResults.EccEncPK,keyResults.EccEncSK,types.KeyTypeMerchant)
		if err != nil{
			logrus.Error("Add merchant key pair error: %v",err)
		}
		err = models.AddKeyPair(ui.PID,ui.Phone,keyResults.SigningKey,keyResults.VerifyingKey,types.KeyType3rd)
		if err != nil{
			logrus.Error("Add 3rd key pair error: %v",err)
		}
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "OK",
	})

}

type keyGenResult struct {
	EccEncPK     string `json:"ecc_enc_pk"`
	EccEncSK     string `json:"ecc_enc_sk"`
	SigningKey   string `json:"signing_key"`
	VerifyingKey string `json:"verifying_key"`
}

// KeyGen TODO
func keyGen() (*keyGenResult,error){
	//TODO create (sk,pk)(sk2,pk2)
	path := viper.GetString("enc.server") + "/keygen"
	param := struct {
		Type string `json:"type"`
		Role string `json:"role"`
	}{"ECC", "sender"}
	bs, _ := json.Marshal(&param)
	body := bytes.NewReader(bs)
	request ,_ := http.NewRequest("POST",path,body)
	request.Header.Add("token",viper.GetString("enc.token"))
	request.Header.Add("Content-Type","application/json")
	client := &http.Client{}
	resp,err := client.Do(request)
	if err != nil {
		msg := fmt.Sprintf("Generate key pair error: %v", err)
		logrus.Error(msg)
		return nil,fmt.Errorf(msg)
	} else {
		var keyPairs keyGenResult
		err = json.NewDecoder(resp.Body).Decode(&keyPairs)
		if err != nil{
			msg := fmt.Sprintf("Decode key result error: %v",err)
			logrus.Error(msg)
			return nil,fmt.Errorf(msg)
		}
		return &keyPairs,nil
	}
}
func login(c *gin.Context) {
	var cre types.Credential
	if err := c.BindJSON(&cre); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "login request data is not json",
		})
		return
	}

	//query in db to find user
	valid, err := models.CheckUser(cre.PID, cre.Phone, cre.Password)
	if err != nil {
		msg := fmt.Sprintf("Internal db error: %v", err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}

	if valid {
		t, err := utils.GenToken(cre.Phone)
		if err != nil {
			msg := fmt.Sprintf("Internal error while gen token: %v", err)
			logrus.Error(msg)
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  msg,
			})
			return
		}
		c.JSON(200, gin.H{
			"code":  0,
			"msg":   "OK",
			"token": t,
		})
	} else {
		msg := fmt.Sprintf("手机号密码不匹配")
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
	}
}
