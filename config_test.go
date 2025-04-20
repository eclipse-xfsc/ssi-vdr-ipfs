package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"testing"
)

func TestGetLogger(t *testing.T) {
	viper.Set("LOG_LEVEL", log.DebugLevel.String())
	logger := getLogger()
	if logger.Level != log.DebugLevel {
		t.Errorf("Logger not configured with expected level: %s actual level: %s", log.DebugLevel.String(), logger.Level.String())
	}
}
