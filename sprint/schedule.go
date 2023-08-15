package sprint

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tjudice/util/go/generic"
)

func (s *Sprint) schedulePendingTasks(ctx context.Context, id TaskID, newBlock int64) error {
	collectorInfo, ok := s.collectors[id]
	if !ok {
		return fmt.Errorf("Error: collector with id %s not found", id)
	}
	lastBlock, err := s.manager.GetTaskScheduledProgress(ctx, id)
	if err != nil {
		return err
	}
	nextRangeStart := lastBlock + 1
	if nextRangeStart > newBlock {
		return nil
	}
	rangesToSchedule := generic.DivideRangeInclusive(nextRangeStart, newBlock, s.config.BlocksPerStage)
	for _, rng := range rangesToSchedule {
		s.manager.InsertBatchJob(id, rng.Start, rng.End)
	}
	// TODO: Implement collector info update
	_ = collectorInfo
	return nil
}

func (s *Sprint) scheduleNewCollectorRanges(ctx context.Context) error {
	cctx, cn := context.WithTimeout(ctx, 30*time.Second)
	defer cn()
	blockToScheduleWith, err := s.manager.LiveBlock(cctx)
	if err != nil {
		return fmt.Errorf("Error while getting live block: %s", err.Error())
	}
	if blockToScheduleWith <= 0 {
		return fmt.Errorf("Error while scheduling: Received %d", blockToScheduleWith)
	}
	s.m.Lock()
	defer s.m.Unlock()
	for id := range s.collectors {
		err := s.schedulePendingTasks(cctx, id, blockToScheduleWith)
		if err != nil {
			log.Printf("Error while scheduling new tasks for collector id %s: %s\n", id, err.Error())
		}
	}
	return nil
}
