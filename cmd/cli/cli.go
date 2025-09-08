package cli

import (
	"github.com/KatrinSalt/backend-challenge-go/cmd/cli/commands"
	"github.com/KatrinSalt/backend-challenge-go/cmd/cli/output"
	"github.com/urfave/cli/v2"
)

const (
	name = "backend-challenge-cli"
)

var (
	companyName string
	departments string
	slackConn   string
	emailConn   string
	verbose     bool
)

func CLI(args []string) int {
	app := &cli.App{
		Name:  name,
		Usage: "A CLI to handle invoice processing and manage approvers and workflow rules.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "company",
				Aliases:     []string{"c"},
				Usage:       "Company name for the workflow service",
				Value:       "Light", // Default value
				Destination: &companyName,
			},
			&cli.StringFlag{
				Name:        "departments",
				Aliases:     []string{"d"},
				Usage:       "Comma-separated list of departments (e.g., 'Finance,Marketing')",
				Value:       "Marketing,Finance", // Default value
				Destination: &departments,
			},
			&cli.StringFlag{
				Name:        "slack-connection-string",
				Usage:       "Connection string for Slack notifications (e.g., 'xoxb-token')",
				Destination: &slackConn,
			},
			&cli.StringFlag{
				Name:        "email-connection-string",
				Usage:       "Connection string for email notifications (e.g., 'smtp://user:pass@host:port')",
				Destination: &emailConn,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "Enable verbose output",
				Destination: &verbose,
			},
		},
		Commands: []*cli.Command{
			commands.ProcessInvoice(),
			// Approver commands
			commands.CreateApprover(),
			commands.UpdateApprover(),
			commands.DeleteApprover(),
			commands.GetApproverByID(),
			commands.ListApprovers(),
			// Workflow Rule commands
			commands.CreateWorkflowRule(),
			commands.UpdateWorkflowRule(),
			commands.DeleteWorkflowRule(),
			commands.GetWorkflowRuleByID(),
			commands.ListWorkflowRules(),
		},
		CustomAppHelpTemplate: `NAME:
	{{.HelpName}} - {{.Usage}}

USAGE:
   	{{.HelpName}} [global options] command [command options]

GLOBAL OPTIONS:
   	{{range .VisibleFlags}}{{.}}
   	{{end}}

COMMANDS:
   	{{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   	{{end}}

{{if .Commands}}
COMMAND DETAILS:
{{range .Commands}}
	{{.Name}}, {{join .Aliases ", "}}

    OPTIONS:
	  {{range .VisibleFlags}}    {{.}}
	  {{end}}{{if .UsageText}}
		EXAMPLE:
	  {{.UsageText}}{{end}}

{{end}}{{end}}`,
	}

	if err := app.Run(args); err != nil {
		output.PrintlnErr(err)
		return 1
	}
	return 0
}
