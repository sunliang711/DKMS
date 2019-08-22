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
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "token")
	router.Use(cors.New(config))

	//admin register its user
	router.POST("/register", register)
	//user login
	router.POST("/login", login)
	router.POST("/pre", PRE)
	router.POST("/check_pop", CheckPop)
	router.POST("/auth_claim", AuthClaim)
	router.POST("/check_auth", CheckAuth)

	router.GET("/get_crypto_config", getCryptoConfig)
	router.POST("/update_api", updateAPI)
	router.GET("/get_expired_timestamp2", getExpiredtimestamp)

	//the following need authMiddleware
	router.Use(authMiddleware)
	router.GET("/get_pk_sk_pairs", GetPkSkPairs)
	router.POST("/update_auth_key", updateAuthKey)
	router.POST("/update_expired_timestamp", updateExpiredTimestamp)
	router.GET("/get_expired_timestamp", getExpiredtimestamp)

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
