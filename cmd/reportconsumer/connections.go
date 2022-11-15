package main

import (
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/config"
	"gitlab.ozon.dev/egor.linkinked/kartashov-egor/internal/grpcapi/client_interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func mustConnectToDb(cfg *config.Service) *sqlx.DB {
	return sqlx.MustConnect("postgres", cfg.Dsn())
}

func connectToGrpcServer(cfg *config.Service) (*grpc.ClientConn, error) {
	metricsInterceptor := client_interceptors.NewMetrics(ServiceName)
	return grpc.Dial(
		cfg.SendReportAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
			client_interceptors.UnaryAuth(cfg.SendReportClientSecret()),
			metricsInterceptor.Unary(),
		),
	)
}
