package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           Crowdfunding API
// @version         1.0
// @description     REST API for crowdfunding campaigns and donations.
// @license.name    MIT
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Type "Bearer" followed by a space and your JWT, e.g. "Bearer eyJhbGci..."
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app, err := InitializeApp()
	if err != nil {
		logger.Error("failed to start application", "error", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    ":" + app.Config.AppPort,
		Handler: app.Router,
	}

	go func() {
		logger.Info("server listening", "port", app.Config.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for an interrupt, then give in-flight requests a few seconds to drain.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
