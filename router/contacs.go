package router

import (
	"github.com/labstack/echo"
	"chattoo/user"
	"net/http"
)

type contactsResponse struct {
	Contacts []user.User `json:"contacts"`
}

func contacts(c echo.Context) error {
	userStore := c.(*userStoreContext).Store()
	users, err := userStore.FindAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &contactsResponse{users})
}
