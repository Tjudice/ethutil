package sprint

import "time"

type Config struct {
	// The number of blocks per individual stage. Setting this number too high
	// can cause RPC issues due to the number of events being returned. Balancing this
	// along with the number of concurrent workers is important.
	BlocksPerStage uint64
	// The number of concurrent workers to use when fetching data. This number should
	// If using a rate-limited RPC, set this slightly below the rate limit, since
	// the scheduler must request the current block number.
	Workers int
	// This sets the maximum queue size for the task execution queue. Too large a number can cause large memory usage, whereas too small
	// can cause the validator to be unable to keep up with the sprint.
	ExecutionQueueSize int
	// The interval between when the scheduler should check if it needs to schedule a new task.
	// This number should correspond directly to the block time and should be reduced
	// if lower latency is desired.
	ScheduleInterval time.Duration
	// The time between when the scheduler finishes executing a task and when it should
	// attempt to execute the next task. Setting this to a high value increases latency.
	ExecuteInterval time.Duration
	// Block to start filtering add. It is reccomended to set this to the first block
	// that any of the contracts were deployed.
	StartBlock uint64
	// The number of validators that should be spawned to verify past event ranges to ensure reorg saftey.
	ValidatorCount int
	// The spacing of the validators. First validator starts at head - ValidatorSpacing, and each successive validator will run
	// at head - (ValidatorSpacing * i)
	ValidatorSpacing int
	// This sets the maximum queue size for the validator queue. Too large a number can cause large memory usage, whereas too small
	// can cause the validator to be unable to keep up with the sprint.
	ValidatorQueueSize int
	// Sets verbosity
	Verbose bool
	// Disables automatic killswitch
	DisableKillSwitch bool
}
