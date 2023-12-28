package authentication

import (
	"adzi-clean-architecture/authentication/http"
	"adzi-clean-architecture/authentication/service"
	"adzi-clean-architecture/config"
	"adzi-clean-architecture/user/repository"

	"github.com/gofiber/fiber/v2"
)

func InitAuthentication(app *fiber.App, dbConfig config.MongoDbInterface) {

	//repositori
	userCollection := dbConfig.GetCollection("user")
	userRepo := repository.NewUserRepo(userCollection)

	// Service
	authService := service.NewauthService(userRepo)

	http.NewRouteAuth(app, authService)

}
