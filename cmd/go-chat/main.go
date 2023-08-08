package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	httproutes "github.com/vaberof/go-chat/internal/app/entrypoint/http"
	"github.com/vaberof/go-chat/internal/app/entrypoint/websocket"
	"github.com/vaberof/go-chat/internal/domain/chat/auth"
	"github.com/vaberof/go-chat/internal/domain/chat/user"
	"github.com/vaberof/go-chat/internal/domain/websocket/websocketserver"
	"github.com/vaberof/go-chat/internal/infra/storage/postgres"
	"github.com/vaberof/go-chat/internal/infra/storage/postgres/postgresuser"
	"github.com/vaberof/go-chat/pkg/http/httpserver"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"os"
)

var appConfigPaths = flag.String("app.config.files", "not-found.yaml", "List of application config files separated by comma")
var loggerConfigPath = flag.String("logger.config.file", "not-found.yaml", "Path to logger config file")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

func main() {
	flag.Parse()

	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	logs, err := logs.New(*loggerConfigPath)
	if err != nil {
		panic(err)
	}

	loggerName := "main"
	logger := logs.WithName(loggerName)

	appConfig := getAppConfig(*appConfigPaths)
	appConfig.PostgresDatabase.User = os.Getenv("POSTGRES_USER")
	appConfig.PostgresDatabase.Password = os.Getenv("POSTGRES_PASSWORD")

	fmt.Printf("%+v\n", appConfig)

	postgresDb, err := postgres.New(&appConfig.PostgresDatabase)
	if err != nil {
		panic(err)
	}

	err = postgresDb.AutoMigrate(&postgresuser.User{})
	if err != nil {
		panic(err.Error())
	}

	userStorage := postgresuser.NewUserStorage(postgresDb, logs)

	userService := user.NewUserService(userStorage, logs)
	authService := auth.NewAuthService(userService, &appConfig.AuthService)

	appServer := httpserver.New(&appConfig.HttpServer, logs)
	appServer.Server.Group(httproutes.RegisterRoute(authService))
	appServer.Server.Group(httproutes.LoginRoute(authService))

	websocketServer := websocketserver.New(&appConfig.WebsocketServer, logs)
	websocketServer.Server.Group(websocket.ServeWebsocketRoute(websocketServer.Hub, authService, logs))

	appServerStarter := appServer.StartAsync()
	websocketServerStarter := websocketServer.StartAsync()

	select {
	case appServerErr := <-(*appServerStarter):
		logger.Errorf("Cannot start app server: %v", appServerErr)
	case websocketErr := <-(*websocketServerStarter):
		logger.Errorf("Cannot start websocket server: %v", websocketErr)
	}
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
