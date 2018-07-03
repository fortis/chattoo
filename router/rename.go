package router

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/dgrijalva/jwt-go"
)

func rename(c echo.Context) error {
	req := map[string]string{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	username := req["username"]
	userStore := c.(*userStoreContext).Store()
	if userStore.IsExists(username) {
		return echo.NewHTTPError(http.StatusNotAcceptable, "Username already taken. Please try another.")
	}

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	id := int64(claims["id"].(float64))
	userStore.UpdateUsername(id, username)

	return c.JSON(http.StatusOK, "")
}
