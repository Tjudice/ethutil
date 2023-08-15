package sprint

import "context"

type ProgressManager interface {
	// Should get the last succesful block that sprint has processed. The last successful block should be the highest block height
	// sprint has processed and completed, and not the most recent block that sprint has ran.
	GetLastSuccessfulBlock(ctx context.Context) (int64, error)
}

type ChainManager interface {
	// Should give the current chain head block. This is typically a wrapped for a call to eth_blockNumber.
	// This method tells sprint what block to schedule up to. This method can be memoized to reduce RPC request frequency
	LiveBlock(ctx context.Context) (int64, error)
}
