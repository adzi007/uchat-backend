package middleware

import (
	"adzi-clean-architecture/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Auth(ctx *fiber.Ctx) error {

	token := ctx.Get("Authorization")
	splitToken := strings.Split(token, "Bearer ")
	token = splitToken[1]

	if token == "" {

		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	claims, err := utils.DecodeToken(token)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	// role := claims["role"].(string)

	// if role != "admin" {
	// 	return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"message": "forbiden access",
	// 	})
	// }

	ctx.Locals("userInfo", claims)

	return ctx.Next()
}
