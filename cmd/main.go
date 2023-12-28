package main

import (
	"adzi-clean-architecture/authentication"
	"adzi-clean-architecture/chat"
	"adzi-clean-architecture/config"
	"adzi-clean-architecture/pkg/logger"
	"adzi-clean-architecture/user"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
)

func main() {

	config.LoadConfig()

	app := fiber.New()

	mylog := logger.NewLogger()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &mylog,
	}))

	dbConnectionCongig := config.NewDatabase()
	_, err := dbConnectionCongig.ConnectDB()

	if err != nil {

		logger.Fatal().Err(err)
	}

	user.InitUser(app, dbConnectionCongig)
	authentication.InitAuthentication(app, dbConnectionCongig)
	chat.InitChat(app, dbConnectionCongig)

	if err := app.Listen(":" + config.ENV.PORT_AP); err != nil {
		logger.Fatal().Err(err).Msg("Fiber app error")
	}

}
