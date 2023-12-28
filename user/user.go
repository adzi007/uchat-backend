package user

import (
	"adzi-clean-architecture/config"
	"adzi-clean-architecture/user/delivery/http"
	"adzi-clean-architecture/user/repository"
	"adzi-clean-architecture/user/service"

	"github.com/gofiber/fiber/v2"
)

func InitUser(app *fiber.App, dbConfig config.MongoDbInterface) {

	// Repository
	userCollection := dbConfig.GetCollection("user")
	userRepo := repository.NewUserRepo(userCollection)

	// Service
	userService := service.NewUserService(userRepo)

	http.NewRouteUser(app, userService)

}
