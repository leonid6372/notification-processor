package domains

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/leonid6372/notification-processor/pkg/errs"
)

const (
	NotificationTypeOrderCreated    = "order_created"
	NotificationTypePaymentReceived = "payment_received"
	NotificationTypeOrderShipped    = "order_shipped"

	NotificationStatusCreated    = "created"
	NotificationStatusInProgress = "in_progress"
	NotificationStatusCompleted  = "completed"
	NotificationStatusFailed     = "failed"
)

type NotificationsRepo interface {
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotificationsToSend(ctx context.Context, limit int) ([]*Notification, error)
	UpdateNotificationStatus(ctx context.Context, id uuid.UUID, status string) error
	ResetZombieNotifications(ctx context.Context, timeout time.Duration) error
}

type Notification struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Type       string
	Payload    *Payload
	RawPayload []byte
	SendAt     time.Time
	Status     string
}

type Payload struct {
	OrderID         uuid.UUID `json:"order_id"`
	PaymentID       uuid.UUID `json:"payment_id"`
	ShippingAddress string    `json:"shipping_address"`
}

func (n *Notification) GetTitleAndText() (string, string, error) {
	switch n.Type {
	case NotificationTypeOrderCreated:
		text := fmt.Sprintf("Your order %s has been created successfully.", n.Payload.OrderID)

		return "Order Created", text, nil
	case NotificationTypePaymentReceived:
		text := fmt.Sprintf("We have received your payment %s.", n.Payload.PaymentID)

		return "Payment Received", text, nil
	case NotificationTypeOrderShipped:
		text := fmt.Sprintf("Your order %s has been shipped in %s.", n.Payload.OrderID, n.Payload.ShippingAddress)

		return "Order Shipped", text, nil
	default:
		return "", "", errs.NewStack(fmt.Errorf("invalid notification type: %s id: %s", n.Type, n.ID))
	}
}
