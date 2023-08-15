package v2

import "time"

type BatchLog interface {
	SetMetadata(updated BatchMetadata)
	BatchLogViewer
}

type BatchLogViewer interface {
	Key() string
	BlockStart() int64
	BlockEnd() int64
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

type Indexable interface {
	Index() uint64
}

type Event struct{}
