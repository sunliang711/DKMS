package main

import (
	"fmt"
	_ "github.com/sunliang711/DKMS/config"
	"github.com/sunliang711/DKMS/handlers"
	"github.com/sunliang711/DKMS/models"
	"github.com/spf13/viper"
)

func main() {
	dsn := viper.GetString("mysql.dsn")
	models.InitMysql(dsn)

	addr := fmt.Sprintf(":%d", viper.GetInt("server.port"))
	tls := viper.GetBool("tls.enable")
	certFile := viper.GetString("tls.certFile")
	keyFile := viper.GetString("tls.keyFile")
	handlers.StartServer(addr, tls, certFile, keyFile)
}
