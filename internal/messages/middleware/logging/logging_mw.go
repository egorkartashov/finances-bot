package logging

import (
	"context"

	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"go.uber.org/zap"
)

func Middleware(handler messages.MessageHandler) messages.MessageHandler {
	return messages.NewHandleFunc(
		func(ctx context.Context, msg messages.Message) messages.HandleResult {
			logger.Info(
				"received message",
				zap.Int64("userId", msg.UserID),
				zap.String("text", msg.Text),
				zap.String("callback", msg.CallbackData),
			)

			res := handler.Handle(ctx, msg)

			if res.Err != nil {
				logError(res)
			}

			return res
		},
	)
}

func logError(res messages.HandleResult) {
	if res.HandlerName != "" {
		logger.Info(
			"error processing message:",
			zap.Error(res.Err), zap.String("handlerName", res.HandlerName),
		)
	} else {
		logger.Info("error processing message:", zap.Error(res.Err))
	}
}
