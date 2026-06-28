package sender

import (
	"github.com/google/uuid"

	"github.com/leonid6372/notification-processor/pkg/errs"
)

func (s *Sender) SendSMS(userID uuid.UUID, title, text string) error {
	return errs.NewStack(emulateError(s.r))
}
