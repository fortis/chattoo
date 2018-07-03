package main

import (
	"chattoo/server"
	"github.com/spf13/viper"
	"chattoo/router"
)

func main() {
	viper.SetConfigFile("config/app.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	wsServer := server.NewServer()
	e := router.Load(wsServer)

	go wsServer.Listen()

	addr := viper.GetString("app.addr")
	e.Logger.Fatal(e.Start(addr))
}
