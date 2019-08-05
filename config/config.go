package config

import "github.com/spf13/viper"
import "github.com/spf13/pflag"
import "github.com/sirupsen/logrus"
import "github.com/sunliang711/DKMS/utils"

func init() {
	configPath := pflag.StringP("config", "c", "config.toml", "config file")
	pflag.Parse()
	logrus.Info("Config file: ",*configPath)

	viper.SetConfigFile(*configPath)
	viper.ReadInConfig()

	loglevel := viper.GetString("log.level")
	logrus.Infoln("Log level: ", loglevel)
	logrus.SetLevel(utils.LogLevel(loglevel))
}
