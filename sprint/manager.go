package sprint

import "context"

type Manager interface {
	Scheduler
	Storage
	Chain
}

type Scheduler interface {
	// Should get the last succesful block that sprint has processed. The last successful block should be the highest block height
	// sprint has processed and completed, and not the most recent block that sprint has ran.
	GetTaskCompletedProgress(ctx context.Context, id TaskID) (int64, error)
	// Gets the last scheduled block for a given task. This is the last block for a given task ID that sprint has scheduled a task for
	GetTaskScheduledProgress(ctx context.Context, id TaskID) (int64, error)
	// Creates a new batch log with given ID, blockStart and blockEnd. The blockStart, blockEnd and ID are chosen by sprint.
	// Subsequent calls to GetTaskScheduledProgress should reflect the new job inserted by InsertNewBatchJob
	InsertBatchJob(id TaskID, blockStart, blockEnd int64) BatchLog
}

type Storage interface {
	// Inserts the collect information for a given block range. Each information should be inserted in a transaction, if using a database, or
	// some form of isolated, concurrency safe write. Once the batch is inserted successfully, the progress table will be updated and this block
	// range will be considered successful
	InsertBatch(ctx context.Context, progressLog BatchLog, batch BatchData) error
}

type Chain interface {
	// Should give the current chain head block. This is typically a wrapped for a call to eth_blockNumber.
	// This method tells sprint what block to schedule up to. This method can be memoized to reduce RPC request frequency
	LiveBlock(ctx context.Context) (int64, error)
}
