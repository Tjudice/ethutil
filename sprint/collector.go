package v2

import "context"

type Call struct{}

type Collector interface {
	BeforeRange(ctx context.Context, log BatchLog, stage StageInteractor) error
	AfterRange(ctx context.Context, log BatchLog, stage StageInteractor) error
}

type SprintStage struct {
	scheduleLog BatchLog
	Calls       []Call
}

type StageInteractor interface {
	// AddCall
}
