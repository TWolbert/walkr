package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("APP_ENV") == "local" {
		log.Fatal(app.Listen(":8000"))
	} else {
		log.Fatal(app.ListenTLS(":8000", os.Getenv("CERT_PATH"), os.Getenv("KEY_PATH")))
	}

}
