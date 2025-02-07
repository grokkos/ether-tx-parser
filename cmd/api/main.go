package main

import (
	"context"
	"fmt"
	"github.com/grokkos/ether-tx-parser/internal/api/http/handler"
	"github.com/grokkos/ether-tx-parser/internal/api/http/server"
	"github.com/grokkos/ether-tx-parser/internal/application/parser"
	"github.com/grokkos/ether-tx-parser/internal/infastructure/ethereum"
	"github.com/grokkos/ether-tx-parser/internal/infastructure/storage"
	"github.com/grokkos/ether-tx-parser/pkg/config"
	"github.com/grokkos/ether-tx-parser/pkg/logger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logger.GetLogger()
	defer logger.Sync()

	// Initialize dependencies
	client := ethereum.NewClient(cfg.Ethereum.RPCURL)
	store := storage.NewMemoryStore()
	service := parser.NewService(store, client)

	// Setup HTTP server
	parserHandler := handler.NewParserHandler(service)
	srv := server.NewServer(parserHandler)
	srv.SetupRoutes()

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logger.Info("Shutting down gracefully...")
		cancel()
	}()

	// Start parsing blocks in a goroutine
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info("Stopping block parser")
				return
			case <-ticker.C:
				if err := service.ParseBlocks(); err != nil {
					logger.Error("Error parsing blocks", zap.Error(err))
				}
			}
		}
	}()

	// Start HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Starting server", zap.String("address", addr))

	server := &http.Server{
		Addr:    addr,
		Handler: srv,
	}

	// Server shutdown on context cancellation
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("Server error", zap.Error(err))
	}
}
