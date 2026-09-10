package websocket

import (
	"context"
	"log/slog"
	"sync"

	json "github.com/goccy/go-json"
)

// outbound is a message queued for delivery. An empty WorkspaceID means the
// message is global (e.g. system health) and goes to every client; otherwise
// it is delivered only to clients registered in that workspace (issue #74).
type outbound struct {
	workspaceID string
	data        []byte
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	// mu guards clients. The run loop mutates it; ClientCount reads it from
	// other goroutines (metrics/logging), so reads and writes must be
	// synchronized to stay race-free under -race.
	mu         sync.RWMutex
	clients    map[*Client]bool
	broadcast  chan outbound
	register   chan *Client
	unregister chan *Client
}

// NewHub creates a new Hub and starts its run loop.
func NewHub(ctx context.Context) *Hub {
	h := &Hub{
		broadcast:  make(chan outbound, 256),
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
			h.mu.Lock()
			h.clients[client] = true
			n := len(h.clients)
			h.mu.Unlock()
			slog.Debug("websocket client connected", "total", n, "workspace", client.workspaceID)
		case client := <-h.unregister:
			h.mu.Lock()
			_, ok := h.clients[client]
			if ok {
				delete(h.clients, client)
				close(client.send)
			}
			n := len(h.clients)
			h.mu.Unlock()
			if ok {
				slog.Debug("websocket client disconnected", "total", n, "workspace", client.workspaceID)
			}
		case msg := <-h.broadcast:
			h.mu.RLock()
			targets := make([]*Client, 0, len(h.clients))
			for client := range h.clients {
				targets = append(targets, client)
			}
			h.mu.RUnlock()
			for _, client := range targets {
				// Tenant isolation: a workspace-scoped message is delivered
				// only to clients in that workspace. Global messages (empty
				// workspace) reach everyone.
				if msg.workspaceID != "" && client.workspaceID != msg.workspaceID {
					continue
				}
				select {
				case client.send <- msg.data:
				default:
					// Client send buffer full; drop it
					h.mu.Lock()
					if _, ok := h.clients[client]; ok {
						close(client.send)
						delete(h.clients, client)
					}
					h.mu.Unlock()
				}
			}
		}
	}
}

// Register adds a client to the hub. Called from the HTTP upgrade handler.
func (h *Hub) Register(c *Client) {
	h.register <- c
}

// BroadcastEvent marshals an Event and sends it to every connected client.
//
// Use this only for deployment-wide events that carry no tenant data (e.g.
// system health). Tenant-visible events MUST use BroadcastWorkspaceEvent, or
// they leak identifiers across tenants (issue #74).
func (h *Hub) BroadcastEvent(event Event) {
	h.enqueue("", event)
}

// BroadcastWorkspaceEvent marshals an Event and delivers it only to clients
// registered in the given workspace. Events carrying photo identifiers,
// filenames, or paths must go through this method (issue #74).
func (h *Hub) BroadcastWorkspaceEvent(workspaceID string, event Event) {
	if workspaceID == "" {
		// Fail closed: a tenant event without a workspace must not be
		// broadcast globally — that is precisely the leak this guards.
		slog.Warn("dropping workspace websocket event with empty workspace id", "type", event.Type)
		return
	}
	h.enqueue(workspaceID, event)
}

func (h *Hub) enqueue(workspaceID string, event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal websocket event", "type", event.Type, "error", err)
		return
	}
	msg := outbound{workspaceID: workspaceID, data: data}
	select {
	case h.broadcast <- msg:
	default:
		// Broadcast channel full; drop oldest
		select {
		case <-h.broadcast:
		default:
		}
		h.broadcast <- msg
	}
}

// ClientCount returns the number of currently connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// WorkspaceClientCount returns the number of clients connected in a
// workspace (approximate, for logging/metrics).
func (h *Hub) WorkspaceClientCount(workspaceID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	n := 0
	for client := range h.clients {
		if client.workspaceID == workspaceID {
			n++
		}
	}
	return n
}
