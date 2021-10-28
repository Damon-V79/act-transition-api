package repo

import "context"

type Cleaner interface {
	Remove(ctx context.Context, eventIDs []uint64) error
}

type cleaner struct {
	repo EventRepo
}

func NewDbCleaner(repo EventRepo) Cleaner {
	return &cleaner{
		repo: repo,
	}
}

func (c *cleaner) Remove(ctx context.Context, eventIDs []uint64) error {
	return c.repo.Remove(ctx, eventIDs)
}
