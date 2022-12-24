package client_interceptors

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/grpcapi/utils"
	"google.golang.org/grpc"
)

type Metrics struct {
	grpcClientCounter *prometheus.CounterVec
}

func NewMetrics(subsystem string) *Metrics {
	return &Metrics{
		grpcClientCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "finances",
				Subsystem: subsystem,
				Name:      "grpc_method_client_counter",
			},
			[]string{"method", "code"},
		),
	}
}
func (m *Metrics) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		err := invoker(ctx, method, req, reply, cc, opts...)

		code := utils.GetResponseCode(err)
		m.grpcClientCounter.WithLabelValues(method, code).Inc()

		return err
	}
}
