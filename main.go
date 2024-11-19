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
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{}) // Initialize ORM Connection

	if err != nil { // Check for database connection error
	  panic("failed to connect database")
	}

	db.AutoMigrate(&User{}) // Migrate the User model into the database schema

	app := fiber.New()

	/* Create Home route */
	app.Get("/", func (c *fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})

	app.Post("/register", func (context *fiber.Ctx) error {
		payload := new(RegisterPayload)

		if err := context.BodyParser(payload); err != nil {
			return err
		}

		// Add user to database
		result := db.Create(&User{
			Username: payload.Username,
			Email: payload.Email,
			Password: payload.Password,
		})

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