package handlers
package handlers

import (
	"log"

	"chatapp/config"
	"chatapp/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewUpgrader(cfg *config.Config) websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return cfg.AllowAllOrigins
		},
	}
}

func ServeWs(hub *ws.Hub, upgrader websocket.Upgrader) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("user")

		if username == "" {
			c.String(400, "missing required query parameter: user")
			return
		}
		if !hub.AllowedUsers[username] {
			c.String(403, "unknown user for this chat")
			return
		}
		if hub.IsUsernameTaken(username) {
			c.String(409, "this user is already connected elsewhere")
			return
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("websocket upgrade failed: %v", err)
			return
		}
		client := ws.NewClient(hub, conn, username)
		hub.Register(client)

		go client.WritePump()
		go client.ReadPump()
	}
}

func HealthCheck(c *gin.Context) {
	c.String(200, "ok")
}