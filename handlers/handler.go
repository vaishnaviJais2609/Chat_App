package handlers

import (
	"log"
	"net/http"

	"chatapp/config"
	"chatapp/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleWebSocket(c *gin.Context, hub *ws.Hub, cfg *config.Config) {
	username := c.Query("user")

	if username == "" {
		c.String(http.StatusBadRequest, "Username is required")
		return
	}

	if username != cfg.UserA && username != cfg.UserB {
		c.String(http.StatusForbidden, "Unknown user")
		return
	}

	upgrader := websocket.Upgrader{}

	if cfg.AllowAllOrigins {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("websocket upgrade failed:", err)
		return
	}

	client := ws.NewClient(hub, conn, username)

	if !hub.Register(client) {
		conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.ClosePolicyViolation,
				"Chat is full or this user is already connected",
			),
		)
		conn.Close()
		return
	}

	go client.WritePump()
	client.ReadPump()
}

func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
