package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"github.com/astianmuchui/go-auth/models"
)

func main() {



	app := fiber.New()

	/* Create Home route */
	app.Get("/", func (c *fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})

	app.Post("/register", func (context *fiber.Ctx) error {
		payload := new(models.User)

		if err := context.BodyParser(payload); err != nil {
			return err
		}

		// Add user to database
		result := models.CreateUser(payload)

		if result.Error != nil {
			return context.SendStatus(fiber.StatusCreated)
		}

		log.Println("Username:", payload.Username)
		log.Println("Email:", payload.Email)
		log.Println("Password:", payload.Password)

		return context.SendString("Username: " + payload.Username + " Email: " + payload.Email + "Password: " + payload.Password)
	})

	app.Listen(":8081")
}