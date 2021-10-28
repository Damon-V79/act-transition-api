package producer

import (
	"context"
	"sync"
	"time"

	"github.com/Damon-V79/act-transition-api/internal/app/repo"
	"github.com/Damon-V79/act-transition-api/internal/app/sender"
	"github.com/Damon-V79/act-transition-api/internal/model"

	"github.com/gammazero/workerpool"
)

type Producer interface {
	Start(ctx context.Context)
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan []model.TransitionEvent

	cleaner repo.Cleaner
	updater repo.Updater

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool
}

// todo for students: add repo
func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan []model.TransitionEvent,
	cleaner repo.Cleaner,
	updater repo.Updater,
	workerPool *workerpool.WorkerPool,
) Producer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &producer{
		n:          n,
		sender:     sender,
		events:     events,
		cleaner:    cleaner,
		updater:    updater,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
	}
}

func (p *producer) Start(ctx context.Context) {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case events := <-p.events:
					var success []uint64
					var failed []uint64
					for _, event := range events {
						if err := p.sender.Send(ctx, &event); err != nil {
							success = append(success, event.ID)
						} else {
							failed = append(failed, event.ID)
						}
					}
					if len(success) > 0 {
						p.workerPool.Submit(func() {
							p.cleaner.Remove(ctx, success)
						})
					}
					if len(failed) > 0 {
						p.workerPool.Submit(func() {
							p.updater.Unlock(ctx, failed)
						})
					}
				case <-p.done:
					return
				}
			}
		}()
	}
}

func (p *producer) Close() {
	close(p.done)
	p.wg.Wait()
}
