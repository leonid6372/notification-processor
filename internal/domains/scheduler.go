package domains

import "context"

type Scheduler interface {
	StartSendings(ctx context.Context)
	WaitToStopSending()
	StartCleaning(ctx context.Context)
}
