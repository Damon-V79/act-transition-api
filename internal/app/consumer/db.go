package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/Damon-V79/act-transition-api/internal/app/repo"
	"github.com/Damon-V79/act-transition-api/internal/model"
)

type Consumer interface {
	Start(ctx context.Context)
	Close()
}

type consumer struct {
	n      uint64
	events chan<- []model.TransitionEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	done chan bool
	wg   *sync.WaitGroup
}

type Config struct {
	n         uint64
	events    chan<- []model.TransitionEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- []model.TransitionEvent) Consumer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
		done:      done,
	}
}

func (c *consumer) Start(ctx context.Context) {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(ctx, c.batchSize)
					if err != nil {
						continue
					}
					c.events <- events
				case <-c.done:
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	close(c.done)
	c.wg.Wait()
}
