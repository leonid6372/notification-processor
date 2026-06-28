package sender

import (
	"github.com/leonid6372/notification-processor/pkg/errs"
)

func (s *Sender) SendEmail(userID int, title, text string) error {
	return errs.NewStack(emulateError(s.r))
}
