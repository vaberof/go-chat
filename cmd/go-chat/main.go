package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	httproutes "github.com/vaberof/go-chat/internal/app/entrypoint/http"
	"github.com/vaberof/go-chat/internal/domain/auth"
	"github.com/vaberof/go-chat/internal/domain/message"
	"github.com/vaberof/go-chat/internal/domain/room"
	"github.com/vaberof/go-chat/internal/domain/user"
	"github.com/vaberof/go-chat/internal/infra/storage/postgres"
	"github.com/vaberof/go-chat/internal/websocket"
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

	appConfig := getAppConfig(*appConfigPaths)
	appConfig.PostgresDatabase.User = os.Getenv("POSTGRES_USER")
	appConfig.PostgresDatabase.Password = os.Getenv("POSTGRES_PASSWORD")

	fmt.Printf("%+v\n", appConfig)

	postgresDb, err := postgres.New(&appConfig.PostgresDatabase)
	if err != nil {
		panic(err)
	}

	err = postgresDb.AutoMigrate(&postgres.User{}, &postgres.Room{}, &postgres.Member{}, &postgres.Message{})
	if err != nil {
		panic(err.Error())
	}

	userStorage := postgres.NewUserStorage(postgresDb, logs)
	roomStorage := postgres.NewRoomStorage(postgresDb, logs)
	messageStorage := postgres.NewMessageStorage(postgresDb, logs)

	userService := user.NewUserService(userStorage, logs)
	roomService := room.NewRoomService(roomStorage, logs)
	messageService := message.NewMessageService(messageStorage, logs)
	authService := auth.NewAuthService(userService, &appConfig.AuthService)

	appServer := httpserver.New(&appConfig.HttpServer, logs)
	hub := websocket.NewHub(messageService, roomService, userService, logs)

	go hub.Run()

	httpHandler := httproutes.NewHandler(hub, authService, userService, roomService, messageService)
	httpHandler.InitRoutes(appServer.Server, logs)

	appServerStarter := appServer.StartAsync()

	<-(*appServerStarter)
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
