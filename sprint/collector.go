package sprint

import "context"

type Collector interface {
	// The ID for the given collector. This should be unique, so that the manager can properly request up to date information for the given
	// collector ID. If a new item is added to the collection scheme after the collector has progressed from the start block,
	// it will never be repopulated unless the entire range is reset or validated.
	ID() TaskID
	// This is a hook to run before the given block range is run. This is meant to be used to preprocess a range to determine what
	// calls should be added/modified for the specified block range. BeforeRange hooks only effect the currently executing task range, and
	// do not persist across seperate tasks.
	BeforeRange(ctx context.Context, log BatchLog, stage StageInteractor) error
	// This is a hook that runs after a block range has been ran. This allows users to utilize data collected in a given range to either
	// postprocess data, update a cache, delete unwanted data, etc...
	AfterRange(ctx context.Context, log BatchLog, stage StageInteractor, collected StageData) error
	// Gets the log filters for a given collector. These filters should not be modified after the collector has already began progressing,
	// otherwise the stored state may become inconsistent
	GetLogFilters() []EventFilter
	// Gets the calls to be made by this collector. Calls are ran at each block in the specified range. It is up to the implementation to filter calls
	// if it is unneccsary to request data at every block in a task range
	GetCalls() []Call
}

// A call is any function (specifically one that requests data at a block) that a collector should process for every block range
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

type StageData interface{}
