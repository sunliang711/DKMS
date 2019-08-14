package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/sunliang711/DKMS/models"
	"github.com/sunliang711/DKMS/types"
	"github.com/sunliang711/DKMS/utils"
)

func register(c *gin.Context) {
	var ui types.RegisterObj
	if err := c.BindJSON(&ui); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "register data format error",
		})
		return
	}
	exist, err := models.ExistAdmin(ui.PID, ui.Token)
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
		msg := fmt.Sprintf("No such pid or token : %v", ui.PID)
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

	msg := ""
	keyResults, err := keyGen()
	if err != nil {
		msg = fmt.Sprintf("Generate key failed: %v", err)
		logrus.Error(msg)
	} else {
		err = models.AddKeyPair(ui.PID, ui.Phone, keyResults.EccEncPK, keyResults.EccEncSK, types.KeyTypeMerchant)
		if err != nil {
			msg = fmt.Sprintf("Add merchant key pair error: %v", err)
			logrus.Error(msg)
		}
		err = models.AddKeyPair(ui.PID, ui.Phone, keyResults.VerifyingKey, keyResults.SigningKey, types.KeyType3rd)
		if err != nil {
			msg = fmt.Sprintf("Add 3rd key pair error: %v", err)
			logrus.Error(msg)
		}
	}

	if msg == "" {
		pkM, _, err := models.GetKeyPair(ui.PID, ui.Phone, types.KeyTypeMerchant)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "register ok,but get key pair failed",
			})
			logrus.Debugf("register ok,but get key pair failed")
			return
		}
		pk3rd, _, err := models.GetKeyPair(ui.PID, ui.Phone, types.KeyType3rd)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "register ok,but get key pair failed",
			})
			logrus.Debugf("register ok,but get key pair failed")
			return
		}
		//gen RK
		err = rkGen(ui.PID, ui.Phone)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 2,
				"msg":  "register ok,gen key pair ok,but gen rk failed",
			})
			logrus.Debugf("register ok,gen key pair ok,but gen rk failed")
			return
		}

		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "OK",
			"data": struct {
				PKMerchant string `json:"pk_merchant"`
				// SKMerchant string `json:"sk_merchant"`
				PK3rd string `json:"pk_3rd"`
				// SK3rd string `json:"sk_3rd"`
			}{
				pkM,
				// skM,
				pk3rd,
				// sk3rd,
			},
		})
		logrus.Info("Everything is OK")
	} else {
		c.JSON(200, gin.H{
			"code": 2,
			"msg":  "register ok,but key gen failed",
		})
		return
	}

}

type keyGenResult struct {
	EccEncPK     string `json:"ecc_enc_pk"`
	EccEncSK     string `json:"ecc_enc_sk"`
	SigningKey   string `json:"signing_key"`
	VerifyingKey string `json:"verifying_key"`
}

// KeyGen TODO
func keyGen() (*keyGenResult, error) {
	path := viper.GetString("enc.server") + "/keygen"
	param := struct {
		Type string `json:"type"`
		Role string `json:"role"`
	}{"ECC", "sender"}
	bs, _ := json.Marshal(&param)
	body := bytes.NewReader(bs)
	request, _ := http.NewRequest("POST", path, body)
	request.Header.Add("token", viper.GetString("enc.token"))
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		msg := fmt.Sprintf("Generate key pair error: %v", err)
		logrus.Error(msg)
		return nil, fmt.Errorf(msg)
	} else {
		var keyPairs keyGenResult
		err = json.NewDecoder(resp.Body).Decode(&keyPairs)
		if err != nil {
			msg := fmt.Sprintf("Decode key result error: %v", err)
			logrus.Error(msg)
			return nil, fmt.Errorf(msg)
		}
		return &keyPairs, nil
	}
}

// rkGen TODO
// 2019/08/13 14:54:07
func rkGen(pid, phone string) error {
	admins, err := models.GetAllAdmins()
	if err != nil {
		msg := fmt.Sprintf("Get all admins error: %v", err)
		logrus.Error(msg)
		return fmt.Errorf(msg)
	}
	//get encSK
	_, skM, err := models.GetKeyPair(pid, phone, types.KeyTypeMerchant)
	if err != nil {
		msg := fmt.Sprintf("Get merchant key pair error: %v", err)
		logrus.Error(msg)
		return fmt.Errorf(msg)
	}
	logrus.Debugf("sk: %v", skM)
	//get signSK
	_, sk3rd, err := models.GetKeyPair(pid, phone, types.KeyType3rd)
	if err != nil {
		return fmt.Errorf("Get 3rd key pair error: %v", err)
	}
	logrus.Debugf("sk3rd: %v", sk3rd)

	path := viper.GetString("enc.server") + "/kfraggen"
	logrus.Debugf("path: %v", path)
	client := &http.Client{}
	//get pk of every admin of admins
	for _, admin := range admins {
		pkM, _, err := models.GetKeyPair(admin.PID, admin.Phone, types.KeyTypeMerchant)
		if err != nil {
			msg := fmt.Sprintf("Get admin key pair of pid: %v phone: %v error: %v", admin.PID, admin.Phone, err)
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}
		t, n, err := models.Threshold(admin.PID)
		if err != nil {
			return err
		}
		param := struct {
			DelegateKey string `json:"delegatekey"`
			SignerSK    string `json:"signersk"`
			PublicKey   string `json:"publickey"`
			Threshold   int    `json:"threshold"`
			N           int    `json:"N"`
		}{skM, sk3rd, pkM, t, n}
		logrus.Debugf("param: %+v", param)
		bs, err := json.Marshal(&param)
		if err != nil {
			return fmt.Errorf("Marshal param eror: %v", err)
		}
		body := bytes.NewReader(bs)
		request, err := http.NewRequest("POST", path, body)
		if err != nil {
			return fmt.Errorf("NewRequest to path: %v error: %v", path, err)
		}
		request.Header.Add("token", viper.GetString("enc.token"))
		request.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(request)
		if err != nil {
			return fmt.Errorf("Post path: %v error: %v", path, err)
		}
		defer resp.Body.Close()
		//decode rk
		var rk []string
		rkBytes, _ := ioutil.ReadAll(resp.Body)
		logrus.Debugf("rkbytes: %v", string(rkBytes))
		err = json.NewDecoder(bytes.NewReader(rkBytes)).Decode(&rk)
		if err != nil {
			msg := fmt.Sprintf("decode rk error: %v", err)
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}

		//add rk to db
		err = models.AddRK(pid, phone, admin.PID, t, n, rk)
		if err != nil {
			msg := fmt.Sprintf("add rk error: %v", err)
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}
	}

	return nil
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
		t, err := utils.GenToken(cre.PID, cre.Phone)
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

// updateExpiredTimestamp TODO
// 2019/08/08 17:40:28
func updateExpiredTimestamp(c *gin.Context) {
	// models.UpdateUserTimestamp()
	var ui types.UserInfo
	if err := c.BindJSON(&ui); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "request data format error",
		})
		return
	}

	err := models.UpdateUserTimestamp(&ui)
	if err != nil {
		msg := fmt.Sprintf("update user timestamp error: %v", err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "OK",
	})
}

func getExpiredtimestamp(c *gin.Context) {
	var ui types.Identity
	if err := c.BindQuery(&ui); err != nil {
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  "request data format error",
		})
		return
	}

	ts, ts3rd, err := models.GetUserTimestamp(&ui)
	if err != nil {
		msg := fmt.Sprintf("get user timestamp error: %v", err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "OK",
		"data": struct {
			ExpiredTimestamp    int `json:"expired_timestamp"`
			ExpiredTimestamp3rd int `json:"expired_timestamp3rd"`
		}{ts, ts3rd},
	})
}
