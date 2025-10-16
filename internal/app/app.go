package app

import (
	"context"
	"fmt"
	"github.com/Roflan4eg/auth-serivce/config"
	"github.com/Roflan4eg/auth-serivce/internal/app/grpc"
	"github.com/Roflan4eg/auth-serivce/internal/interfaces/grpc/handlers"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"github.com/Roflan4eg/auth-serivce/internal/repository"
	"github.com/Roflan4eg/auth-serivce/internal/services"
	"github.com/Roflan4eg/auth-serivce/internal/storage"
	"sync"
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
	Name() string
}

type App struct {
	cfg        *config.Config
	logger     *logger.Logger
	storage    *storage.Container
	repository *repository.Container
	services   *services.Container
	handlers   *handlers.Container
	servers    []Server
	closer     *Closer
}

func New(cfg *config.Config, logger *logger.Logger) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
		closer: &Closer{},
	}
}

func (a *App) Setup() error {
	a.closer.Add(a.stopServers)
	stor, err := storage.NewContainer(a.cfg, a.logger)
	if err != nil {
		return fmt.Errorf("storage setup: %w", err)
	}
	a.storage = stor
	a.closer.Add(a.storage.Close)

	a.repository = repository.NewContainer(a.storage, a.cfg, a.logger)
	a.services = services.NewContainer(a.repository, a.cfg, a.logger)
	a.handlers = handlers.NewContainer(a.services, a.cfg, a.logger)

	if err = a.setupServers(); err != nil {
		return fmt.Errorf("server setup: %w", err)
	}

	return nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, server := range a.servers {
		go func(server Server) {
			if err := server.Start(); err != nil {
				a.logger.Error("Failed to start server",
					a.logger.String("server", server.Name()),
					a.logger.Any("error", err),
				)
			}
		}(server)
	}
	go StartShutdownListener(ctx, cancel, a.logger)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
	defer cancel()

	if err := a.closer.Close(shutdownCtx); err != nil {
		return err
	}

	return nil
}

func (a *App) setupServers() error {
	grpcServer := grpc.NewServer(
		a.handlers,
		a.logger,
		a.cfg.GRPC.Port,
	)
	a.servers = append(a.servers, grpcServer)

	return nil
}

func (a *App) stopServers(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(len(a.servers))
	for _, server := range a.servers {
		go func(server Server) {
			a.logger.Info("Stopping server", a.logger.String("server", server.Name()))
			if err := server.Stop(ctx); err != nil {
				a.logger.Warn("Failed to stop server",
					a.logger.String("server", server.Name()),
					a.logger.Any("error", err),
				)
			}
			wg.Done()
		}(server)
	}
	wg.Wait()
	return nil
}
