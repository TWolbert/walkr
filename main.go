package main

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	sqldb "wlbt.nl/walkr/db"
	"wlbt.nl/walkr/db/models"
	database "wlbt.nl/walkr/db/sqlc"
)

func main() {
	log.SetPrefix("Walkr: ")

	if err := sqldb.Run(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Get("/api/users", func(c *fiber.Ctx) error {
		ctx := c.Context()
		users, err := sqldb.Queries.ListUsers(ctx)

		if err != nil {
			log.Println(err)
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "No users in database",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Something wernt wrong",
			})
		}

		if users == nil {
			users = []database.User{}
		}

		return c.JSON(users)
	})

	app.Get("/user/:id", func(c *fiber.Ctx) error {
		ctx := c.Context()
		userId, err := strconv.ParseInt(c.Params("id"), 10, 64)

		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse route parameters",
			})
		}

		user, err := sqldb.Queries.GetUserById(ctx, userId)

		if err != nil {
			log.Println(err)
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Something went wrong",
			})
		}

		return c.JSON(user)
	})

	type CreateUserRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	app.Post("/api/user", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		ctx := c.Context()

		var req CreateUserRequest

		if err := c.BodyParser(&req); err != nil {
			log.Println(err)
			// Body wasn't valid JSON or couldn't be parsed into the struct
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Invalid request body",
				"detail": err.Error(),
				"for":    "any",
			})
		}

		req.Username = strings.TrimSpace(req.Username)
		req.Email = strings.TrimSpace(req.Email)

		if len(req.Username) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Invalid request body",
				"detail": "Username can't be empty",
				"for":    "username",
			})
		}

		if len(req.Email) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Invalid request body",
				"detail": "Email can't be empty",
				"for":    "email",
			})
		}

		if len(req.Password) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Invalid request body",
				"detail": "Password can't be empty",
				"for":    "password",
			})
		}

		if len(req.Username) < 3 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Invalid request body",
				"detail": "Username too short",
				"for":    "username",
			})
		}

		if len(req.Password) < 8 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Invalid request body",
				"detail": "Password too short",
				"for":    "password",
			})
		}

		user, err := models.CreateUser(ctx, req.Username, req.Email, req.Password)

		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":  "Username or email already taken",
				"detail": err.Error(),
				"for":    "any",
			})
		}

		return c.JSON(user)
	})

	log.Fatal(app.Listen(":8000"))
}
