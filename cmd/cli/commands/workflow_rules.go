package commands

import (
	"fmt"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/cmd/cli/output"
	"github.com/urfave/cli/v2"
)

func CreateWorkflowRule() *cli.Command {
	return &cli.Command{
		Name:    "create-workflow-rule",
		Aliases: []string{"cwr"},
		Usage:   "Create a new workflow rule",
		UsageText: ` 
		    backend-challenge-cli create-workflow-rule --min-amount 100 --max-amount 500 --department "Finance" --approver-id 1 --approval-channel 0 --manager-approval 1
		    backend-challenge-cli cwr -min 100 -max 500 -d "Marketing" -aid 1 -ac 0 -ma 0`,
		Flags: []cli.Flag{
			&cli.Float64Flag{
				Name:    "min-amount",
				Aliases: []string{"min"},
				Usage:   "Minimum amount for the rule (optional)",
			},
			&cli.Float64Flag{
				Name:    "max-amount",
				Aliases: []string{"max"},
				Usage:   "Maximum amount for the rule (optional)",
			},
			&cli.StringFlag{
				Name:    "department",
				Aliases: []string{"d"},
				Usage:   "Department for the rule (optional)",
			},
			&cli.IntFlag{
				Name:     "approver-id",
				Aliases:  []string{"aid"},
				Usage:    "ID of the approver for this rule, required",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "approval-channel",
				Aliases:  []string{"ac"},
				Usage:    "Approval channel (0=Slack, 1=Email), required",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "manager-approval",
				Aliases: []string{"ma"},
				Usage:   "Whether manager approval is required (0=No, 1=Yes, optional)",
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

			// Create workflow rule
			rule := api.WorkflowRule{
				CompanyID:       1, // Light company ID from sample data
				ApproverID:      c.Int("approver-id"),
				ApprovalChannel: c.Int("approval-channel"),
			}

			// Set optional fields
			if c.IsSet("min-amount") {
				minAmount := c.Float64("min-amount")
				rule.MinAmount = &minAmount
			}
			if c.IsSet("max-amount") {
				maxAmount := c.Float64("max-amount")
				rule.MaxAmount = &maxAmount
			}
			if c.IsSet("department") {
				department := c.String("department")
				rule.Department = &department
			}
			if c.IsSet("manager-approval") {
				rule.IsManagerApprovalRequired = c.Int("manager-approval")
			} else {
				// Set default value to 0 (No) if not specified
				rule.IsManagerApprovalRequired = 0
			}

			createdRule, err := services.Management.CreateWorkflowRule(rule)
			if err != nil {
				return fmt.Errorf("failed to create workflow rule: %w", err)
			}

			message := fmt.Sprintf("✅ Workflow rule created successfully!\n"+
				"ID: %d\n"+
				"Min Amount: %s\n"+
				"Max Amount: %s\n"+
				"Department: %s\n"+
				"Manager Approval Required: %s\n"+
				"Approver ID: %d\n"+
				"Approval Channel: %s",
				createdRule.ID,
				formatFloatPtr(createdRule.MinAmount),
				formatFloatPtr(createdRule.MaxAmount),
				formatStringPtr(createdRule.Department),
				formatManagerApproval(createdRule.IsManagerApprovalRequired),
				createdRule.ApproverID,
				formatApprovalChannel(createdRule.ApprovalChannel))
			output.Println(message)
			return nil
		},
	}
}

func UpdateWorkflowRule() *cli.Command {
	return &cli.Command{
		Name:    "update-workflow-rule",
		Aliases: []string{"uwr"},
		Usage:   "Update an existing workflow rule",
		UsageText: ` 
		    backend-challenge-cli update-workflow-rule --id 1 --min-amount 200 --max-amount 600 --department "Finance" --approver-id 2 --approval-channel 1 --manager-approval 0
		    backend-challenge-cli uwr -i 1 -min 200 -max 600 -d "Marketing" -aid 2 -ac 1 -ma 1`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the workflow rule to update, required",
				Required: true,
			},
			&cli.Float64Flag{
				Name:    "min-amount",
				Aliases: []string{"min"},
				Usage:   "Minimum amount for the rule (optional)",
			},
			&cli.Float64Flag{
				Name:    "max-amount",
				Aliases: []string{"max"},
				Usage:   "Maximum amount for the rule (optional)",
			},
			&cli.StringFlag{
				Name:    "department",
				Aliases: []string{"d"},
				Usage:   "Department for the rule (optional)",
			},
			&cli.IntFlag{
				Name:     "approver-id",
				Aliases:  []string{"aid"},
				Usage:    "ID of the approver for this rule, required",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "approval-channel",
				Aliases:  []string{"ac"},
				Usage:    "Approval channel (0=Slack, 1=Email), required",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "manager-approval",
				Aliases: []string{"ma"},
				Usage:   "Whether manager approval is required (0=No, 1=Yes, optional)",
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

			// Update workflow rule
			rule := api.WorkflowRule{
				ID:              c.Int("id"),
				ApproverID:      c.Int("approver-id"),
				ApprovalChannel: c.Int("approval-channel"),
			}

			// Set optional fields
			if c.IsSet("min-amount") {
				minAmount := c.Float64("min-amount")
				rule.MinAmount = &minAmount
			}
			if c.IsSet("max-amount") {
				maxAmount := c.Float64("max-amount")
				rule.MaxAmount = &maxAmount
			}
			if c.IsSet("department") {
				department := c.String("department")
				rule.Department = &department
			}
			if c.IsSet("manager-approval") {
				rule.IsManagerApprovalRequired = c.Int("manager-approval")
			} else {
				// Set default value to 0 (No) if not specified
				rule.IsManagerApprovalRequired = 0
			}

			err = services.Management.UpdateWorkflowRule(rule)
			if err != nil {
				return fmt.Errorf("failed to update workflow rule: %w", err)
			}

			message := fmt.Sprintf("✅ Workflow rule updated successfully!\n"+
				"ID: %d\n"+
				"Min Amount: %s\n"+
				"Max Amount: %s\n"+
				"Department: %s\n"+
				"Manager Approval Required: %s\n"+
				"Approver ID: %d\n"+
				"Approval Channel: %s",
				rule.ID,
				formatFloatPtr(rule.MinAmount),
				formatFloatPtr(rule.MaxAmount),
				formatStringPtr(rule.Department),
				formatManagerApproval(rule.IsManagerApprovalRequired),
				rule.ApproverID,
				formatApprovalChannel(rule.ApprovalChannel))
			output.Println(message)
			return nil
		},
	}
}

func DeleteWorkflowRule() *cli.Command {
	return &cli.Command{
		Name:    "delete-workflow-rule",
		Aliases: []string{"dwr"},
		Usage:   "Delete a workflow rule by ID",
		UsageText: ` 
		    backend-challenge-cli delete-workflow-rule --id 1
		    backend-challenge-cli dwr -i 1`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the workflow rule to delete, required",
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

			// Delete workflow rule
			id := c.Int("id")
			err = services.Management.DeleteWorkflowRule(id)
			if err != nil {
				return fmt.Errorf("failed to delete workflow rule: %w", err)
			}

			output.Println(fmt.Sprintf("✅ Workflow rule with ID %d deleted successfully!", id))
			return nil
		},
	}
}

func GetWorkflowRuleByID() *cli.Command {
	return &cli.Command{
		Name:    "get-workflow-rule",
		Aliases: []string{"gwr"},
		Usage:   "Get a workflow rule by ID",
		UsageText: ` 
		    backend-challenge-cli get-workflow-rule --id 1
		    backend-challenge-cli gwr -i 1`,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "ID of the workflow rule to fetch, required",
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

			// Get workflow rule
			id := c.Int("id")
			rule, err := services.Management.GetWorkflowRuleByID(id)
			if err != nil {
				return fmt.Errorf("failed to get workflow rule: %w", err)
			}

			message := fmt.Sprintf("✅ Workflow rule found!\n"+
				"ID: %d\n"+
				"Min Amount: %s\n"+
				"Max Amount: %s\n"+
				"Department: %s\n"+
				"Manager Approval Required: %s\n"+
				"Approver ID: %d\n"+
				"Approval Channel: %s",
				rule.ID,
				formatFloatPtr(rule.MinAmount),
				formatFloatPtr(rule.MaxAmount),
				formatStringPtr(rule.Department),
				formatManagerApproval(rule.IsManagerApprovalRequired),
				rule.ApproverID,
				formatApprovalChannel(rule.ApprovalChannel))
			output.Println(message)
			return nil
		},
	}
}

func ListWorkflowRules() *cli.Command {
	return &cli.Command{
		Name:    "list-workflow-rules",
		Aliases: []string{"lwr"},
		Usage:   "List all workflow rules for the company",
		UsageText: ` 
		    backend-challenge-cli list-workflow-rules
		    backend-challenge-cli lwr`,
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

			// List workflow rules
			rules, err := services.Management.ListWorkflowRules()
			if err != nil {
				return fmt.Errorf("failed to list workflow rules: %w", err)
			}

			if len(rules) == 0 {
				output.Println("No workflow rules found for this company.")
			} else {
				output.Println(fmt.Sprintf("Found %d workflow rule(s):", len(rules)))
				for _, rule := range rules {
					message := fmt.Sprintf("ID: %d | Min: %s | Max: %s | Dept: %s | Manager: %s | Approver: %d | Channel: %s",
						rule.ID,
						formatFloatPtr(rule.MinAmount),
						formatFloatPtr(rule.MaxAmount),
						formatStringPtr(rule.Department),
						formatManagerApproval(rule.IsManagerApprovalRequired),
						rule.ApproverID,
						formatApprovalChannel(rule.ApprovalChannel))
					output.Println(message)
				}
			}
			return nil
		},
	}
}

// Helper functions for formatting optional fields
func formatFloatPtr(f *float64) string {
	if f == nil {
		return "Any"
	}
	return fmt.Sprintf("%.2f", *f)
}

func formatStringPtr(s *string) string {
	if s == nil {
		return "Any"
	}
	return *s
}

func formatManagerApproval(approval int) string {
	switch approval {
	case 0:
		return "No"
	case 1:
		return "Yes"
	default:
		return "Any"
	}
}

func formatApprovalChannel(channel int) string {
	switch channel {
	case 0:
		return "Slack"
	case 1:
		return "Email"
	default:
		return fmt.Sprintf("Unknown (%d)", channel)
	}
}
