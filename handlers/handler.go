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
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

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

	client := ws.NewClient(hub, conn)

	if !hub.Register(client) {
		log.Println("registration rejected: chat is full")
		// Send a JSON message the client can parse before closing
		errMsg := map[string]string{"error": "Chat is full."}
		conn.WriteJSON(errMsg)
		conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.ClosePolicyViolation,
				"Chat is full",
			),
		)
		conn.Close()
		return
	}

	log.Printf("client registered, total clients: %d", hub.ClientCount())
	go client.WritePump()
	client.ReadPump()
}

func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
