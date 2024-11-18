package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"log"
)

type User struct {
	gorm.Model
	Username string
	Email string
	Password string
}

type RegisterPayload struct {
	Username string
	Email string
	Password string
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
	  panic("failed to connect database")
	}

	db.AutoMigrate(&User{})

	app := fiber.New()

	/* Create Home route */
	app.Get("/", func (c *fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})

	app.Post("/register", func(context *fiber.Ctx) error {
		payload := new(RegisterPayload)

		if err := context.BodyParser(payload); err != nil {
			return err
		}
		db.Create(&User{
			Username: payload.Username,
			Email: payload.Email,
			Password: payload.Password,
		})

		log.Println("Username:", payload.Username)
		log.Println("Email:", payload.Email)
		log.Println("Password:", payload.Password)

		return context.SendString("Username: " + payload.Username + " Email: " + payload.Email + "Password: " + payload.Password)
	})

	app.Listen(":8081")
}