package router

import (
	"github.com/labstack/echo"
	"chattoo/user"
	"golang.org/x/crypto/bcrypt"
	"chattoo/jwt"
	"github.com/spf13/viper"
	"net/http"
)

type signInResponse struct {
	User  user.User `json:"user"`
	Token string    `json:"token"`
}

func signIn(c echo.Context) error {
	cred := new(user.Credentials)
	if err := c.Bind(cred); err != nil {
		return err
	}

	if cred.Username == "" || cred.Password == "" {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Please fill required fields.")
	}

	userStore := c.(*userStoreContext).Store()
	var u user.User
	userStore.FindOneByName(cred.Username, &u)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(cred.Password))
	if err != nil {
		return echo.ErrUnauthorized
	}

	t, err := jwt.New(u, viper.GetString("jwt.secret"))
	if err != nil {
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, &signInResponse{u, t})
}
