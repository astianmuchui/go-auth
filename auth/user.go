package auth

import (
	"github.com/astianmuchui/go-auth/models"
	"log"
)


func Login(user *models.User) (bool) {
	var u models.User

	models.Connect()

	result := models.DB.First(&u, "username = ?", user.Username)
	if result.RowsAffected == 0 {
		log.Println("User not found")
		return (false)
	} else {
		// Verify password
		if models.Password_verify(u.Password, []byte(user.Password)) == true {
			log.Println("Login successful")

			return (true)
		} else {
			log.Println("Login not successful")
			return (false)
		}
	}
}