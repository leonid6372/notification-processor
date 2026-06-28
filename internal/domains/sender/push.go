package sender

import (
	"github.com/google/uuid"

	"github.com/leonid6372/notification-processor/pkg/errs"
)

func (s *Sender) SendPush(userID uuid.UUID, title, text string) error {
	return errs.NewStack(emulateError(s.r))
}
