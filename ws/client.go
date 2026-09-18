package ws

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	username string
	send     chan Message
}

func NewClient(hub *Hub, conn *websocket.Conn, username string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		username: username,
		send:     make(chan Message),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	for {
		_, data, err := c.conn.ReadMessage()

		if err != nil {
			break
		}

		var message Message

		if err := json.Unmarshal(data, &message); err != nil {
			continue
		}

		message.User = c.username

		c.hub.broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()

	for message := range c.send {
		data, err := json.Marshal(message)

		if err != nil {
			continue
		}

		err = c.conn.WriteMessage(websocket.TextMessage, data)

		if err != nil {
			return
		}
	}
}
