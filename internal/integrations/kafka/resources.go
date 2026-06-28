package kafka

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Notification struct {
	EventID   uuid.UUID       `json:"event_id"`
	UserID    int             `json:"user_id"`
	EventType string          `json:"event_type"`
	Timestamp string          `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}
