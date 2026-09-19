package ws

const maxClients = 2

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

type Hub struct {
	clients    map[*Client]bool
	register   chan RegisterRequest
	unregister chan *Client
	broadcast  chan Message
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
		broadcast:  make(chan Message),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case request := <-h.register:
			if len(h.clients) >= maxClients || h.usernameTaken(request.client.username) {
				request.result <- false
				continue
			}

			h.clients[request.client] = true
			request.result <- true

		case client := <-h.unregister:
			h.removeClient(client)

		case message := <-h.broadcast:
			for client := range h.clients {
				if client.username == message.User {
					continue
				}
				select {
				case client.send <- message:
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
	}
}
func (h *Hub) usernameTaken(username string) bool {
	for client := range h.clients {
		if client.username == username {
			return true
		}
	}

	return false
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
