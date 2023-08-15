package sprint

import (
	"context"
	"time"

	"gfx.cafe/open/ghost"
	"tuxpa.in/a/zlog/log"
)

type EventBatch struct {
	progressLog *StageProgressLog
	Events      []*ghost.ErigonLog
	Blocks      map[int]*BlockInfo
}

type task struct {
	finished         chan *EventBatch
	StageProgressLog *StageProgressLog
}

type executionQueue struct {
	s               *Sprint
	taskQueue       chan *task
	validationQueue chan *task
	// double channels are required so we can fetch task block ranges concurrently but ensure we sequence and
	// upload the results of those ranges in order
	uploadSequence chan chan *EventBatch
	// validation must also be done in order since a factory event that is reorged out/in
	// may have later effects on events
	validationSequence chan chan *EventBatch
}

func (s *Sprint) newExecutionQueue() *executionQueue {
	return &executionQueue{
		s:                  s,
		taskQueue:          make(chan *task, s.c.ExecutionQueueSize),
		validationQueue:    make(chan *task, s.c.ValidatorQueueSize),
		uploadSequence:     make(chan chan *EventBatch, s.c.ExecutionQueueSize),
		validationSequence: make(chan chan *EventBatch, s.c.ValidatorQueueSize),
	}
}

func (e *executionQueue) run(ctx context.Context, numWorkers int) error {
	// adds all workers so they can process tasks
	for i := 0; i < numWorkers; i++ {
		e.addWorker(ctx)
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case validationWait := <-e.validationSequence:
				// ensure we sequence the validator results correctly. First channel wait will take in
				// the sequenced task, the second waits for the results
				batch := <-validationWait
				// we cannot validate future ranges until all past ranges are validated, so must loop here
				for {
					err := e.s.validateBatch(ctx, batch)
					if err != nil {
						log.Println(err)
						continue
					}
					break
				}
				// update sprint validator progress
				e.s.updateValidatorProgress(batch.progressLog.ValidatorPasses, batch.progressLog.EndBlock)
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case uploadWait := <-e.uploadSequence:
			// ensure we sequence the validator results correctly. First channel wait will take in
			// the sequenced task, the second waits for the results
			batch := <-uploadWait
			// we cannot insert future ranges until all past ranges are inserted
			for {
				err := e.s.uploadBatch(ctx, batch)
				if err != nil {
					log.Println(err)
					continue
				}
				break
			}
			// update sprint progress
			e.s.lastSuccessfulBlock.Store(int64(batch.progressLog.EndBlock))
		}
	}
}

func (e *executionQueue) addTask(ctx context.Context, taskParams *StageProgressLog) {
	// create response channel for upload goroutine to read on completion
	finished := make(chan *EventBatch, 1)
	// update start timestamp for logging in progress table
	taskParams.StartTs = time.Now()
	// send response channel to upload queue
	e.uploadSequence <- finished
	// send task to task queue with response channel included for sending on completion
	e.taskQueue <- &task{
		finished:         finished,
		StageProgressLog: taskParams,
	}
}

func (e *executionQueue) addValidationTask(ctx context.Context, taskParams *StageProgressLog) {
	// create response channel for validation goroutine to read on completion
	finished := make(chan *EventBatch, 1)
	// update start timestamp for logging in progress table
	taskParams.StartTs = time.Now()
	// send response channel to validation queue
	e.validationSequence <- finished
	// send task to task queue with response channel included for sending on completion
	e.validationQueue <- &task{
		finished:         finished,
		StageProgressLog: taskParams,
	}
}

func (e *executionQueue) addWorker(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case task := <-e.taskQueue:
				// must ensure tasks are sequenced correctly when scheduling
				res := e.s.executeTask(ctx, task)
				task.finished <- res
			case task := <-e.validationQueue:
				// must ensure tasks are sequenced correctly when scheduling
				res := e.s.validateTask(ctx, task)
				task.finished <- res
			}
		}
	}()
}
