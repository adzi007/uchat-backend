package http

import (
	"adzi-clean-architecture/domain"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AuthService domain.AuthService
}

func NewRouteAuth(r *fiber.App, us domain.AuthService) {

	handler := &AuthHandler{
		AuthService: us,
	}

	r.Post("/login", handler.Login)

}

func (uh *AuthHandler) Login(ctx *fiber.Ctx) error {

	loginRequest := new(domain.LoginRequest)

	if err := ctx.BodyParser(loginRequest); err != nil {
		return err
	}

	// Validation Login Request

	var validate = validator.New()

	errValidate := validate.Struct(loginRequest)

	if errValidate != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "failed",
			"error":   errValidate.Error(),
		})
	}

	// Authetication service

	token, statusCode, err := uh.AuthService.Login(loginRequest)

	if err != nil {
		return ctx.Status(statusCode).JSON(fiber.Map{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"pesan": "test login request",
		"token": token,
	})
}
