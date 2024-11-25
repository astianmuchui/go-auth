package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"github.com/astianmuchui/go-auth/models"
	"github.com/gofiber/template/django/v3"
	"github.com/astianmuchui/go-auth/auth"

)

func main() {

	engine := django.New("./templates", ".django")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	/* Home route */
	app.Get("/", func (c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello world",
		})
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

		return context.SendString("Username: " + payload.Username + " Email: " + payload.Email + " Password: " + payload.Password)
	})

	app.Get("/login", func (c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})
	app.Post("/signin", func (c *fiber.Ctx) error {
		payload := new(models.User)

		if err := c.BodyParser(payload); err != nil {
			return err
		}

		var user_verified bool = auth.Login(payload)
		if user_verified == true {
			return c.Redirect("/login?logged_in")
		}

		log.Println("Username:", payload.Username)
		log.Println("Password:", payload.Password)
		return c.SendString("Username: " + payload.Username + " Password: " + payload.Password)

	})
	app.Listen(":8081")
}