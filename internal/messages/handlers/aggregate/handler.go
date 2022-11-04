package aggregate

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages/handlers/utils"
)

type Aggregate struct {
	handlers []messages.MessageHandler
}

func (a *Aggregate) Handle(ctx context.Context, msg messages.Message) messages.HandleResult {
	for _, h := range a.handlers {
		res := h.Handle(ctx, msg)
		if res.Skipped {
			continue
		}
		if namedHandler, ok := h.(messages.NamedHandler); ok {
			res.HandlerName = namedHandler.Name()
		}
		return res
	}

	return utils.HandleSkipped
}

func NewAggregate(handlers ...messages.MessageHandler) *Aggregate {
	return &Aggregate{
		handlers: handlers,
	}
}
