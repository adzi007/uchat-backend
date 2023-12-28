package http

import (
	"adzi-clean-architecture/domain"
	"adzi-clean-architecture/user/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Us domain.UserService
}

func NewRouteUser(r *fiber.App, us domain.UserService) {

	handler := &UserHandler{
		Us: us,
	}

	r.Get("/user-test", middleware.Auth, handler.GetUserHandler)

}

func (uh *UserHandler) GetUserHandler(ctx *fiber.Ctx) error {

	data := uh.Us.GetUserAll()

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"pesan": "lorem ipsum dolor sit amet",
		"data":  data,
	})
}
