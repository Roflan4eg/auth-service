package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Auth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		authenticatedCtx, err := authenticate(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "authentication failed")
		}

		return handler(authenticatedCtx, req)
	}
}

func isPublicMethod(method string) bool {
	publicMethods := map[string]bool{
		"/auth.AuthService/Login":    true,
		"/auth.AuthService/Register": true,
		"/auth.AuthService/Health":   true,
	}
	return publicMethods[method]
}

func authenticate(ctx context.Context) (context.Context, error) {
	// TODO
	return ctx, nil
}
