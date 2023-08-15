package sprint

import (
	"context"
	"sync/atomic"
)

type Sprinter interface {
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
	collectors []Collector
}

func (s *Sprint) Run(ctx context.Context) {
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
