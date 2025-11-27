package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"go-api-template/internal"
	"go-api-template/internal/logger"
	"go-api-template/internal/version"
)

func main() {
	// Get version information
	version, err := version.Get()
	if err != nil {
		panic(err)
	}

	// Initialize structured logging
	logger.Init(version)

	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		slog.Error("application error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv := internal.NewServer()

	// Use a configurable port or default to 8080
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("", port),
		Handler: srv,
	}

	go func() {
		slog.Info("starting server",
			slog.String("addr", httpServer.Addr),
			slog.String("port", port),
		)

		slog.Info(fmt.Sprintf("listening on address: %s", httpServer.Addr),
			slog.String("addr", httpServer.Addr),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error",
				slog.String("error", err.Error()),
			)
			slog.Error(fmt.Sprintf("error listening and serving: %s", err.Error()),
				slog.String("error", err.Error()),
			)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		slog.Info("shutting down server gracefully")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("error during shutdown",
				slog.String("error", err.Error()),
			)
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		} else {
			slog.Info("server shutdown complete")
		}
	}()
	wg.Wait()
	return nil
}
