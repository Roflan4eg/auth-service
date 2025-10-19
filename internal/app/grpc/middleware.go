package grpc

import (
	"github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/interceptors"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"google.golang.org/grpc"
)

func WithInterceptors(logger *logger.Logger) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.MetricsInterceptor(),
			interceptors.Validation(),
			interceptors.Logging(logger),
			interceptors.Auth(),
			interceptors.Recovery(logger),
		),
	}
}
