package main

import (
	"log"

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

	app.Get("/api/user", services.GetAllUsers)
	app.Get("/api/user/:id", services.GetUserById)
	app.Post("/api/user", services.CreateUser)
	app.Get("/api/self", services.GetSelf)

	app.Get("/", services.ServeHTML)
	app.Get("*", services.ServeClient)

	log.Fatal(app.Listen(":8000"))
}
