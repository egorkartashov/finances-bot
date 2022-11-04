package messages

import "context"

type handlerFunc struct {
	f func(ctx context.Context, msg Message) HandleResult
}

func (h *handlerFunc) Handle(ctx context.Context, msg Message) HandleResult {
	return h.f(ctx, msg)
}

func NewHandleFunc(f func(ctx context.Context, msg Message) HandleResult) MessageHandler {
	return &handlerFunc{f}
}
