package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Logger *log.Logger

func setConfig() {
	viper.AutomaticEnv()
	viper.SetTypeByDefaultValue(true)
	Logger = getLogger()
}

func getLogger() *log.Logger {
	lvl := viper.GetString("IPFS_LOG_LEVEL")
	level, err := log.ParseLevel(lvl)
	var logger = log.New()
	if err != nil {
		logger.Errorf("Failed to parse log level value `%s`. Setting to `%s`", lvl, log.InfoLevel.String())
	} else {
		logger.SetLevel(level)
	}
	logger.SetReportCaller(true)
	return logger
}
