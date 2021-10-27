package repo

import (
	"github.com/Damon-V79/act-transition-api/internal/model"
)

type EventRepo interface {
	Lock(n uint64) ([]model.TransitionEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []model.TransitionEvent) error
	Remove(eventIDs []uint64) error
}
