package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// StartServer starts gin server
func StartServer(addr string, tls bool, certFile string, keyFile string) {
	//MUST SetMode first
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(cors.Default())

	//admin register its user
	router.POST("/register", register)
	//user login
	router.POST("/login", login)

	//the following need authMiddleware
	router.Use(authMiddleware)
	router.GET("/get_pk_sk_pairs", GetPkSkPairs)
	router.POST("/update_auth_key", updateAuthKey)

	logrus.Infof("Start server on %v, tls enabled: %v", addr, tls)
	if tls {
		router.RunTLS(addr, certFile, keyFile)
	} else {
		router.Run(addr)
	}

}

func todo(c *gin.Context) {
	t := c.Request.Header.Get("token")
	p := token2User[t]
	logrus.Debugf("phone: %v", p)
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "OK",
	})
}
