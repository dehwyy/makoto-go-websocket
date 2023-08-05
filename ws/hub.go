package ws

import (
	"fmt"

	l "github.com/dehwyy/makoto-go-websocket/logger"
)

type Hub struct {
	// UID
	UID string
	// channels for:
	// - registering user
	// - unregistering user
	// - broadcasting messages
	reg       chan (*Client)
	unreg     chan (*Client)
	broadcast chan []byte
	// Map for each client
	clients map[*Client]bool
}

func NewHub(uid string) *Hub {
	return &Hub{
		UID:       uid,
		reg:       make(chan *Client),
		unreg:     make(chan *Client),
		broadcast: make(chan []byte),
		clients:   make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {

		// when someone connecting to the Hub
		case client := <-h.reg:
			l.Log(fmt.Sprintf("New client connected: %v on hub '%s'", client.conn.RemoteAddr().String(), h.UID))
			h.clients[client] = true

			// on disconnect
		case client := <-h.unreg:

			// clarifying whether client exists in map
			if _, ok := h.clients[client]; !ok {
				continue
				// log.Fatalf("Client does not exist: %v", client)
			}

			// delete from map AND close client.send chan
			delete(h.clients, client)
			close(client.send)

		// on broadcast
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				// if buffered channel is not full
				case client.send <- message:
					// Already writed to the channel

				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}
