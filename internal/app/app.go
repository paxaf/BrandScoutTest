package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	storage "github.com/paxaf/BrandScoutTest/internal/repo/engine"
	"github.com/paxaf/BrandScoutTest/internal/usecase"
)

const (
	appHost                      = "0.0.0.0"
	appPort                      = "8080"
	defaultTimeout time.Duration = 5 * time.Second
)

type App struct {
	apiServer *http.Server
}

func New() (*App, error) {
	app := &App{}
	repo, err := storage.NewEngine()
	if err != nil {
		return nil, fmt.Errorf("failed init repo: %w", err)
	}
	_ = usecase.New(repo)
	addr := net.JoinHostPort(appHost, appPort)
	app.apiServer = &http.Server{
		Addr:              addr,
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: defaultTimeout,
	}
	return app, nil
}

func (app *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("API server started successfully. " + "Address: " + app.apiServer.Addr)
		if err := app.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to start server")
		}
	}()

	<-ctx.Done()
	log.Printf("Received shutdown signal")

	return nil
}

func (app *App) Close() error {
	err := app.apiServer.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}
