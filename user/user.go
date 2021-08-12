package user

import "gorm.io/gorm"

var DB *gorm.DB
var err error

type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
}