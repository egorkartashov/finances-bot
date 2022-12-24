package main

import (
	"github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/logger"
	"go.uber.org/zap"
)

func initTracing(serviceName string) {
	tracingCfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := tracingCfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("failed to init tracer", zap.Error(err))
	}
}
