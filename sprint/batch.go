package sprint

import "time"

// A Batch log is a interface that defines a stage meant for processing by sprint. It contains data related to the block range
// of the task, and metadata such as the start time, end time, and any additional data that should be stored.
type BatchLog interface {
	// Sets the metadata for the given batch log. Sprint will call this function with task execution metadata that is logged
	// throughout the collection process. It is up to the implementation if this data is stored or not.
	SetMetadata(updated BatchMetadata)
	// View functions. Batch logs will be passed from sprint to the manager implementation when the log information is meant not to be modified
	BatchLogViewer
}

type TaskID string

type BatchLogViewer interface {
	// The Unique Task ID of the given batch. This is passed from the Collector interface the task is set to executee.
	ID() TaskID
	// The start and end block the given bactch log is executing
	Range() (startBlock, endBlock int64)
	// Task Metadata
	Metadata() BatchMetadata
}

type BatchMetadata struct {
	StartTime time.Time
	EndTime   time.Time
	Msg       string
	ExtraData map[string]any
}

type BatchData struct {
	CallData  []interface{}
	EventLogs []Event
}
