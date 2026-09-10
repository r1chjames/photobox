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
//
// The connection is bound to the caller's resolved workspace: the hub only
// delivers that workspace's events to it (issue #74). A request without a
// workspace context is rejected before the upgrade — fail closed.
func (h *WebSocketHandler) HandleUpgrade(c *gin.Context) {
	wsCtx := GetWorkspaceContext(c)
	if wsCtx == nil || wsCtx.WorkspaceID == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "No workspace context"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Warn("websocket upgrade failed", "error", err)
		return
	}

	h.hub.Register(ws.NewClient(h.hub, conn, wsCtx.WorkspaceID))

	slog.Debug("websocket connection upgraded", "remote", c.Request.RemoteAddr, "workspace", wsCtx.WorkspaceID)
}
