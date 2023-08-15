package sprint

import (
	"fmt"
	"time"

	"tuxpa.in/a/zlog"
	"tuxpa.in/a/zlog/log"
)

// loops and logs status of sprint and validator every 5 seconds
func (s *Sprint) logLoop() {
	if !s.c.Verbose {
		return
	}
	logInterval := time.NewTicker(5 * time.Second)
	for {
		<-logInterval.C
		if !s.c.Verbose {
			continue
		}
		statusLog := log.Debug().Int64("chain_head", int64(s.m.CurrentBlock())).
			Int64("sprint_head", s.lastSuccessfulBlock.Load())
		for i := range s.validatorBlockStatus {
			statusLog = statusLog.Int64(fmt.Sprintf("validator_%d", i), s.validatorBlockStatus[i].Load())
		}
		statusLog.Int("task_queue_size", len(s.executionQueue.taskQueue)+len(s.executionQueue.uploadSequence)).
			Int("validation_queue_size", len(s.executionQueue.validationQueue)+len(s.executionQueue.validationSequence)).
			Msg("status")
	}
}

func (s *Sprint) logTaskSuccess(b *EventBatch) {
	if !s.c.Verbose {
		return
	}
	getProgressLog(b).Msg("range complete")
}

func (s *Sprint) logValidationSuccess(b *EventBatch, reorgDetected bool) {
	if !s.c.Verbose {
		return
	}
	baseInfo := getProgressLog(b)
	baseInfo.Bool("reorg_detected", reorgDetected).
		Msg("validation complete")
}

func getProgressLog(b *EventBatch) *zlog.Event {
	prog := b.progressLog
	logInfo := log.Info().Int("start_block", prog.StartBlock).
		Int("end_block", prog.EndBlock).
		Float64("duration (s)", prog.EndTs.Sub(prog.StartTs).Seconds()).
		Float64("block_time (s)", calculateBlockTime(prog.StartTs, prog.EndTs, prog.StartBlock, prog.EndBlock)).
		Int("event_count", len(b.Events))
	return logInfo
}

func calculateBlockTime(start, end time.Time, startBlock, endBlock int) float64 {
	// add 1 since block ranges are inclusive
	return end.Sub(start).Seconds() / float64(endBlock-startBlock+1)
}
