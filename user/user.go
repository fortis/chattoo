package user

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func New(credentials *Credentials) User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), 8)
	return User{
		Username: credentials.Username,
		Password: string(hashedPassword),
	}
}
