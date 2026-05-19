package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	ws "gitlab.com/r1chjames/photobox/api/internal/components/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // CORS is handled by Gin middleware
	},
}

// WebSocketHandler handles upgrading HTTP connections to WebSocket.
type WebSocketHandler struct {
	hub *ws.Hub
}

// NewWebSocketHandler creates a new WebSocketHandler.
func NewWebSocketHandler(hub *ws.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

// HandleUpgrade upgrades the HTTP connection to a WebSocket and registers
// the client with the hub. The connection is kept alive for the duration
// of the client's session via read/write pumps.
func (h *WebSocketHandler) HandleUpgrade(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Warn("websocket upgrade failed", "error", err)
		return
	}

	h.hub.Register(buildClient(h.hub, conn))

	slog.Debug("websocket connection upgraded", "remote", c.Request.RemoteAddr)
}

// buildClient creates a Client with the proper internals set.
// This avoids exporting the Client struct fields while keeping
// the package boundary clean.
func buildClient(hub *ws.Hub, conn *websocket.Conn) *ws.Client {
	return ws.NewClient(hub, conn)
}
