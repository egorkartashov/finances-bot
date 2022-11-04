package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/messages"
	"go.uber.org/zap"
)

type cfg interface {
	ServiceName() string
}

func initTracing(cfg cfg) {
	tracingCfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := tracingCfg.InitGlobalTracer(cfg.ServiceName())
	if err != nil {
		logger.Fatal("failed to init tracer", zap.Error(err))
	}
}

func Middleware(handler messages.MessageHandler, cfg cfg) messages.MessageHandler {
	initTracing(cfg)
	return messages.NewHandleFunc(
		func(ctx context.Context, msg messages.Message) messages.HandleResult {
			span, ctx := opentracing.StartSpanFromContext(ctx, "handleMessage")
			defer span.Finish()

			res := handler.Handle(ctx, msg)
			if res.HandlerName != "" {
				span.SetTag("handler", res.HandlerName)
			}

			return res
		},
	)
}
