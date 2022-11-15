package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/logger"
	"go.uber.org/zap"
)

func StartMetricsHttpServer(ctx context.Context, port int) {
	srv := http.Server{
		Addr: fmt.Sprintf(":%v", port),
	}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("failed to shutdown server", zap.Error(err))
		}
	}()

	logger.Info(fmt.Sprintf("Starting HTTP metrics server on port %v", port))
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("failed to start http server", zap.Error(err))
	}
}
