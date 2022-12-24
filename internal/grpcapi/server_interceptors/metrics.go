package server_interceptors

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/grpcapi/utils"
	"google.golang.org/grpc"
)

type Metrics struct {
	grpcServerCounter *prometheus.CounterVec
}

func NewMetrics(subsystem string) *Metrics {
	return &Metrics{
		grpcServerCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "finances",
				Subsystem: subsystem,
				Name:      "grpc_method_server_counter",
			},
			[]string{"method", "code"},
		),
	}
}

func (m *Metrics) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)

		code := utils.GetResponseCode(err)
		m.grpcServerCounter.WithLabelValues(info.FullMethod, code).Inc()

		return resp, err
	}
}
