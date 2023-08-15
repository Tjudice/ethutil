package sprint

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/upper/db/v4"
	"tuxpa.in/a/zlog/log"
)

type TaskStatus int

var (
	StatusScheduled   TaskStatus = -1
	StatusUnscheduled TaskStatus = 0
	StatusFinished    TaskStatus = 1
)

type StageProgressLog struct {
	PrimaryKey      string     `db:"primary_key" json:"_key"`
	Stage           int        `db:"stage" json:"stage"`
	StartBlock      int        `db:"start_block" json:"start_block"`
	EndBlock        int        `db:"end_block" json:"end_block"`
	StartTs         time.Time  `db:"start_ts" json:"start_ts"`
	EndTs           time.Time  `db:"end_ts" json:"end_ts"`
	Success         TaskStatus `db:"success" json:"success"`
	Error           string     `db:"error" json:"error"`
	Msg             string     `db:"msg" json:"msg"`
	ValidatorPasses int        `db:"validator_passes" json:"validator_passes"`
}

func (s *Sprint) schedule(ctx context.Context) error {
	scheduleHead, err := s.currentScheduleHead(ctx)
	if err != nil {
		return err
	}
	chainHead := s.m.CurrentBlock()
	if chainHead < s.c.StartBlock || scheduleHead < s.c.StartBlock {
		return fmt.Errorf("sanity check: live block is too low")
	}
	scheduleHead = scheduleHead + 1
	rngs := divideStageRanges(scheduleHead, chainHead, s.c.BlocksPerStage)
	toSchedule, err := createScheduleLogs(rngs)
	if err != nil {
		return err
	}
	if len(toSchedule) == 0 {
		return nil
	}
	return s.insertScheduledRanges(ctx, toSchedule)
}

func (s *Sprint) insertScheduledRanges(ctx context.Context, toSchedule []*StageProgressLog) error {
	err := s.db.TxContext(ctx, func(d db.Session) error {
		b := s.db.SQL().InsertInto(s.progressTable).Amend(onConflictDoNothing).Batch(1000)
		go func() {
			defer b.Done()
			for _, next := range toSchedule {
				log.Info().Int("start_block", next.StartBlock).Int("end_block", next.EndBlock).Msg("scheduled stage")
				b.Values(next)
			}
		}()
		return b.Wait()
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sprint) currentScheduleHead(ctx context.Context) (uint64, error) {
	type Blk struct {
		EndBlock uint64 `db:"end_block"`
	}
	b := &Blk{}
	err := s.db.SQL().Select("end_block").From(s.progressTable).OrderBy("-end_block").Limit(1).One(b)
	if err != nil {
		if strings.Contains(err.Error(), "no more rows in this result set") {
			return s.c.StartBlock, nil
		}
		return 0, err
	}
	return b.EndBlock, nil
}

func (s *Sprint) getNextRangeToExecute(ctx context.Context) (*StageProgressLog, error) {
	prog := &StageProgressLog{}
	err := s.db.TxContext(ctx, func(sess db.Session) error {
		// should probs be able to configure limit
		err := sess.SQL().SelectFrom(s.progressTable).Where("success = 0 AND stage = 2").OrderBy("start_block ASC").Limit(s.c.Workers).One(prog)
		if err != nil {
			return err
		}
		prog.StartTs = time.Now()
		prog.Success = StatusScheduled
		return sess.Collection(s.progressTable).UpdateReturning(prog)
	}, nil)
	return prog, err
}

func (s *Sprint) getCurrentStartHeadBlock(ctx context.Context) (int, error) {
	prog := &StageProgressLog{}
	err := s.db.SQL().SelectFrom(s.progressTable).Where("success >= 1 AND stage = 2").OrderBy("-end_block").Limit(1).One(prog)
	if err != nil {
		if strings.Contains(err.Error(), "no more rows in this result set") {
			return int(s.c.StartBlock), nil
		}
		return 0, err
	}
	return prog.StartBlock, nil
}

type ValidatorBlock struct {
	StartBlock       int `db:"start_block"`
	ValidatorPassess int `db:"validator_passes"`
}

func (s *Sprint) getValidatorTaskRanges(ctx context.Context, currHeadStart, validatorCount, validatorSpacing int) ([]*StageProgressLog, error) {
	if validatorCount == 0 {
		return nil, fmt.Errorf("validator count is 0")
	}
	// for first validator, we want to start at the most-validated block + 1
	// for all other validators, we want to start at the second most-validated block - validatorSpacing
	res := make([]*StageProgressLog, 0, validatorCount)
	lastScheduledValidatorBlock := currHeadStart - validatorSpacing
	s.db.TxContext(ctx, func(sess db.Session) error {
		var maxes []*ValidatorBlock
		err := sess.SQL().Select("validator_passes", db.Raw("MIN(start_block) as start_block")).
			From(s.progressTable).Where("success = 1 AND validator_passes >= 0 AND validator_passes < ?", validatorCount).
			GroupBy("validator_passes").OrderBy("validator_passes").
			Limit(validatorCount).
			All(&maxes)
		if err != nil {
			return err
		}
		for _, max := range maxes {
			if max.ValidatorPassess >= validatorCount {
				continue
			}
			currValidatorTask := &StageProgressLog{}
			err := s.db.SQL().SelectFrom(s.progressTable).
				Where("success = 1 AND validator_passes <= ? AND validator_passes >= 0 AND start_block >= ? AND end_block <= ?",
					max.ValidatorPassess, max.StartBlock, lastScheduledValidatorBlock).
				OrderBy("-validator_passes", "end_block").Limit(1).One(currValidatorTask)
			if err != nil {
				continue
			}
			lastScheduledValidatorBlock = currValidatorTask.StartBlock - validatorSpacing
			currValidatorTask.ValidatorPasses = -currValidatorTask.ValidatorPasses - 1
			currValidatorTask.StartTs = time.Now()
			err = s.db.Collection(s.progressTable).UpdateReturning(currValidatorTask)
			if err != nil {
				return err
			}
			res = append(res, currValidatorTask)
		}
		return nil
	}, nil)
	reversed := make([]*StageProgressLog, len(res))
	for i := range res {
		reversed[len(res)-1-i] = res[i]
	}
	return reversed, nil
}

// Resets the success of tasks that were terminated while running
func (s *Sprint) resetProgressSuccess(ctx context.Context) error {
	_, err := s.db.SQL().Update(s.progressTable).Set("success", StatusUnscheduled).Where("success = -1 AND stage = 2").ExecContext(ctx)
	if err != nil {
		return err
	}
	_, err = s.db.SQL().Update(s.progressTable).Set("validator_passes", db.Raw("-validator_passes - 1")).Where("validator_passes < 0").ExecContext(ctx)
	return err
}

func createScheduleLogs(stageRanges []stageFilterRange) ([]*StageProgressLog, error) {
	toSchedule := make([]*StageProgressLog, 0, len(stageRanges))
	for _, nextRange := range stageRanges {
		toSchedule = append(toSchedule, &StageProgressLog{
			PrimaryKey: fmt.Sprintf("%d-%d", nextRange.start, nextRange.end),
			Stage:      2,
			StartBlock: int(nextRange.start),
			EndBlock:   int(nextRange.end),
			StartTs:    time.Time{},
			EndTs:      time.Time{},
			Error:      "",
			Msg:        fmt.Sprintf("Scheduled at %s", time.Now().String()),
		})
	}
	return toSchedule, nil
}

type stageFilterRange struct {
	start uint64
	end   uint64
}

func divideStageRanges(start, stop, size uint64) []stageFilterRange {
	if size == 0 || start == stop {
		return []stageFilterRange{{start: start, end: stop}}
	}
	if start > stop {
		return nil
	}
	totalRange := stop - start
	res := make([]stageFilterRange, 0, totalRange/size+1)
	for i := start; i < stop; i += (size + 1) {
		end := i + size
		if end > stop {
			end = stop
		}
		res = append(res, stageFilterRange{start: i, end: end})
	}
	return res
}
