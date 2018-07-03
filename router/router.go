package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/dgrijalva/jwt-go"
	"chattoo/user"
	"chattoo/server"
)

func Load(wsServer *server.Server) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{viper.GetString("cors.origin")},
		AllowMethods: []string{echo.GET, echo.POST},
	}))
	e.Use(userStoreMiddleware)

	e.POST("/signin", signIn)
	e.POST("/signup", signUp)
	e.File("/", "web/index.html")
	e.Static("/static", "web/static")

	// Private api. Token required.
	a := e.Group("/api")
	a.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString("jwt.secret")),
		TokenLookup: "header:Authorization",
	}))
	a.Use(userStoreMiddleware)
	a.GET("/contacts", contacts)
	a.POST("/rename", rename)

	// Private routes. Token required.
	p := e.Group("/private")
	p.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(viper.GetString("jwt.secret")),
		TokenLookup: "query:itok",
	}))
	p.Use(userStoreMiddleware)
	p.GET("/ws", func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		name := claims["name"].(string)

		var u user.User
		c.(*userStoreContext).Store().FindOneByName(name, &u)
		wsServer.HandleWS(u).ServeHTTP(c.Response(), c.Request())
		return nil
	})

	return e
}
