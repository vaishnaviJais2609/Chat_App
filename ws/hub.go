package ws

import (
	"encoding/json"
	"log"
	"time"
)

type Message struct {
	From      string `json:"from"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
}
type Hub struct {
	clients map[string]*Client

	register   chan *Client
	unregister chan *Client
	broadcast  chan Message

	AllowedUsers map[string]bool
}

func NewHub(userA, userB string) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
		AllowedUsers: map[string]bool{
			userA: true,
			userB: true,
		},
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Username] = client
			log.Printf("client connected: %s (total: %d)", client.Username, len(h.clients))
			h.notifySystem(client.Username + " has joined the chat")

		case client := <-h.unregister:
			if _, ok := h.clients[client.Username]; ok {
				delete(h.clients, client.Username)
				close(client.Send)
				log.Printf("client disconnected: %s (total: %d)", client.Username, len(h.clients))
				h.notifySystem(client.Username + " has left the chat")
			}

		case msg := <-h.broadcast:
			h.deliver(msg)
		}
	}
}
func (h *Hub) deliver(msg Message) {
	payload, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return
	}

	for username, client := range h.clients {
		if username == msg.From {
			continue
		}
		select {
		case client.Send <- payload:
		default:
			close(client.Send)
			delete(h.clients, username)
		}
	}
}
func (h *Hub) notifySystem(text string) {
	msg := Message{
		From:      "system",
		Content:   text,
		Timestamp: time.Now().Unix(),
		Type:      "system",
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return
	}
	for username, client := range h.clients {
		select {
		case client.Send <- payload:
		default:
			close(client.Send)
			delete(h.clients, username)
		}
	}
}
func (h *Hub) Register(c *Client) {
	h.register <- c
}
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}
func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}
func (h *Hub) IsFull() bool {
	return len(h.clients) >= 2
}
func (h *Hub) IsUsernameTaken(username string) bool {
	_, ok := h.clients[username]
	return ok
}
