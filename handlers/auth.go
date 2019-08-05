package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sunliang711/DKMS/types"

	//"github.com/gin-contrib/cors"
)

const (
	tokenExpired = "Token expired"
)

var (
	token2User = make(map[string]string)
)

func authMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("token")
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(types.Key), nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				c.JSON(200, gin.H{
					"code": 1,
					"msg":  tokenExpired,
				})
			default:
				c.JSON(200, gin.H{
					"code": 1,
					"msg":  "parse token error",
				})
			}
		default:
			c.JSON(200, gin.H{
				"code": 1,
				"msg":  "invalid token",
			})
		}
		c.Abort()
		return
	}
	token2User[t.Raw] = t.Claims.(jwt.MapClaims)["phone"].(string)
	logrus.Debugf("token2User: %v", token2User)
	c.Next()
}
