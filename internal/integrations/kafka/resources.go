package kafka

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Notification struct {
	EventType string          `json:"event_type"`
	Timestamp string          `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
	UserID    int             `json:"user_id"`
	EventID   uuid.UUID       `json:"event_id"`
}
