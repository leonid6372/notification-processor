package kafka

import "github.com/google/uuid"

type Notification struct {
	EventID   uuid.UUID `json:"event_id"`
	UserID    uuid.UUID `json:"user_uuid"`
	EventType string    `json:"event_type"`
	Timestamp string    `json:"timestamp"`
	Payload   []byte    `json:"payload"`
}
