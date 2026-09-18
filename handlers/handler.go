package handlers

import (
	"log"
	"net/http"

	"chatapp/config"
	"chatapp/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func HandleWebSocket(c *gin.Context, hub *ws.Hub, cfg *config.Config) {
	username := c.Query("user")

	if username == "" {
		c.String(400, "Username is required")
		return
	}

	if username != cfg.UserA && username != cfg.UserB {
		c.String(403, "Unknown user")
		return
	}

func HealthCheck(c *gin.Context) {
	c.String(200, "ok")
}
