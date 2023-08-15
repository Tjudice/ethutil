package sprint

type ErrDuplicateCollectorID struct {
	id TaskID
}

func (e *ErrDuplicateCollectorID) Error() string {
	return "Error: Duplicate Collector Id: " + string(e.id)
}
