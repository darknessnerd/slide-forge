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

	"github.com/mark3labs/mcp-go/server"
	"github.com/your-org/your-repo/internal/handler"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	gen := handler.NewGeneratorAdapter()
	mcpServer := handler.NewMCPServer(gen)

	transport := os.Getenv("MCP_TRANSPORT")
	if transport == "" {
		transport = "stdio"
	}

	switch transport {
	case "stdio":
		logger.Info("slide-forge MCP server starting", "transport", "stdio")
		if err := server.ServeStdio(mcpServer); err != nil {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}

	case "http":
		addr := os.Getenv("ADDR")
		if addr == "" {
			addr = ":8080"
		}

		httpServer := server.NewStreamableHTTPServer(mcpServer,
			server.WithEndpointPath("/mcp"),
		)

		logger.Info("slide-forge MCP server starting", "transport", "http", "addr", addr, "endpoint", "/mcp")

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		serverErr := make(chan error, 1)
		go func() {
			if err := httpServer.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
			if err := httpServer.Shutdown(ctx); err != nil {
				logger.Error("shutdown error", "err", err)
				os.Exit(1)
			}
			logger.Info("shutdown complete")
		}

	default:
		logger.Error("unknown transport", "MCP_TRANSPORT", transport, "valid", "stdio|http")
		os.Exit(1)
	}
}
