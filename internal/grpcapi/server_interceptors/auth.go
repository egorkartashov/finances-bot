package server_interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuth(clientSecrets []string) grpc.UnaryServerInterceptor {
	secretsMap := toMap(clientSecrets)
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			err = status.Error(codes.Unauthenticated, "metadata not found")
			return
		}
		clientSecret, ok := md["client_secret"]
		if !ok || len(clientSecret) < 1 {
			err = status.Error(codes.Unauthenticated, "client secret not found in metadata")
			return
		}
		if _, ok := secretsMap[clientSecret[0]]; !ok {
			err = status.Error(codes.Unauthenticated, "invalid client secret")
			return
		}
		return handler(ctx, req)
	}
}

func toMap(secrets []string) map[string]struct{} {
	m := make(map[string]struct{}, len(secrets))
	for _, s := range secrets {
		m[s] = struct{}{}
	}
	return m
}
