package sprint

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Sprinter interface {
	AddCollector(startBlock int64, c Collector) error
	Run(ctx context.Context)
	Abort()
	Status() Status
}

func NewSprint(config Config, m Manager, collectors []Collector) Sprinter {
	return &Sprint{
		config:  config,
		manager: m,
	}
}

type Sprint struct {
	// Signals sprint is done
	done chan struct{}
	// Tells us if sprint is active or not
	active atomic.Bool
	// Configuration Information
	config Config
	// Manages the sprint fetcher
	manager Manager
	// Collectors, used for getting information in between a given block range
	collectors map[TaskID]*collectorInfo
	// Mutex
	m sync.Mutex
}

func (s *Sprint) AddCollector(startBlock int64, c Collector) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.collectors[c.ID()]; ok {
		return &ErrDuplicateCollectorID{c.ID()}
	}
	s.collectors[c.ID()] = &collectorInfo{
		Status: CollectorStatus{
			HeadBlock: startBlock - 1,
		},
		collector: c,
	}
	return nil
}

func (s *Sprint) Run(ctx context.Context) {
	tickSchedule := time.NewTicker(s.config.ScheduleInterval)
	tickExecute := time.NewTicker(s.config.ExecuteInterval)
	for {
		select {
		case <-s.done:
			log.Println("Done signal received, exiting")
			return
		case <-ctx.Done():
			log.Println("Context cancelled, exiting")
			return
		case <-tickSchedule.C:
			err := s.scheduleNewCollectorRanges(ctx)
			if err != nil {
				log.Printf("Error while scheduling new collection tasks: %s\n", err.Error())
			}
		case <-tickExecute.C:
		}
	}
}

func (s *Sprint) Abort() {
	s.done <- struct{}{}
}

func (s *Sprint) Status() Status {
	active := s.active.Load()
	return Status{
		Active: active,
	}
}
