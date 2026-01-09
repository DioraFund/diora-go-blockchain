package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DioraFund/diora-go-blockchain/internal/cli"
	"github.com/DioraFund/diora-go-blockchain/internal/config"
	"github.com/DioraFund/diora-go-blockchain/internal/logger"
)

const (
	// Version information
	Version = "1.0.0"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	// Initialize context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Initialize logger
	log := logger.NewLogger("diora-cli")
	log.Info("Starting Diora CLI", "version", Version, "commit", Commit, "date", Date)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize CLI
	cliApp, err := cli.NewCLI(ctx, cfg, log)
	if err != nil {
		log.Error("Failed to initialize CLI", "error", err)
		os.Exit(1)
	}

	// Run CLI application
	if err := cliApp.Run(); err != nil {
		log.Error("CLI application failed", "error", err)
		os.Exit(1)
	}

	log.Info("Diora CLI shutdown complete")
}
