package models

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"log"
	"golang.org/x/crypto/bcrypt"
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

func password_hash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost) // Hash password
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func password_verify(hashed string, raw_pwd []byte) bool {
	bytehash := []byte(hashed)
	err :=bcrypt.CompareHashAndPassword(bytehash, raw_pwd) // Compare hash
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CreateUser(user_data *User) (*gorm.DB) {
	connect() // Connect to database

	result := db.Create(&User{
			Username: user_data.Username,
			Email: user_data.Email,
			Password: password_hash([]byte(user_data.Password)), // Enter hashed password
		})

	return result
}