package interceptors

import (
	"context"
	"errors"
	"github.com/Roflan4eg/auth-serivce/internal/domain"
	l "github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/Roflan4eg/auth-serivce/internal/services"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func Logging(logger *l.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		span := trace.SpanFromContext(ctx)
		traceID := span.SpanContext().TraceID()
		var traceIDString string

		if traceID.IsValid() {
			traceIDString = traceID.String()
		} else {
			traceIDString = "unknown"
		}

		ctx = l.WithTraceID(ctx, traceIDString)
		ctx = l.WithMethod(ctx, info.FullMethod)
		start := time.Now()

		logger.DebugContext(ctx,
			"gRPC request started",
			logger.Any("request", req))

		resp, err := handler(ctx, req)
		duration := time.Since(start)
		if err != nil {
			errCtx := l.ErrorCtx(ctx, err)
			logger.ErrorContext(
				errCtx,
				"gRPC request failed",
				logger.Duration("duration", duration),
				logger.String("error", err.Error()),
				//logger.String("status_code", statusCodeStr),
			)
		} else {
			logger.InfoContext(ctx,
				"gRPC request completed",
				logger.Duration("duration", duration),
				//logger.String("status_code", statusCodeStr),
			)
		}

		clientErr := translateError(err)

		return resp, clientErr
	}
}

func translateError(err error) error {
	if err == nil {
		return nil
	}

	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	translatedErr := translateCustomError(err)
	if translatedErr != nil {
		return translatedErr
	}

	return status.Error(codes.Internal, err.Error())
}

func translateCustomError(err error) error {

	errorMapping := map[error]codes.Code{
		services.ErrInvalidRefreshToken: codes.Unauthenticated,
		services.ErrTokenExpired:        codes.Unauthenticated,
		services.ErrTokenMalformed:      codes.Unauthenticated,
		services.ErrTokenExpired:        codes.Unauthenticated,
		domain.ErrSessionExpired:        codes.Unauthenticated,
		domain.ErrSessionNotFound:       codes.NotFound,
		domain.ErrUserNotFound:          codes.NotFound,
		domain.ErrSessionAlreadyExists:  codes.AlreadyExists,
		domain.ErrUserAlreadyExists:     codes.AlreadyExists,
		domain.ErrPermissionDenied:      codes.PermissionDenied,
	}
	originalErr := l.OriginalError(err)
	for knownErr, code := range errorMapping {
		if errors.Is(originalErr, knownErr) {
			return status.Errorf(code, err.Error())
		}
	}

	return nil
}
