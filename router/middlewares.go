package router

import (
	"github.com/labstack/echo"
	"chattoo/store"
	"github.com/spf13/viper"
)

type userStoreContext struct {
	echo.Context
	userStore *store.UserStore
}

func (c *userStoreContext) Store() *store.UserStore {
	return c.userStore
}

func userStoreMiddleware(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userStore, _ := store.NewUserStore(viper.GetString("database.addr"))
		return handler(&userStoreContext{c, userStore})
	}
}
