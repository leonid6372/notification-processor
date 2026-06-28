package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/leonid6372/notification-processor/internal/domains"
	"github.com/leonid6372/notification-processor/pkg/errs"
)

type Notification struct {
	SendAt     time.Time       `db:"send_at"`
	Payload    *Payload        `db:"-"`
	StartedAt  *sql.NullString `db:"started_at"`
	Type       string          `db:"type"`
	Status     string          `db:"status"`
	RawPayload json.RawMessage `db:"payload"`
	UserID     int             `db:"user_id"`
	TriesCount int             `db:"tries_count"`
	ID         uuid.UUID       `db:"id"`
}

type Payload struct {
	OrderID         *uuid.UUID `json:"order_id"`
	PaymentID       *uuid.UUID `json:"payment_id"`
	ShippingAddress *string    `json:"shipping_address"`
}

func (n *Notification) ToDomain() (*domains.Notification, error) {
	payload := new(domains.Payload)

	if n.RawPayload != nil {
		if err := json.Unmarshal(n.RawPayload, &n.Payload); err != nil {
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
