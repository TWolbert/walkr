package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	sqldb "wlbt.nl/walkr/db"
	"wlbt.nl/walkr/services"
	v "wlbt.nl/walkr/validation"
)

func main() {
	log.SetPrefix("Walkr: ")

	if err := sqldb.Run(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		Prefork: true,
	})

	app.Get("/api/user", v.IsAdmin, services.GetAllUsers)
	app.Get("/api/user/:id", services.GetUserById)
	app.Post("/api/user", services.CreateUser)
	app.Post("/api/user/login", services.UserLogin)
	app.Get("/api/self", v.IsLoggedIn, services.GetSelf)

	app.Get("/", services.ServeHTML)
	app.Get("*", services.ServeClient)

	log.Fatal(app.Listen(":8000"))
}
