package models

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"log"
)

var db *gorm.DB;

type User struct {
	gorm.Model
	Username string
	Email string
	Password string
}

func connect() {
	var err error;

	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to connect to database")
	}
	db.AutoMigrate(&User{}) // Migrate the User model into the database schema

}

func CreateUser(user_data *User) (*gorm.DB) {
	connect() // Connect to database

	result := db.Create(&User{
			Username: user_data.Username,
			Email: user_data.Email,
			Password: user_data.Password,
		})

	return result
}