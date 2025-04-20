package main

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetRpcEndpointURL(endpoint string) string {
	host := viper.GetString("IPFS_HOST")
	port := viper.GetString("IPFS_RPC_API_PORT")
	endpoint, _ = strings.CutPrefix(endpoint, "/")
	return fmt.Sprintf("http://%s:%s/%s", host, port, endpoint)
}

func debugEntry(r any) {
	if Logger.Level == log.DebugLevel {
		data, err := json.Marshal(&r)
		if err != nil {
			Logger.Debugf("Could not unmarshal")
			return
		}
		Logger.Debugf("returned data \n%s", string(data))
	}
}
