package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ranas-mukminov/kkt-54fz-monitoring/internal/config"
	"github.com/ranas-mukminov/kkt-54fz-monitoring/pkg/logger"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	version := flag.Bool("version", false, "Print version information")
	flag.Parse()

	if *version {
		fmt.Printf("kkt-monitor version %s\n", Version)
		fmt.Printf("Build time: %s\n", BuildTime)
		fmt.Printf("Git commit: %s\n", GitCommit)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.New(cfg.Logging.Level, cfg.Logging.Format)
	log.Info("Starting KKT 54-FZ Monitoring System",
		"version", Version,
		"build_time", BuildTime,
		"commit", GitCommit,
	)

	// Create context for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// TODO: Initialize collectors
	// TODO: Initialize exporter
	// TODO: Initialize AI subsystem
	// TODO: Start HTTP server

	log.Info("KKT Monitor started successfully", "port", cfg.Server.Port)

	// Wait for shutdown signal
	<-sigChan
	log.Info("Shutdown signal received, stopping...")

	// Graceful shutdown
	cancel()
	log.Info("KKT Monitor stopped")
}
