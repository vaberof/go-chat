package main

import (
	"errors"
	"github.com/vaberof/go-chat/internal/domain/websocket"
	"github.com/vaberof/go-chat/pkg/config"
	"github.com/vaberof/go-chat/pkg/http/httpserver"
)

type AppConfig struct {
	HttpServer      httpserver.HttpServerConfig
	WebSocketServer websocket.WebSocketServerConfig
}

func getAppConfig(sources ...string) AppConfig {
	config, err := tryGetAppConfig(sources...)
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic(errors.New("config cannot be nil"))
	}

	return *config
}

func tryGetAppConfig(sources ...string) (*AppConfig, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least 1 source must be set for app config")
	}

	provider := config.MergeConfigs(sources)

	var httpServerConfig httpserver.HttpServerConfig
	err := config.ParseConfig(provider, "app.http.server", &httpServerConfig)
	if err != nil {
		return nil, err
	}

	var websocketServerConfig websocket.WebSocketServerConfig
	err = config.ParseConfig(provider, "app.websocket.server", &websocketServerConfig)
	if err != nil {
		return nil, err
	}

	config := AppConfig{
		HttpServer:      httpServerConfig,
		WebSocketServer: websocketServerConfig,
	}

	return &config, nil
}
