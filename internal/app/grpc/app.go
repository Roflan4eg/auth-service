package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/handlers"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	gRPCServer *grpc.Server
	port       string
	log        *logger.Logger
	name       string
}

func NewServer(
	handlers *handlers.Container,
	logger *logger.Logger,
	port string,
) *Server {
	opts := WithInterceptors(logger)
	grpcServer := grpc.NewServer(opts...)

	handlers.UserService.RegisterHandler(grpcServer)
	handlers.AuthService.RegisterHandler(grpcServer)

	return &Server{
		gRPCServer: grpcServer,
		port:       port,
		log:        logger,
	}
}

func (s *Server) Start() error {
	const op = "grpcserver.Run"
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("gRPC server starting",
		s.log.String("addr", l.Addr().String()),
		s.log.String("port", s.port),
	)

	if err = s.gRPCServer.Serve(l); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.gRPCServer.GracefulStop()
	return nil
}

func (s *Server) Name() string {
	return s.port
}
