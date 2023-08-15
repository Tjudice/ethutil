package sprint

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"gfx.cafe/open/ghost"
	"gfx.cafe/open/jrpc"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tjudice/util/go/concurrency"
	"github.com/upper/db/v4"
	"golang.org/x/sync/errgroup"
	"tuxpa.in/a/zlog/log"
)

type SprintConfig struct {
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
}

func checkConfig(c *SprintConfig) error {
	if c.BlocksPerStage <= 0 {
		return errors.New("blocks per stage must be greater than 0")
	}
	if c.Workers <= 0 {
		return errors.New("workers must be greater than 0")
	}
	if c.ExecutionQueueSize < 0 {
		return errors.New("execution queue size must be greater than or equal to 0")
	}
	if c.ScheduleInterval < 0 {
		return errors.New("schedule interval must be greater than or equal to 0")
	}
	if c.ExecuteInterval <= 0 {
		return errors.New("execute interval must be greater than 0")
	}
	if c.ValidatorSpacing <= 0 && c.ValidatorCount > 0 {
		return errors.New("validator spacing must be greater than 0")
	}
	if c.ValidatorQueueSize <= 0 && c.ValidatorCount > 0 {
		return errors.New("validator queue size must be greater than 0")
	}
	return nil
}

type SprintManager interface {
	CurrentBlock() uint64
	Insert(ctx context.Context, d db.Session, startBlock int, endBlock int, eventData *EventBatch) error
	Validate(ctx context.Context, d db.Session, startBlock int, endBlock int, eventData *EventBatch) (bool, error)
	GetStageFilters(blockStart, blockEnd int) []ethereum.FilterQuery
}

type Sprint struct {
	// task scheduling table
	progressTable string
	// manages the insertion and validation logic
	m SprintManager
	// for logging
	lastSuccessfulBlock  atomic.Int64
	validatorBlockStatus []atomic.Int64
	// clients
	rpc        jrpc.Conn
	db         db.Session
	blockCache *blockCache

	// task execution sequencing
	executionQueue *executionQueue

	// config
	c        *SprintConfig
	isActive atomic.Bool
}

var ErrSprintActive = errors.New("cannot add new event filter while sprint is active")

func NewSprint(ctx context.Context, db db.Session, rpc jrpc.Conn, config *SprintConfig, progressTable string, manager SprintManager) (*Sprint, error) {
	err := checkConfig(config)
	if err != nil {
		return nil, err
	}
	s := &Sprint{
		progressTable:        progressTable,
		db:                   db,
		rpc:                  rpc,
		c:                    config,
		isActive:             atomic.Bool{},
		blockCache:           NewBlockCache(ctx, rpc, config.Workers),
		m:                    manager,
		validatorBlockStatus: make([]atomic.Int64, config.ValidatorCount),
	}
	s.executionQueue = s.newExecutionQueue()
	return s, nil
}

func (s *Sprint) Run(ctx context.Context) error {
	s.isActive.Store(true)
	defer s.isActive.Store(false)
	// start listening for new sprint/validation tasks in background
	go s.executionQueue.run(ctx, s.c.Workers)
	// In the event the program exited while a task was queued, must reset
	// so it gets scheduled, otherwise we can skip blockranges
	err := s.resetProgressSuccess(ctx)
	if err != nil {
		return err
	}
	// set the intial last block for logging
	lastScheduledBlock, err := s.getCurrentStartHeadBlock(ctx)
	if err != nil {
		return err
	}
	s.lastSuccessfulBlock.Store(int64(lastScheduledBlock))
	go s.logLoop()
	// create tickers based on config
	scheduleTicker := time.NewTicker(s.c.ScheduleInterval)
	executeTicker := time.NewTicker(s.c.ExecuteInterval)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-scheduleTicker.C:
			// schedule new ranges when most recent block number is not contained
			// in set of scheduled block ranges
			err := s.schedule(ctx)
			if err != nil {
				log.Err(err).Msg("error scheduling stages")
			}
		case <-executeTicker.C:
			// query for next range to execute from database
			next, err := s.getNextRangeToExecute(ctx)
			// dont want to error when we don't have any new ranges to execute,
			// since we may still have validator tasks to execute
			if err == nil {
				s.executionQueue.addTask(ctx, next)
			}
			lastHeadBlock := s.lastSuccessfulBlock.Load()
			// to prevent scheduling from stalling in the event the validation queue is full
			if s.c.ValidatorCount == 0 || len(s.executionQueue.validationQueue) >= s.c.ValidatorQueueSize-s.c.ValidatorCount {
				continue
			}
			// gets validator task ranges to execute
			validatorTasks, err := s.getValidatorTaskRanges(ctx, int(lastHeadBlock), s.c.ValidatorCount, s.c.ValidatorSpacing)
			if err != nil {
				log.Err(err).Msg("error getting validator tasks")
				continue
			}
			// add validator tasks to queue
			for _, task := range validatorTasks {
				s.executionQueue.addValidationTask(ctx, task)
			}
		}
	}
}

func (s *Sprint) uploadBatch(ctx context.Context, batch *EventBatch) error {
	// execute in transaction to ensure we successfully can update and insert new events
	// and to preserve read consistency for api
	return s.db.TxContext(ctx, func(tx db.Session) error {
		// insert first so we can update end timestamp correctly
		err := s.m.Insert(ctx, tx, batch.progressLog.StartBlock, batch.progressLog.EndBlock, batch)
		if err != nil {
			return err
		}
		// update progress log values to reflect completed range
		batch.progressLog.Success = StatusFinished
		batch.progressLog.EndTs = time.Now()
		// if for some reason the progress log is deleted, this will fail
		err = tx.Collection(s.progressTable).UpdateReturning(batch.progressLog)
		if err != nil {
			return err
		}
		// log success
		s.logTaskSuccess(batch)
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (s *Sprint) validateBatch(ctx context.Context, batch *EventBatch) error {
	// execute in transaction to ensure we successfully can update and insert new events
	// and to preserve read consistency for api
	return s.db.TxContext(ctx, func(tx db.Session) error {
		// validate events in range. Validate should delete and insert correct events
		// when it detects inconsistencies
		didReorg, err := s.m.Validate(ctx, tx, batch.progressLog.StartBlock, batch.progressLog.EndBlock, batch)
		if err != nil {
			return err
		}
		// update validator passes to allow multiple validators to function correctly
		batch.progressLog.ValidatorPasses = -batch.progressLog.ValidatorPasses
		if didReorg {
			// this will be overwritten if a validator detects and fixes a reorg that is reorged out again
			batch.progressLog.Msg = fmt.Sprintf("reorg detected at time: %s after validator pass %d", time.Now().String(), batch.progressLog.ValidatorPasses)
		}
		batch.progressLog.EndTs = time.Now()
		err = tx.Collection(s.progressTable).UpdateReturning(batch.progressLog)
		if err != nil {
			return err
		}
		// log success
		s.logValidationSuccess(batch, didReorg)
		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
}

func (s *Sprint) executeTask(ctx context.Context, t *task) *EventBatch {
	var stageLogs []*ghost.ErigonLog
	var blockInfo map[int]*BlockInfo
	var err error
	for {
		stageLogs, err = s.getStageEventLogs(ctx, t.StageProgressLog.StartBlock, t.StageProgressLog.EndBlock)
		if err != nil {
			log.Err(err).Msg("error getting event logs")
			continue
		}
		blockInfo, err = s.getStageBlockInfo(ctx, stageLogs)
		if err != nil {
			log.Err(err).Msg("error getting block info")
			continue
		}
		break
	}
	return &EventBatch{
		progressLog: t.StageProgressLog,
		Events:      stageLogs,
		Blocks:      blockInfo,
	}
}

func (s *Sprint) getStageBlockInfo(ctx context.Context, eventLogs []*ghost.ErigonLog) (map[int]*BlockInfo, error) {
	blockSet := make(map[int][]common.Hash)
	for _, eventLog := range eventLogs {
		blockSet[int(eventLog.BlockNumber)] = append(blockSet[int(eventLog.BlockNumber)], eventLog.TxHash)
	}
	blocks := make(map[int]*BlockInfo, len(blockSet))
	wg := errgroup.Group{}
	m := sync.Mutex{}
	for block := range blockSet {
		block := block
		wg.Go(func() error {
			blk, err := s.blockCache.BlockAt(ctx, block)
			if err != nil {
				return err
			}
			if !checkAllTxsExist(blk, blockSet[block]) {
				return fmt.Errorf("missing transaction in block: %d", block)
			}
			m.Lock()
			defer m.Unlock()
			blocks[block] = blk
			return nil
		})
	}
	return blocks, wg.Wait()
}

func (s *Sprint) getStageEventLogs(ctx context.Context, start, stop int) ([]*ghost.ErigonLog, error) {
	filters := s.m.GetStageFilters(start, stop)
	var allLogs []*ghost.ErigonLog
	m := sync.Mutex{}
	getLogs := func(ctx context.Context, f ethereum.FilterQuery) error {
		logs, err := s.getLogs(ctx, f)
		m.Lock()
		defer m.Unlock()
		allLogs = append(allLogs, logs...)
		return err
	}
	err := concurrency.DoContext(ctx, len(filters), getLogs, filters...)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(allLogs, func(i, j int) bool {
		if allLogs[i].BlockNumber != allLogs[j].BlockNumber {
			return allLogs[i].BlockNumber < allLogs[j].BlockNumber
		}
		return allLogs[i].Index < allLogs[j].Index
	})
	return allLogs, nil
}

func (s *Sprint) getLogs(ctx context.Context, f ethereum.FilterQuery) ([]*ghost.ErigonLog, error) {
	filter := toGethFilter(f)
	var logs []*ghost.ErigonLog
	err := s.rpc.Do(ctx, &logs, "eth_getLogs", []any{filter})
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (s *Sprint) validateTask(ctx context.Context, t *task) *EventBatch {
	var stageLogs []*ghost.ErigonLog
	var blockInfo map[int]*BlockInfo
	var err error
	for {
		stageLogs, err = s.getStageEventLogs(ctx, t.StageProgressLog.StartBlock, t.StageProgressLog.EndBlock)
		if err != nil {
			log.Err(err).Msg("error getting event logs")
			continue
		}
		blockInfo, err = s.getStageBlockInfo(ctx, stageLogs)
		if err != nil {
			log.Err(err).Msg("error getting block info")
			continue
		}
		break
	}
	return &EventBatch{
		progressLog: t.StageProgressLog,
		Events:      stageLogs,
		Blocks:      blockInfo,
	}
}

func (s *Sprint) updateValidatorProgress(validatorID int, endBlock int) {
	if validatorID < 1 || validatorID > len(s.validatorBlockStatus) {
		log.Error().Int("validator_index", validatorID).Msg("invalid validator index")
		return
	}
	s.validatorBlockStatus[validatorID-1].Store(int64(endBlock))
}
