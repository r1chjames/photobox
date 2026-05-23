package websocket

import (
	"context"
	"log/slog"

	json "github.com/goccy/go-json"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// NewHub creates a new Hub and starts its run loop.
func NewHub(ctx context.Context) *Hub {
	h := &Hub{
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	go h.run(ctx)
	return h
}

func (h *Hub) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case client := <-h.register:
			h.clients[client] = true
			slog.Debug("websocket client connected", "total", len(h.clients))
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				slog.Debug("websocket client disconnected", "total", len(h.clients))
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client send buffer full; drop it
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Register adds a client to the hub. Called from the HTTP upgrade handler.
func (h *Hub) Register(c *Client) {
	h.register <- c
}

// BroadcastEvent marshals an Event to JSON and sends it to all connected clients.
// Returns immediately; the hub's run loop handles actual delivery.
func (h *Hub) BroadcastEvent(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal websocket event", "type", event.Type, "error", err)
		return
	}
	select {
	case h.broadcast <- data:
	default:
		// Broadcast channel full; drop oldest
		select {
		case <-h.broadcast:
		default:
		}
		h.broadcast <- data
	}
}

// ClientCount returns the number of currently connected clients.
func (h *Hub) ClientCount() int {
	// Approximate: we don't lock, but it's only for metrics/logging.
	return len(h.clients)
}
