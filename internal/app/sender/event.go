package sender

import (
	"github.com/Damon-V79/act-transition-api/internal/model"
)

type EventSender interface {
	Send(transition *model.TransitionEvent) error
}
