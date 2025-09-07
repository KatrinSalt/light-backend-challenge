package config

import (
	"fmt"

	"github.com/KatrinSalt/backend-challenge-go/cli"
	"github.com/KatrinSalt/backend-challenge-go/common"
	"github.com/KatrinSalt/backend-challenge-go/db"
	"github.com/KatrinSalt/backend-challenge-go/db/sqlite"
	"github.com/KatrinSalt/backend-challenge-go/management"
	"github.com/KatrinSalt/backend-challenge-go/notification/email"
	"github.com/KatrinSalt/backend-challenge-go/notification/slack"
	"github.com/KatrinSalt/backend-challenge-go/workflow"
)

// Services contains the services for the application.
type services struct {
	CLI        cli.Service
	Management management.Service
}

func SetUpServices(log common.Logger, cfg Configuration) (*services, error) {
	dbSvc, err := setUpDatabaseService(cfg.Services.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to create database service: %v", err)
	}
	// Create CLI service.
	cliSvc, err := setUpCLIService(log, dbSvc, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create CLI service: %v", err)
	}
	// Create management service.
	managementSvc, err := setUpManagementService(log, dbSvc, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create management service: %v", err)
	}

	return &services{
		CLI:        cliSvc,
		Management: managementSvc,
	}, nil

}

// setUpCLIService creates and configures a CLI service.
func setUpCLIService(log common.Logger, dbSvc db.Service, cfg Configuration) (cli.Service, error) {
	// Setup workflow service.
	workflowSvc, err := setUpWorkflowService(log, dbSvc, cfg)
	if err != nil {
		return nil, err
	}
	// Create CLI service.
	return cli.NewService(workflowSvc, cli.WithLogger(log))
}

// setUpManagementService creates and configures a management service.
func setUpManagementService(log common.Logger, dbSvc db.Service, cfg Configuration) (management.Service, error) {
	// Create management service with company name.
	return management.NewService(log, dbSvc, cfg.Services.Company.Name)
}

func setUpWorkflowService(log common.Logger, dbSvc db.Service, cfg Configuration) (workflow.Service, error) {
	// Create slack notification service.
	slackSvc, err := slack.NewService(cfg.Services.Slack.ConnectionString, slack.WithLogger(log))
	if err != nil {
		return nil, fmt.Errorf("failed to create slack notification service: %v", err)
	}

	// Create email notification service.
	emailSvc, err := email.NewService(cfg.Services.Email.ConnectionString, email.WithLogger(log))
	if err != nil {
		return nil, fmt.Errorf("failed to create email notification service: %v", err)
	}

	// Create workflow service.
	workflowCompany := workflow.Company{
		Name:        cfg.Services.Company.Name,
		Departments: cfg.Services.Company.Departments,
	}
	workflowSvc, err := workflow.NewService(workflowCompany, dbSvc, slackSvc, emailSvc, workflow.WithLogger(log))
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow service: %v", err)
	}

	if err := workflowSvc.ValidateCompany(); err != nil {
		return nil, fmt.Errorf("failed to start workflow service for company %s: %v", cfg.Services.Company.Name, err)
	}

	return workflowSvc, nil
}

func setUpDatabaseService(cfg Database) (db.Service, error) {
	// Create SQL client.
	client, err := sqlite.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create sql client: %v", err)
	}

	// Set defaults if schema is nil or empty
	if len(cfg.Schema) == 0 {
		cfg.Schema = sqlite.NewDBSchema()
	}

	// Set defaults if sample data is nil or has no data
	if cfg.SampleData == nil ||
		(cfg.SampleData.Companies == nil && cfg.SampleData.Approvers == nil && cfg.SampleData.WorkflowRules == nil) {
		cfg.SampleData = db.NewSampleData()
	}

	svc, err := db.NewService(client, db.WithSchema(cfg.Schema), db.WithSampleData(cfg.SampleData))
	if err != nil {
		return nil, fmt.Errorf("failed to create database service: %v", err)
	}

	// Initialize database schema
	if err := svc.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %v", err)
	}

	// Seed with sample data
	if err := svc.SeedSampleData(); err != nil {
		return nil, fmt.Errorf("failed to seed sample data: %v", err)
	}

	return svc, nil
}
