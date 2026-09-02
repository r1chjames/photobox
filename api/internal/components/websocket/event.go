package websocket

// Event types broadcast to connected clients.
const (
	EventThumbnailReady = "thumbnail:ready"
	EventIndexProgress  = "index:progress"
	EventIndexComplete  = "index:complete"
	EventPhotoNew       = "photo:new"
	EventSystemHealth   = "system:health"
	EventImportProgress = "import:progress"
	EventImportComplete = "import:complete"
)

// Event is the JSON envelope for all WebSocket messages.
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// ThumbnailReadyPayload is sent when a thumbnail has been generated.
type ThumbnailReadyPayload struct {
	PhotoID string `json:"photoId"`
	Size    string `json:"size"`
}

// IndexProgressPayload is sent periodically during photo indexing.
type IndexProgressPayload struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	Phase   string `json:"phase"`
}

// IndexCompletePayload is sent when photo indexing finishes.
type IndexCompletePayload struct {
	Total  int      `json:"total"`
	Errors []string `json:"errors,omitempty"`
}

// PhotoNewPayload is sent when a new photo is detected.
type PhotoNewPayload struct {
	PhotoID string `json:"photoId"`
}

// SystemHealthPayload is sent as a periodic heartbeat.
type SystemHealthPayload struct {
	Status  string `json:"status"`
	Uptime  int64  `json:"uptime"`
	Version string `json:"version"`
}

// ImportProgressPayload is sent periodically during a Google Takeout import.
type ImportProgressPayload struct {
	Status   string `json:"status"` // "running" | "complete" | "error"
	Phase    string `json:"phase"`  // "copying" | "indexing" | "done"
	Current  int    `json:"current"`
	Total    int    `json:"total"`
	Imported int    `json:"imported"`
	Skipped  int    `json:"skipped"`
}

// ImportCompletePayload is sent when a Google Takeout import finishes.
type ImportCompletePayload struct {
	Imported int      `json:"imported"`
	Skipped  int      `json:"skipped"`
	Errors   []string `json:"errors,omitempty"`
}
