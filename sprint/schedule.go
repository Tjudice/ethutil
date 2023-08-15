package sprint

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (s *Sprint) schedulePendingTasks(ctx context.Context, cInfo *collectorInfo, newBlock int64) error {
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
	for id, cInfo := range s.collectors {
		err := s.schedulePendingTasks(cctx, cInfo, blockToScheduleWith)
		if err != nil {
			log.Printf("Error while scheduling new tasks for collector id %s: %s\n", id, err.Error())
		}
	}
	return nil
}
