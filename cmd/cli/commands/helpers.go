package commands

import (
	"fmt"

	"github.com/KatrinSalt/backend-challenge-go/common"
	"github.com/KatrinSalt/backend-challenge-go/config"
)

// Config holds the CLI configuration
type Config struct {
	Company     string
	Departments string
	SlackConn   string
	EmailConn   string
	Verbose     bool
}

// setupServicesWithConfig loads configuration using CLI config and sets up all required services
func setupServicesWithConfig(cliConfig *Config) (*config.Services, error) {
	// Build flags from CLI config
	flags := []string{
		"--company", cliConfig.Company,
	}

	// Add optional flags if they are set
	if cliConfig.Departments != "" {
		flags = append(flags, "--departments", cliConfig.Departments)
	}
	if cliConfig.SlackConn != "" {
		flags = append(flags, "--slack-connection-string", cliConfig.SlackConn)
	}
	if cliConfig.EmailConn != "" {
		flags = append(flags, "--email-connection-string", cliConfig.EmailConn)
	}

	// Parse flags
	parsedFlags, err := config.ParseFlags(flags)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// Create logger
	logger := common.NewLogger()

	// Load configuration
	cfg, err := config.New(config.WithFlags(parsedFlags))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup services
	services, err := config.SetUpServices(logger, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to setup services: %w", err)
	}

	return services, nil
}
