package main

import (
	"log"

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

	r.GET("/health", handlers.HealthCheck)

	r.GET("/ws", func(c *gin.Context) {
		handlers.HandleWebSocket(c, hub, cfg)
	})

	r.StaticFile("/", "./static/index.html")

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("server failed: ", err)
	}
}
