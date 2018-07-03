package user

import "gopkg.in/go-playground/validator.v9"

type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CredentialsValidator struct {
	Validator *validator.Validate
}

func (cv *CredentialsValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
