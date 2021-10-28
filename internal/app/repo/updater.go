package repo

import "context"

type Updater interface {
	Unlock(ctx context.Context, eventIDs []uint64) error
}

type updater struct {
	repo EventRepo
}

func NewDbUpdater(repo EventRepo) Updater {
	return &updater{
		repo: repo,
	}
}

func (u *updater) Unlock(ctx context.Context, eventIDs []uint64) error {
	return u.repo.Unlock(ctx, eventIDs)
}
