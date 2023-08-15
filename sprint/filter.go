package sprint

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EventFilter interface {
	// Populates a given type with the information from an event log. Returns a populated event,
	// a boolean of whether or not this event should be ignored, and an error which will trigger
	// a stage repopulation if it is non-nil.
	Populate(log types.Log) (populated Event, ok bool, err error)
	// Allows for hooks to be added when new events occur. This is useful for collecting information that
	// may change when a given event is emitted by a contract.
	// Users can customize the implementation of the Call to allow for processing of specific blocks
	Hook(Event) []Call
	// Hash of the event. This will be used to add filters to event logs for a given collector. Should always return a valid event hash
	EventHash() common.Hash
}
