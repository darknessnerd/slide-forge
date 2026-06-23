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

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/your-org/your-repo/internal/config"
	"github.com/your-org/your-repo/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		level = slog.LevelInfo
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))

	gen := handler.NewGeneratorAdapter()
	mcpServer := handler.NewMCPServer(gen)

	switch cfg.Transport {
	case "stdio":
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		logger.Info("slide-forge MCP server starting", "transport", "stdio")
		if err := mcpServer.Run(ctx, mcp.NewStdioTransport()); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}

	case "http":
		mux := http.NewServeMux()
		mux.HandleFunc("GET /health", healthHandler)
		mux.HandleFunc("GET /ready", readyHandler)
		mux.Handle("/mcp/", mcp.NewStreamableHTTPHandler(
			func(r *http.Request) *mcp.Server { return mcpServer }, nil,
		))

		srv := &http.Server{Addr: cfg.Addr, Handler: mux}
		logger.Info("slide-forge MCP server starting", "transport", "http", "addr", cfg.Addr, "endpoint", "/mcp/")

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		serverErr := make(chan error, 1)
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				serverErr <- err
			}
		}()

		select {
		case err := <-serverErr:
			logger.Error("server error", "err", err)
			os.Exit(1)
		case sig := <-stop:
			logger.Info("shutting down", "signal", sig)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				logger.Error("shutdown error", "err", err)
				os.Exit(1)
			}
			logger.Info("shutdown complete")
		}

	default:
		logger.Error("unknown transport", "MCP_TRANSPORT", cfg.Transport, "valid", "stdio|http")
		os.Exit(1)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}
