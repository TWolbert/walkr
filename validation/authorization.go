package v

import (
	"github.com/gofiber/fiber/v2"
	"wlbt.nl/walkr/db/models"
)

func IsAdmin(c *fiber.Ctx) error {
	token := c.Cookies("walkr-session")

	if token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "No session token",
		})
	}

	user, err := models.GetUserByToken(c.Context(), token)

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "Forbidden",
			"detail": "Token didn't resolve an authorized user",
		})
	}

	role, err := models.UserGetRole(c.Context(), user)

	if role.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "Forbidden",
			"detail": "You are not allowed to access this route",
		})
	}

	return c.Next()
}

func IsLoggedIn(c *fiber.Ctx) error {
	token := c.Cookies("walkr-session")

	if token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "No session token",
		})
	}

	user, err := models.GetUserByToken(c.Context(), token)

	if err != nil || user == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "Forbidden",
			"detail": "Token didn't resolve an authorized user",
		})
	}

	return c.Next()
}
