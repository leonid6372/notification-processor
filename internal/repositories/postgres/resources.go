package postgres

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/leonid6372/notification-processor/internal/domains"
	"github.com/leonid6372/notification-processor/pkg/errs"
)

type Notification struct {
	ID         uuid.UUID        `db:"id"`
	UserID     uuid.UUID        `db:"user_id"`
	Type       string           `db:"type"`
	Payload    *Payload         `db:"-"`
	RawPayload *json.RawMessage `db:"payload"`
	SendAt     time.Time        `db:"send_at"`
	Status     string           `db:"status"`
	TriesCount int              `db:"tries_count"`
	StartedAt  string           `db:"started_at"`
}

type Payload struct {
	OrderID         *uuid.UUID `json:"order_id"`
	PaymentID       *uuid.UUID `json:"payment_id"`
	ShippingAddress *string    `json:"shipping_address"`
}

func (n *Notification) ToDomain() (*domains.Notification, error) {
	payload := new(domains.Payload)

	if n.Payload != nil {
		if err := json.Unmarshal(*n.RawPayload, &n.Payload); err != nil {
			return nil, errs.NewStack(err)
		}

		if n.Payload.OrderID != nil {
			payload.OrderID = *n.Payload.OrderID
		}

		if n.Payload.PaymentID != nil {
			payload.PaymentID = *n.Payload.PaymentID
		}

		if n.Payload.ShippingAddress != nil {
			payload.ShippingAddress = *n.Payload.ShippingAddress
		}
	}

	return &domains.Notification{
		ID:      n.ID,
		UserID:  n.UserID,
		Type:    n.Type,
		Payload: payload,
		SendAt:  n.SendAt,
		Status:  n.Status,
	}, nil

}
