package sender

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/leonid6372/notification-processor/pkg/errs"
)

var (
	errTemporary = errors.New("temporary error")
	errPermament = errors.New("permament error")
)

func emulateError(r *rand.Rand) error {
	chance := r.Float64()

	// 1% permament error [от 0 до 0.01)
	if chance < 0.01 {
		return errs.NewStack(errPermament)
	}

	// 10% temporary error [0.01; 0.11)
	if chance < 0.11 {
		return errs.NewStack(errTemporary)
	}

	// 89% - success
	return nil
}

func (s *Sender) doWithRetry(f func() error) error {
	for attempt := 1; attempt <= s.retryCount; attempt++ {
		err := f()
		if err == nil {
			break
		}

		if attempt == s.retryCount {
			return errs.NewStack(err)
		}

		if errors.Is(err, errPermament) {
			continue
		}

		if errors.Is(err, errTemporary) {
			select {
			case <-s.ctx.Done():
				return s.ctx.Err()

			case <-time.After(calculateDelay(s.minDelay, attempt)):
				continue
			}
		}
	}

	return nil
}

func calculateDelay(minDelay time.Duration, attempt int) time.Duration {
	delay := float64(minDelay) * math.Pow(2.0, float64(attempt-1))

	jitter := rand.Float64() * 0.05 * delay // up to 100 ms

	delay += jitter

	return time.Duration(delay)
}
