package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	sqldb "wlbt.nl/walkr/db"
	"wlbt.nl/walkr/services"
)

func main() {
	log.SetPrefix("Walkr: ")

	if err := sqldb.Run(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	htmlPath := filepath.Join((dir), "client/dist/index.html")
	html, err := os.ReadFile(htmlPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	app.Get("/api/user", services.GetAllUsers)
	app.Get("/api/user/:id", services.GetUserById)
	app.Post("/api/user", services.CreateUser)
	app.Get("/api/self", services.GetSelf)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Type("html").SendString(string(html))
	})

	app.Get("*", func(c *fiber.Ctx) error {
		path := c.Path()

		path = filepath.Join(dir, "client/dist", path)

		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return c.Type("html").SendString(string(html))
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		return c.SendFile(htmlPath, true)
	})

	log.Fatal(app.Listen(":8000"))
}
