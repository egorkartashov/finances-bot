package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/grpcapi/send_report"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/grpcapi/server_interceptors"
	"gitlab.ozon.dev/egor.linkinked/finances-bot/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func startGrpcServer(
	ctx context.Context, port int, clientSecrets []string, sendReportServer send_report.ReportSenderServer,
) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("failed to start grpc server")
	}
	metricsInterceptor := server_interceptors.NewMetrics("bot")

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
			server_interceptors.UnaryAuth(clientSecrets),
			metricsInterceptor.Unary(),
		),
	)
	send_report.RegisterReportSenderServer(s, sendReportServer)

	go func() {
		<-ctx.Done()
		s.GracefulStop()
		if err := lis.Close(); err != nil {
			logger.Fatal("failed to stop listening http for grpc server", zap.Error(err))
		}
	}()

	logger.Info(fmt.Sprintf("gRPC server listening at %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
