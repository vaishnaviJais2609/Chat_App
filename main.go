package main

import (
	"chatapp/config"
	"chatapp/handlers"
	"chatapp/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	hub := ws.NewHub()
	go hub.Run()

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		handlers.HandleWebSocket(c, hub, cfg)
	})

	r.StaticFile("/", "./static/index.html")

	r.Run(":" + cfg.Port)
}
