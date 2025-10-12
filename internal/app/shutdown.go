package app

import (
	"context"
	"fmt"
	"github.com/Roflan4eg/auth-serivce/internal/lib/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Hook func(ctx context.Context) error

type Closer struct {
	mu    sync.Mutex
	hooks []Hook
}

func (c *Closer) Add(hook Hook) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hooks = append(c.hooks, hook)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	errorCh := make(chan string, len(c.hooks))
	readyCh := make(chan struct{}, 1)

	wg := sync.WaitGroup{}
	wg.Add(len(c.hooks))
	go func() {
		for _, hook := range c.hooks {
			go func(hook Hook) {
				defer wg.Done()
				if err := hook(ctx); err != nil {
					errorCh <- fmt.Sprintf("[!] %v\n", err)
				}
			}(hook)
		}
		wg.Wait()
		readyCh <- struct{}{}
		close(errorCh)
	}()

	select {
	case <-readyCh:
		break
	case sig := <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", sig)
	}
	if len(errorCh) > 0 {
		errString := fmt.Sprintf("Some shutdown hooks failed\n")
		for err := range errorCh {
			errString += err
		}
		return fmt.Errorf(errString)
	}
	return nil
}

func StartShutdownListener(ctx context.Context, cancel context.CancelFunc, log *logger.Logger) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	select {
	case sig := <-sigChan:
		log.Info("Received signal, initiating shutdown", log.String("signal", sig.String()))
		cancel()
		return

	case <-ctx.Done():
		log.Debug("Context cancelled, shutting down")
		return
	}
}
