package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/gtkit/logger"
)

// Options defines runtime startup options.
type Options struct {
	EnvFile string
}

// App owns the process lifecycle and assembled runtime dependencies.
type App struct {
	options    Options
	runtime    *Runtime
	httpConfig HTTPConfig
	httpServer httpServer
	workers    *WorkerManager
}

func NewApp(options Options) (*App, error) {
	if options.EnvFile == "" {
		return nil, errors.New("env file is required")
	}

	runtime, err := NewRuntime(options)
	if err != nil {
		return nil, err
	}

	providers, err := newProviders(runtime)
	if err != nil {
		runtime.Close()
		return nil, err
	}

	return &App{
		options:    options,
		runtime:    runtime,
		httpConfig: providers.httpConfig,
		httpServer: providers.httpServer,
		workers:    providers.workers,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.closeResources()

	runCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	workerCtx, cancelWorkers := context.WithCancel(runCtx)
	defer cancelWorkers()

	if err := a.workers.Start(workerCtx); err != nil {
		return fmt.Errorf("start workers: %w", err)
	}

	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- serveHTTPServer(a.httpServer, a.httpConfig)
	}()

	workerErrCh := a.workers.Errors()
	serverExited := false
	var resultErr error

	select {
	case err := <-serverErrCh:
		serverExited = true
		if err != nil {
			resultErr = err
		}
	case err := <-workerErrCh:
		if err != nil {
			resultErr = fmt.Errorf("worker failed: %w", err)
		}
	case <-runCtx.Done():
		logger.Info("app received shutdown signal")
	}

	cancelWorkers()

	if !serverExited {
		if err := shutdownHTTPServer(context.Background(), a.httpServer, a.httpConfig); err != nil && resultErr == nil {
			resultErr = fmt.Errorf("shutdown http server: %w", err)
		}

		if err := <-serverErrCh; err != nil && resultErr == nil {
			resultErr = err
		}
	}

	waitCtx, cancelWait := context.WithTimeout(context.Background(), a.httpConfig.ShutdownTimeout)
	defer cancelWait()
	if err := a.workers.Wait(waitCtx); err != nil && resultErr == nil {
		resultErr = fmt.Errorf("wait workers: %w", err)
	}

	return resultErr
}

func (a *App) closeResources() {
	if a.runtime != nil {
		a.runtime.Close()
	}
}
