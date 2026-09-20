package ws

import (
	"log"
	"sync"
)

const maxClients = 2

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type BroadcastMessage struct {
	Sender  *Client
	Message Message
}

type Hub struct {
	clients    map[*Client]bool
	register   chan RegisterRequest
	unregister chan *Client
	broadcast  chan BroadcastMessage
	mu         sync.RWMutex
	count      int
}

type RegisterRequest struct {
	client *Client
	result chan bool
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan RegisterRequest),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case request := <-h.register:
			if len(h.clients) >= maxClients {
				log.Printf("registration rejected: %d/%d clients connected", len(h.clients), maxClients)
				request.result <- false
				continue
			}

			h.clients[request.client] = true
			h.mu.Lock()
			h.count = len(h.clients)
			h.mu.Unlock()
			log.Printf("client registered: %d/%d clients connected", len(h.clients), maxClients)
			request.result <- true

		case client := <-h.unregister:
			h.removeClient(client)

		case message := <-h.broadcast:
			for client := range h.clients {
				if client == message.Sender {
					continue
				}
				select {
				case client.send <- message.Message:
				default:
					h.removeClient(client)
				}
			}
		}
	}
}

func (h *Hub) removeClient(client *Client) {
	if h.clients[client] {
		delete(h.clients, client)
		close(client.send)
		h.mu.Lock()
		h.count = len(h.clients)
		h.mu.Unlock()
		log.Printf("client unregistered: %d/%d clients connected", len(h.clients), maxClients)
	}
}

func (h *Hub) Register(client *Client) bool {
	result := make(chan bool, 1)

	h.register <- RegisterRequest{
		client: client,
		result: result,
	}

	return <-result
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.count
}
