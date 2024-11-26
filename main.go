package main

import (
	"github.com/astianmuchui/go-auth/auth"
	"github.com/astianmuchui/go-auth/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/django/v3"
	"log"
)

var store = session.New()

func main() {

	engine := django.New("./templates", ".django")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	/* Home route */
	app.Get("/", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		errorMsg := sess.Get("signup_error")
		sess.Delete("signup_error")
		sess.Save()

		return c.Render("index", fiber.Map{
			"signup_error": errorMsg,
		})
	})

	app.Post("/register", func(context *fiber.Ctx) error {
		payload := new(models.User)

		if err := context.BodyParser(payload); err != nil {
			return err
		}
		sess, _ := store.Get(context)

		// Create the user if already does not exist
		if models.UserExists(payload) == false {
			// Add user to database
			result := models.CreateUser(payload)

			if result != nil && result.Error == nil {
				sess.Set("user_email", payload.Email)
				sess.Set("logged_in", true)

				sess.Save()

				return context.Redirect("/dashboard")
			} else {
				sess.Set("signup_error", "Unable to sign up")
				sess.Save()
			}
		} else {
			sess.Set("signup_error", "User already exists")
			log.Println("User found")
			sess.Save()
			return context.Redirect("/")
		}

		log.Println("Username:", payload.Username)
		log.Println("Email:", payload.Email)
		log.Println("Password:", payload.Password)

		return context.SendStatus(fiber.StatusCreated)

	})

	app.Get("/login", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		errorMsg := sess.Get("login_error")
		sess.Delete("login_error")
		sess.Save()

		return c.Render("login", fiber.Map{
			"login_error": errorMsg,
		})
	})

	app.Post("/signin", func(c *fiber.Ctx) error {
		payload := new(models.User)

		if err := c.BodyParser(payload); err != nil {
			return err
		}

		var u models.User
		u = models.GetUserDataByUsername(payload.Username)
		log.Println(u)
		log.Println(payload)

		userVerified := auth.Login(payload)
		if userVerified {
			sess, _ := store.Get(c)
			sess.Set("user_email", u.Email)
			sess.Set("logged_in", true)

			sess.Save()

			return c.Redirect("/dashboard")
		}
		return c.Redirect("/login")
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {

		sess, _ := store.Get(c)

		logged_in := sess.Get("logged_in")
		log.Println(logged_in)
		if logged_in == true {
			userEmail := sess.Get("user_email").(string)
			log.Println(userEmail)
			user_data := models.GetUserDataByEmail(userEmail)
			return c.Render("dashboard", fiber.Map{
				"username": user_data.Username,
				"email":    user_data.Email,
			})
		} else {
			return c.Redirect("/login")
		}
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		sess.Destroy()
		return c.Redirect("/login")
	})

	app.Get("/update", func(c *fiber.Ctx) error {

		sess, _ := store.Get(c)
		logged_in := sess.Get("logged_in")

		if logged_in == true {
			userEmail := sess.Get("user_email").(string)
			var update_err, update_success string
			if sess.Get("update_error") != nil {
				update_err = sess.Get("update_error").(string)
			}
			if sess.Get("update_success") != nil {
				update_success = sess.Get("update_success").(string)
			}
			sess.Delete("update_error")
			sess.Delete("update_success")
			sess.Save()
			log.Println(userEmail)
			user_data := models.GetUserDataByEmail(userEmail)

			return c.Render("update", fiber.Map{
				"username":       user_data.Username,
				"email":          user_data.Email,
				"update_error":   update_err,
				"update_success": update_success,
			})
		} else {
			return c.Redirect("/login")
		}
	})

	app.Post("/update-profile", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		logged_in := sess.Get("logged_in")

		if logged_in == true {
			payload := new(models.User)

			if err := c.BodyParser(payload); err != nil {
				return err
			}

			u := models.GetUserDataByEmail(sess.Get("user_email").(string))
			u.Username = payload.Username
			u.Email = payload.Email

			result := models.DB.Save(u)

			if result.Error == nil {
				sess.Set("update_success", "Data Updated successfully")
				log.Println("Updated data")
				sess.Save()
				return c.Redirect("/update")

			} else {
				log.Println("Did not Update data")

				sess.Set("update_error", "Could not update profile")
				sess.Save()
				return c.Redirect("/update")

			}
		} else {
			return c.Redirect("/login")
		}
	})

	app.Listen(":8081")
}
