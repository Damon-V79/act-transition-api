package repo

import (
	"context"

	"github.com/Damon-V79/act-transition-api/internal/model"
)

type EventRepo interface {
	Lock(ctx context.Context, n uint64) ([]model.TransitionEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error

	Add(ctx context.Context, event []model.TransitionEvent) error
	Remove(ctx context.Context, eventIDs []uint64) error
}
