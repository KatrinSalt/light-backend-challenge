package commands

import (
	"fmt"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/cmd/cli/output"
	"github.com/urfave/cli/v2"
)

func CreateApprover() *cli.Command {
	return &cli.Command{
		Name:    "create-approver",
		Aliases: []string{"ca"},
		Usage:   "Create a new approver",
		UsageText: ` 
		    backend-challenge-cli create-approver --name "John Doe" --role "Manager" --email "john@example.com" --slack-id "U123456"
		    backend-challenge-cli ca -n "Jane Smith" -r "Director" -e "jane@example.com" -s "U789012"`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the approver, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "role",
				Aliases:  []string{"r"},
				Usage:    "Role of the approver, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "email",
				Aliases:  []string{"e"},
				Usage:    "Email address of the approver, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "slack-id",
				Aliases:  []string{"s"},
				Usage:    "Slack ID of the approver, required",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			// Get CLI config from global flags
			cliConfig := &Config{
				Company:     c.String("company"),
				Departments: c.String("departments"),
				SlackConn:   c.String("slack-connection-string"),
				EmailConn:   c.String("email-connection-string"),
				Verbose:     c.Bool("verbose"),
			}

			// Setup services
			services, err := setupServicesWithConfig(cliConfig)
			if err != nil {
				return fmt.Errorf("failed to setup services: %w", err)
			}

			// Create approver
			approver := api.Approver{
				Name:    c.String("name"),
				Role:    c.String("role"),
				Email:   c.String("email"),
				SlackID: c.String("slack-id"),
			}

			createdApprover, err := services.Management.CreateApprover(approver)
			if err != nil {
				return fmt.Errorf("failed to create approver: %w", err)
			}

			message := fmt.Sprintf("✅ Approver created successfully!\n"+
				"ID: %d\n"+
				"Name: %s\n"+
				"Role: %s\n"+
				"Email: %s\n"+
				"Slack ID: %s",
				createdApprover.ID, createdApprover.Name, createdApprover.Role,
				createdApprover.Email, createdApprover.SlackID)
			output.Println(message)
			return nil
		},
	}
}

func UpdateApprover() *cli.Command {
	return &cli.Command{
		Name:    "update-approver",
		Aliases: []string{"ua"},
		Usage:   "Update an existing approver",
		UsageText: ` 
		    backend-challenge-cli update-approver --id 1 --name "John Doe Updated" --role "Senior Manager" --email "john.updated@example.com" --slack-id "U123456"
		    backend-challenge-cli ua -i 1 -n "Jane Smith Updated" -r "VP" -e "jane.updated@example.com" -s "U789012"`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the approver to update, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the approver, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "role",
				Aliases:  []string{"r"},
				Usage:    "Role of the approver, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "email",
				Aliases:  []string{"e"},
				Usage:    "Email address of the approver, required",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "slack-id",
				Aliases:  []string{"s"},
				Usage:    "Slack ID of the approver, required",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			// Get CLI config from global flags
			cliConfig := &Config{
				Company:     c.String("company"),
				Departments: c.String("departments"),
				SlackConn:   c.String("slack-connection-string"),
				EmailConn:   c.String("email-connection-string"),
				Verbose:     c.Bool("verbose"),
			}

			// Setup services
			services, err := setupServicesWithConfig(cliConfig)
			if err != nil {
				return fmt.Errorf("failed to setup services: %w", err)
			}

			// Update approver
			approver := api.Approver{
				ID:      c.Int("id"),
				Name:    c.String("name"),
				Role:    c.String("role"),
				Email:   c.String("email"),
				SlackID: c.String("slack-id"),
			}

			err = services.Management.UpdateApprover(approver)
			if err != nil {
				return fmt.Errorf("failed to update approver: %w", err)
			}

			message := fmt.Sprintf("✅ Approver updated successfully!\n"+
				"ID: %d\n"+
				"Name: %s\n"+
				"Role: %s\n"+
				"Email: %s\n"+
				"Slack ID: %s",
				approver.ID, approver.Name, approver.Role,
				approver.Email, approver.SlackID)
			output.Println(message)
			return nil
		},
	}
}

func DeleteApprover() *cli.Command {
	return &cli.Command{
		Name:    "delete-approver",
		Aliases: []string{"da"},
		Usage:   "Delete an approver by ID",
		UsageText: ` 
		    backend-challenge-cli delete-approver --id 1
		    backend-challenge-cli da -i 1`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the approver to delete, required",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			// Get CLI config from global flags
			cliConfig := &Config{
				Company:     c.String("company"),
				Departments: c.String("departments"),
				SlackConn:   c.String("slack-connection-string"),
				EmailConn:   c.String("email-connection-string"),
				Verbose:     c.Bool("verbose"),
			}

			// Setup services
			services, err := setupServicesWithConfig(cliConfig)
			if err != nil {
				return fmt.Errorf("failed to setup services: %w", err)
			}

			// Delete approver
			id := c.Int("id")
			err = services.Management.DeleteApprover(id)
			if err != nil {
				return fmt.Errorf("failed to delete approver: %w", err)
			}

			output.Println(fmt.Sprintf("✅ Approver with ID %d deleted successfully!", id))
			return nil
		},
	}
}

func GetApproverByID() *cli.Command {
	return &cli.Command{
		Name:    "get-approver",
		Aliases: []string{"ga"},
		Usage:   "Get an approver by ID",
		UsageText: ` 
		    backend-challenge-cli get-approver --id 1
		    backend-challenge-cli ga -i 1`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the approver to fetch, required",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			// Get CLI config from global flags
			cliConfig := &Config{
				Company:     c.String("company"),
				Departments: c.String("departments"),
				SlackConn:   c.String("slack-connection-string"),
				EmailConn:   c.String("email-connection-string"),
				Verbose:     c.Bool("verbose"),
			}

			// Setup services
			services, err := setupServicesWithConfig(cliConfig)
			if err != nil {
				return fmt.Errorf("failed to setup services: %w", err)
			}

			// Get approver
			id := c.Int("id")
			approver, err := services.Management.GetApproverByID(id)
			if err != nil {
				return fmt.Errorf("failed to get approver: %w", err)
			}

			message := fmt.Sprintf("✅ Approver found!\n"+
				"ID: %d\n"+
				"Name: %s\n"+
				"Role: %s\n"+
				"Email: %s\n"+
				"Slack ID: %s",
				approver.ID, approver.Name, approver.Role,
				approver.Email, approver.SlackID)
			output.Println(message)
			return nil
		},
	}
}

func ListApprovers() *cli.Command {
	return &cli.Command{
		Name:    "list-approvers",
		Aliases: []string{"la"},
		Usage:   "List all approvers for the company",
		UsageText: ` 
		    backend-challenge-cli list-approvers
		    backend-challenge-cli la`,
		Action: func(c *cli.Context) error {
			// Get CLI config from global flags
			cliConfig := &Config{
				Company:     c.String("company"),
				Departments: c.String("departments"),
				SlackConn:   c.String("slack-connection-string"),
				EmailConn:   c.String("email-connection-string"),
				Verbose:     c.Bool("verbose"),
			}

			// Setup services
			services, err := setupServicesWithConfig(cliConfig)
			if err != nil {
				return fmt.Errorf("failed to setup services: %w", err)
			}

			// List approvers
			approvers, err := services.Management.ListApprovers()
			if err != nil {
				return fmt.Errorf("failed to list approvers: %w", err)
			}

			if len(approvers) == 0 {
				output.Println("No approvers found for this company.")
			} else {
				output.Println(fmt.Sprintf("Found %d approver(s):", len(approvers)))
				for _, approver := range approvers {
					message := fmt.Sprintf("ID: %d | Name: %s | Role: %s | Email: %s | Slack ID: %s",
						approver.ID, approver.Name, approver.Role, approver.Email, approver.SlackID)
					output.Println(message)
				}
			}
			return nil
		},
	}
}
