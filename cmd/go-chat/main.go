package main

import (
	"fmt"
	"github.com/vaberof/go-chat/pkg/http/httpserver"
	"github.com/vaberof/go-chat/pkg/logging/logs"
)

var defaultAppConfigPath = "./config/application.yaml"
var defaultLoggerConfigPath = "./config/logger.json"

func main() {
	logs, err := logs.New(defaultLoggerConfigPath)
	if err != nil {
		panic(err)
	}

	appConfig := getAppConfig(defaultAppConfigPath)

	fmt.Printf("%+v\n", appConfig)

	appServer := httpserver.New(&appConfig.HttpServer, logs)

	startChannel := appServer.StartAsync()

	<-(*startChannel)
}
