package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
		msg := fmt.Sprintf("Query with pid: %v token: %v failed: %v", ui.PID, ui.Token, err)
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
		return
	}
	if !exist {
		msg := fmt.Sprintf("No such pid:%v or token:%v", ui.PID, ui.Token)
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
	if err != nil {
		msg := fmt.Sprintf("Generate key pair error: %v", err)
		logrus.Error(msg)
		return nil, fmt.Errorf(msg)
	} else {
		defer resp.Body.Close()
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
		logrus.Debugf("kfragen param: %v", string(bs))
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

	_, ts, _, ts3rd, err := models.GetUserTimestamp(&ui)
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

// PRE 代理重加密
// 2019/08/14 17:08:43
func PRE(c *gin.Context) {
	var input types.PREInput
	if err := c.BindJSON(&input); err != nil {
		msg := fmt.Sprintf("bad request format: %v", err)
		c.JSON(400, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	exist, err := models.ExistAdmin(input.PID, input.Token)
	if err != nil {
		msg := fmt.Sprintf("check admin existence error: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}
	if !exist {
		msg := fmt.Sprintf("No such admin with pid: %v,token: %v", input.PID, input.Token)
		c.JSON(400, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	exist, err = models.ExistUser(input.PID, input.Phone)
	if err != nil {
		msg := fmt.Sprintf("check user existence error: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}
	if !exist {
		msg := fmt.Sprintf("No such user with pid: %v,phone: %v", input.PID, input.Phone)
		c.JSON(400, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	//get delegating加密公钥
	delegating, _, err := models.GetKeyPair(input.PID, input.Phone, types.KeyTypeMerchant)
	if err != nil {
		msg := fmt.Sprintf("get delegating error: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	//get receiving  接收者公钥
	receiving, _, err := models.GetAdminKeyPair(input.Receiver)
	if err != nil {
		msg := fmt.Sprintf("get receiving error: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	//get verifying 验签公钥
	verifying, _, err := models.GetKeyPair(input.PID, input.Phone, types.KeyType3rd)
	if err != nil {
		msg := fmt.Sprintf("get verifying error: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	logrus.Debugf("delegating: %v", delegating)
	logrus.Debugf("verifying: %v", verifying)
	logrus.Debugf("receiving: %v", receiving)

	//get rk by pid,phone,receiver
	rks, t, err := models.GetRK(input.PID, input.Phone, input.Receiver)
	if err != nil {
		msg := fmt.Sprintf("Get RK error: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	//call remote PRE api
	path := viper.GetString("enc.server") + "/reencrypt"
	client := http.Client{}
	preReq := struct {
		Addresses  []string `json:"addresses"`
		Capsule    string   `json:"capsule"`
		Delegating string   `json:"delegating"`
		Receiving  string   `json:"receiving"`
		Verifying  string   `json:"verifying"`
		Threshold  int      `json:"threshold"`
	}{rks[:t], input.Capsule, delegating, receiving, verifying, t}

	bs, _ := json.Marshal(&preReq)
	logrus.Infof("pre request data: %+v", string(bs))
	body := bytes.NewReader(bs)
	request, _ := http.NewRequest("POST", path, body)
	request.Header.Add("token", viper.GetString("enc.token"))
	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		msg := fmt.Sprintf("call remote pre api error: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}
	defer resp.Body.Close()
	resBs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("read remote pre response error: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
		})
		logrus.Error(msg)
		return
	}
	logrus.Debugf("response from pre api server: %v", string(resBs))

	preRespon := struct {
		Caddrs  []string `json:"caddrs"`
		Capsule string   `json:"capsule"`
	}{}
	err = json.NewDecoder(bytes.NewReader(resBs)).Decode(&preRespon)
	if err != nil {
		msg := fmt.Sprintf("pre api server response not valid: %v", err)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
		})
		logrus.Error(msg)
		return
	}

	c.JSON(200, preRespon)
}

// CheckPop 是否需要弹框
// 2019/08/15 11:09:38
func CheckPop(c *gin.Context) {
	var chk types.CheckPopInput
	if err := c.BindJSON(&chk); err != nil {
		msg := fmt.Sprintf("bad request format")
		c.JSON(400, types.Response{
			Code: 2,
			Msg:  msg,
		})
		logrus.Error(msg)
		return
	}

	exist, err := models.ExistAdmin(chk.PID, chk.Token)
	if err != nil {
		msg := fmt.Sprintf("check admin existence error: %v", err)
		c.JSON(500, types.Response{
			Code: 2,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}
	if !exist {
		msg := fmt.Sprintf("No such admin with pid: %v,token: %v", chk.PID, chk.Token)
		c.JSON(400, types.Response{
			Code: 2,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	exist, err = models.ExistUser(chk.PID, chk.Phone)
	if err != nil {
		msg := fmt.Sprintf("check user existence error: %v", err)
		c.JSON(500, types.Response{
			Code: 2,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}
	if !exist {
		msg := fmt.Sprintf("No such user with pid: %v,phone: %v", chk.PID, chk.Phone)
		c.JSON(400, types.Response{
			Code: 2,
			Msg:  msg,
			Data: nil,
		})
		logrus.Error(msg)
		return
	}

	begin, expire, _, _, err := models.GetUserTimestamp(&types.Identity{
		PID:   chk.PID,
		Phone: chk.Phone,
	})
	if err != nil {
		msg := fmt.Sprintf("Get user timestamp error with pid: %v,phone: %v", chk.PidPhone, chk.Phone)
		c.JSON(500, types.Response{
			Code: 2,
			Msg:  msg,
		})
		logrus.Error(msg)
		return
	}
	if time.Now().Unix() > int64(begin+expire) {
		c.JSON(200, types.Response{
			Code: 1,
			Msg:  "need pop",
		})
		logrus.Debugf("need pop")
	} else {
		c.JSON(200, types.Response{
			Code: 0,
			Msg:  "not need pop",
		})
		logrus.Debugf("not need pop")
	}

}

// CheckAuth TODO
// 2019/08/15 18:19:45
func CheckAuth(c *gin.Context) {
	var input types.CheckAuthInput
	if err := c.BindJSON(&input); err != nil {
		msg := fmt.Sprintf("bad request format")
		c.JSON(400, types.Response{
			Code: 2,
			Msg:  msg,
		})
		logrus.Error(msg)
		return
	}
	result, _ := models.CheckAuthKey(input.PID, input.Phone, input.AuthKey, input.Period)
	logrus.Debugf("pid: %v phone: %v authKey: %v period: %v", input.PID, input.Phone, input.AuthKey, input.Period)
	if result {
		c.JSON(200, types.Response{
			Code: 0,
			Msg:  "验证通过",
		})
		logrus.Infof("验证通过")
	} else {
		c.JSON(200, types.Response{
			Code: 1,
			Msg:  "验证失败",
		})
		logrus.Infof("验证失败")
	}
}

// AuthClaim 理赔申请授权
// 2019/08/15 19:58:07
func AuthClaim(c *gin.Context) {
	var input types.AuthClaimInput
	if err := c.BindJSON(&input); err != nil {
		msg := fmt.Sprintf("bad reqeust format")
		c.JSON(200, types.Response{
			Code: 1,
			Msg:  msg,
		})
		logrus.Error(msg)
		return
	}

	//check authKey
	logrus.Debugf("check auth key with pid: %v phone: %v, authKey: %v", input.PID, input.Phone, input.AuthKey)
	pass, err := models.CheckAuthKeyOnly(input.PID, input.Phone, input.AuthKey)
	if err == nil && pass {
		// update period to db
		err := models.UpdatePeriod(input.PID, input.Phone, int(time.Now().Unix()), input.Period)
		if err != nil {
			msg := fmt.Sprintf("update period error: %v", err)
			c.JSON(200, types.Response{
				Code: 1,
				Msg:  msg,
			})
			logrus.Error(msg)
		}

		path := viper.GetString("ys.server") + "/authlist"
		data := types.AuthList{
			PidPhone: types.PidPhone{
				PID:   input.PID,
				Phone: input.Phone,
			},
			DeviceList: input.DeviceList,
			Start:      input.Start,
			End:        input.End,
			Period:     input.Period,
			Receiver:   input.Recevier,
		}
		logrus.Debugf("post to path: %v with data: %+v", path, data)
		resp, err := utils.Post(path, map[string]string{"Content-Type": "application/json"}, data)
		if err != nil {
			msg := fmt.Sprintf("call ys server /authlist api error: %v", err)
			logrus.Error(msg)
			c.JSON(500, types.Response{
				Code: 1,
				Msg:  msg,
			})
		}
		defer resp.Body.Close()
		var r types.Response
		err = json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			msg := fmt.Sprintf("decode /authlist api error: %v", err)
			c.JSON(500, types.Response{
				Code: 1,
				Msg:  msg,
			})
			logrus.Error(msg)
			return
		}
		if r.Code == 0 {
			c.JSON(200, types.Response{
				Code: 0,
				Msg:  "OK",
			})
			return
		} else {
			c.JSON(500, types.Response{
				Code: r.Code,
				Msg:  r.Msg,
				Data: r.Data,
			})
			return
		}
		//data":[{"capsule":"QmRuNX95x4ZXqpCcTtJfoukW3jyFo6yLox62wBiZCPYc6t","ciphertext":"QmUid62eMk1Y6oDEYfXo7oiQjTtzv17J889qjhDHnBmaMG","filename":"1565875831_0_sun.png","decrypt_mode":"1","deckey":""},{"capsule":"QmRj3vPidHbXwevBHRecEFuHQYRQozL3D9bHdsX4qc5Q6X","ciphertext":"QmbPmtrwiM487zQyRaVPUwCoAKjJRmax1F5Tyk5d6eqGv1","filename":"1565882396_0_sun.png","decrypt_mode":"1","deckey":""},{"capsule":"QmaEwAfaVgn5Dxq8b1gdpTHeXtG4wjDHhS9mpSqhxjBSbZ","ciphertext":"QmXN6wrTYbjESafbgLqLgARaUoyFE6kqd7yau25CJAmqkj","filename":"1565882400_0_sun.png","decrypt_mode":"1","deckey":""},{"capsule":"QmdMzJT1U8ReKDyZHkVJKA8S87dWjnnc5pxqVnJR3i6D7C","ciphertext":"QmbpqXvphcEsHG979J1g6gtR9jEJdTdiULfLxv1P5CWJ6Q","filename":"1565882403_0_sun.png","decrypt_mode":"1","deckey":""},{"capsule":"QmNmL3iuQLjkkbVyyn3CpShJCmBFxsNkJ9oWpuYn6YAfZA","ciphertext":"QmYSCEdPgkovRbEXjo5X7iQmPLDsJU4iCGRj2WYVLsbe4V","filename":"1565882407_0_sun.png","decrypt_mode":"1","deckey":""}]}
	} else {
		msg := fmt.Sprintf("check auth key failed")
		c.JSON(200, types.Response{
			Code: 1,
			Msg:  msg,
		})
		logrus.Error(msg)
		return
	}
}

// getCryptoConfig TODO
// 2019/08/17 12:17:32
func getCryptoConfig(c *gin.Context) {
	pid := c.Query("pid")
	if pid == "" {
		msg := fmt.Sprintf("Need pid parameter.")
		logrus.Error(msg)
		c.JSON(400, types.Response{
			Code: 1,
			Msg:  msg,
		})
		return
	}
	cc, err := models.GetCryptoConfig(pid)
	if err != nil {
		msg := fmt.Sprintf("Get crypto config error: %v", err)
		logrus.Error(msg)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
		})
		return
	}
	c.JSON(200, types.Response{
		Code: 1,
		Msg:  "OK",
		Data: cc,
	})

}

// updateAPI TODO
// 2019/08/17 12:21:45
func updateAPI(c *gin.Context) {
	//Only use pid and api fileds
	var cc types.CryptoConfig
	if err := c.BindJSON(&cc); err != nil {
		msg := fmt.Sprintf("Bad request format")
		logrus.Error(msg)
		c.JSON(400, types.Response{
			Code: 1,
			Msg:  msg,
		})
		return
	}

	if err := models.UpdateAPI(cc.PID, cc.API); err != nil {
		msg := fmt.Sprintf("Update api error: %v", err)
		logrus.Error(msg)
		c.JSON(500, types.Response{
			Code: 1,
			Msg:  msg,
		})
		return
	} else {
		c.JSON(200, types.Response{
			Code: 0,
			Msg:  "OK",
		})
	}
}
