package commands

import (
	"fmt"

	"github.com/KatrinSalt/backend-challenge-go/cmd/cli/output"
	"github.com/urfave/cli/v2"
)

func ProcessInvoice() *cli.Command {
	return &cli.Command{
		Name:    "process-invoice",
		Aliases: []string{"invoice", "i"},
		Usage:   "Process an invoice through the approval workflow",
		UsageText: ` 
		    backend-challenge-cli process-invoice
		    backend-challenge-cli invoice
		    backend-challenge-cli i`,
		Action: func(c *cli.Context) error {
			// Get CLI config from global flags
			cliConfig := &Config{
				Company:     c.String("company"),
				Departments: c.String("departments"),
				SlackConn:   c.String("slack-connection-string"),
				EmailConn:   c.String("email-connection-string"),
				Verbose:     c.Bool("verbose"),
			}

			// Setup workflow services using CLI config
			services, err := setupServicesWithConfig(cliConfig)
			if err != nil {
				return fmt.Errorf("failed to setup services: %w", err)
			}

			// Show configuration if verbose mode is enabled
			if cliConfig.Verbose {
				output.Printf("🔧 Configuration:\n")
				output.Printf("   Company: %s\n", cliConfig.Company)
				output.Printf("   Departments: %s\n", cliConfig.Departments)
				if cliConfig.SlackConn != "" {
					output.Printf("   Slack Connection: %s\n", cliConfig.SlackConn[:8]+"...")
				}
				if cliConfig.EmailConn != "" {
					output.Printf("   Email Connection: %s\n", cliConfig.EmailConn[:8]+"...")
				}
				output.Println("")
			}

			// Run the interactive workflow
			err = services.Workflow.Run()
			if err != nil {
				return fmt.Errorf("failed to process invoice: %w", err)
			}

			output.Println("✅ Invoice processing completed successfully!")
			return nil
		},
	}
}
