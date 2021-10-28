package sender

import (
	"context"

	"github.com/Damon-V79/act-transition-api/internal/model"
)

type EventSender interface {
	Send(ctx context.Context, transition *model.TransitionEvent) error
}
