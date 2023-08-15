package sprint

import "context"

type Collector interface {
	ID() TaskID
	BeforeRange(ctx context.Context, log BatchLog, stage StageInteractor) error
	AfterRange(ctx context.Context, log BatchLog, stage StageInteractor) error
	GetLogFilters() []EventLogFilter
}

type Call func(ctx context.Context, block int64) (Keyable, error)

type Keyable interface {
	Key() string
}

type SprintStage struct {
	scheduleLog     BatchLog
	AdditionalCalls []Call
}

type StageInteractor interface {
	AddCall(c Call)
}
