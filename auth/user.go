package auth

import (
	"github.com/astianmuchui/go-auth/models"
)

func Login(user *models.User) (bool) {
	var u models.User
	models.Connect()

	result := models.DB.First(&u, "username = ?", u.Username)
	if result.RowsAffected == 0 {
		return (false)
	} else {
		// Verify password
		if models.Password_verify(u.Password, []byte(user.Password)) == true {
			return (true)
		}
	}

	return (false)
}