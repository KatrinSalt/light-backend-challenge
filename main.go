package main

import (
	"log"
	"os"

	"github.com/KatrinSalt/backend-challenge-go/common"
	"github.com/KatrinSalt/backend-challenge-go/config"
)

func main() {
	// Parse command line flags.
	flags, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	// Create logger.
	logger := common.NewLogger()
	logger.Info("Starting CLI service")

	// Load configuration.
	cfg, err := config.New(config.WithFlags(flags))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	logger.Info("Configuration loaded", "config", cfg)

	// Setup CLI service.
	cliService, err := config.SetUpCliService(logger, cfg)
	if err != nil {
		log.Fatalf("Failed to setup CLI service: %v", err)
	}

	// Run the CLI service.
	if err := cliService.Run(); err != nil {
		log.Fatalf("CLI service failed: %v", err)
	}
}
