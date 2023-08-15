package sprint

type Status struct {
	// Whether or not sprint is actively collecting blocks
	Active bool
	// The current sprint block height, which is the smallest block height of the collector with
	BlockHeight int64
	// The statuses of each of the collectors that sprint is managing
	// TODO: Collector status type
	CollectorStatuses []Status
}
