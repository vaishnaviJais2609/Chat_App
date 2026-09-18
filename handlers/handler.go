package handlers

import (
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

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(hub, conn, username)

	if !hub.Register(client) {
		conn.WriteMessage(
			websocket.TextMessage,
			[]byte(`{"user":"system","message":"Chat is full"}`),
		)
		conn.Close()
		return
	}

	go client.ReadPump()
	go client.WritePump()
}
