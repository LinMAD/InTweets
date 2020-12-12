package domain

// WebSocketEvent used for communication
type WebSocketEvent struct {
	// Data from client
	Data string `json:"data,omitempty"`
}
