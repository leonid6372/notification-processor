package domain

import (
	"github.com/google/uuid"
)

const (
	NotificationTypeOrderCreated    = "order_created"
	NotificationTypePaymentReceived = "payment_received"
	NotificationTypeOrderShipped    = "order_shipped"
)

type NotificationsRepository interface {
}

type Notification struct {
	ID uuid.UUID `json:"id"`
}
