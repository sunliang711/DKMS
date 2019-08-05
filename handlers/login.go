package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sunliang711/DKMS/models"
	"github.com/sunliang711/DKMS/types"
	"github.com/sunliang711/DKMS/utils"
)

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
		msg := "Already exist such user"
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
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "OK",
	})
	//TODO create (sk,pk)(sk2,pk2)

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
	valid, err := models.CheckUser(cre.PID,cre.Phone, cre.Password)
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
		msg := fmt.Sprintf("Phone doesn't match password.")
		logrus.Error(msg)
		c.JSON(200, gin.H{
			"code": 1,
			"msg":  msg,
		})
	}
}
