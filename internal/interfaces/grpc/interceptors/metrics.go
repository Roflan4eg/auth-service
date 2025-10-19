package interceptors

import (
	"context"
	"github.com/Roflan4eg/auth-serivce/internal/lib/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		serviceName := "auth-service"
		metrics.GRPCRequestsInFlight.WithLabelValues(serviceName, info.FullMethod).Inc()
		defer metrics.GRPCRequestsInFlight.WithLabelValues(serviceName, info.FullMethod).Dec()

		resp, err := handler(ctx, req)

		if err != nil {
			code := status.Code(err).String()
			metrics.GRPCRequestErrors.WithLabelValues(info.FullMethod, code).Inc()
		}

		duration := time.Since(start).Seconds()
		code := status.Code(err).String()

		metrics.GRPCRequestCount.WithLabelValues(info.FullMethod, code).Inc()
		metrics.GRPCRequestDuration.WithLabelValues(info.FullMethod, code).Observe(duration)

		return resp, err
	}
}
