package services

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	sqldb "wlbt.nl/walkr/db"
	"wlbt.nl/walkr/db/models"
	database "wlbt.nl/walkr/db/sqlc"
	v "wlbt.nl/walkr/validation"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

func CreateUser(c *fiber.Ctx) error {
	c.Accepts("application/json")
	ctx := c.Context()

	var req CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println(err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid request body",
			"detail": err.Error(),
			"for":    "any",
		})
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if err, ok := v.Validate(req.Username, c, v.IsNotEmpty("username"), v.IsMinLength("username", 3), v.IsMaxLength("username", 20)); !ok {
		return err
	}

	if err, ok := v.Validate(req.Email, c, v.IsNotEmpty("email"), v.IsEmail("email")); !ok {
		return err
	}

	if err, ok := v.Validate(req.Password, c, v.IsNotEmpty("password"), v.IsMinLength("password", 8), v.IsStrongPassword("password")); !ok {
		return err
	}

	if user, err := models.CreateUser(ctx, req.Username, req.Email, req.Password); err != nil {
		log.Println(err)

		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":  "Username or email already taken",
			"detail": err.Error(),
			"for":    strings.ToLower(strings.Split(err.Error(), " ")[0]),
		})
	} else {
		if token, err := models.CreateToken(ctx, user.ID); err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":  "Error generating token",
				"detail": err.Error(),
				"for":    "any",
			})
		} else {
			c.Cookie(&fiber.Cookie{
				HTTPOnly: true,
				Name:     "walkr-session",
				MaxAge:   int((time.Hour * 24 * 30).Seconds()),
				Value:    token.Token,
			})
			return c.Status(fiber.StatusCreated).JSON(user)
		}
	}
}

func GetAllUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	if users, err := sqldb.Queries.ListUsers(ctx); err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something wernt wrong",
		})
	} else if users == nil {
		return c.JSON([]database.User{})
	} else {
		return c.JSON(users)
	}
}

func GetUserById(c *fiber.Ctx) error {
	ctx := c.Context()
	userId, err := strconv.ParseInt(c.Params("id"), 10, 64)

	if err != nil {
		log.Println(err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse route parameters",
		})
	}

	if user, err := sqldb.Queries.GetUserById(ctx, userId); err != nil {
		log.Println(err)

		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	} else {
		return c.JSON(user)
	}
}

func GetSelf(c *fiber.Ctx) error {
	token := c.Cookies("walkr-session")

	if token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "No session token",
		})
	}

	if data, err := models.GetUserByToken(c.Context(), token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong",
		})
	} else {
		return c.JSON(data)
	}
}

func UserLogin(c *fiber.Ctx) error {
	c.Accepts("application/json")
	ctx := c.Context()

	var req UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println(err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Invalid request body",
			"detail": err.Error(),
			"for":    "any",
		})
	}

	if err, ok := v.Validate(req.UsernameOrEmail, c, v.IsNotEmpty("username_or_email")); !ok {
		return err
	}

	if err, ok := v.Validate(req.Password, c, v.IsNotEmpty("password")); !ok {
		return err
	}

	err, isEmail := v.Validate(req.UsernameOrEmail, c, v.IsEmail("username_or_email"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":  "Something went wrong",
			"detail": "We failed to determine if you entered a username or a password",
			"for":    "username_or_email",
		})
	}

	var user *database.User

	if isEmail {
		email := req.UsernameOrEmail

		if data, err := models.GetUserByEmail(ctx, email); err != nil || data == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":  "Not found",
				"detail": "We couldn't find a user with that e-mail",
				"for":    "username_or_email",
			})
		} else {
			user = data
		}
	} else {
		username := req.UsernameOrEmail

		if data, err := models.GetUserByName(ctx, username); err != nil || data == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":  "Not found",
				"detail": "We couldn't find a user with that username",
				"for":    "username_or_email",
			})
		} else {
			user = data
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":  "Forbidden",
			"detail": "Password incorrect",
			"for":    "password",
		})
	}

	if token, err := models.CreateToken(ctx, user.ID); err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":  "Error generating token",
			"detail": err.Error(),
			"for":    "any",
		})
	} else {
		c.Cookie(&fiber.Cookie{
			HTTPOnly: true,
			Name:     "walkr-session",
			MaxAge:   int((time.Hour * 24 * 30).Seconds()),
			Value:    token.Token,
		})

		return c.Status(fiber.StatusOK).JSON(user)
	}
}
