package router

import (
	"github.com/labstack/echo"
	"chattoo/user"
	"net/http"
	"chattoo/jwt"
	"github.com/spf13/viper"
)

type signUpResponse struct {
	User  user.User `json:"user"`
	Token string    `json:"token"`
}

func signUp(c echo.Context) error {
	cred := new(user.Credentials)
	if err := c.Bind(cred); err != nil {
		return err
	}

	if cred.Username == "" || cred.Password == "" {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Please fill required fields.")
	}

	userStore := c.(*userStoreContext).Store()
	u := user.New(cred)

	if userStore.IsExists(u.Username) {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Username already taken. Please, try another username.")
	}

	if err := userStore.Insert(u); err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Authorization broken.")
	}

	t, err := jwt.New(u, viper.GetString("jwt.secret"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Authorization broken.")
	}

	return c.JSON(http.StatusOK, signUpResponse{u, t})
}
