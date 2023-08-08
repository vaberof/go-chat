package main

import (
	"errors"
	"github.com/vaberof/go-chat/internal/domain/chat/auth"
	"github.com/vaberof/go-chat/internal/domain/websocket/websocketserver"
	"github.com/vaberof/go-chat/internal/infra/storage/postgres"
	"github.com/vaberof/go-chat/pkg/config"
	"github.com/vaberof/go-chat/pkg/http/httpserver"
)

type AppConfig struct {
	HttpServer       httpserver.HttpServerConfig
	WebsocketServer  websocketserver.WebsocketServerConfig
	AuthService      auth.AuthServiceConfig
	PostgresDatabase postgres.PostgresDatabaseConfig
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

	var websocketServerConfig websocketserver.WebsocketServerConfig
	err = config.ParseConfig(provider, "app.websocket.server", &websocketServerConfig)
	if err != nil {
		return nil, err
	}

	var authServiceConfig auth.AuthServiceConfig
	err = config.ParseConfig(provider, "app.auth", &authServiceConfig)
	if err != nil {
		return nil, err
	}

	var postgresConfig postgres.PostgresDatabaseConfig
	err = config.ParseConfig(provider, "app.database.postgres", &postgresConfig)
	if err != nil {
		return nil, err
	}

	config := AppConfig{
		HttpServer:       httpServerConfig,
		WebsocketServer:  websocketServerConfig,
		AuthService:      authServiceConfig,
		PostgresDatabase: postgresConfig,
	}

	return &config, nil
}
