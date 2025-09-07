package config

import (
	"context"
	"strings"

	"github.com/KatrinSalt/backend-challenge-go/db"
	"github.com/sethvargo/go-envconfig"
)

const (
	// Default company name.
	defaultCompanyName = "Light"
)

var (
	// Default company departments.
	defaultCompanyDepartments = []string{"Marketing", "Finance"}
)

const (
	defaultSlackConnectionString = "slack"
	defaultEmailConnectionString = "email"
)

// Configuration contains the configuration for the application.
type Configuration struct {
	Services Services
}

type Services struct {
	Workflow Workflow
	Database Database
	Slack    Slack
	Email    Email
}

type Workflow struct {
	Company Company
}

type Company struct {
	Name        string   `env:"COMPANY_NAME"`
	Departments []string `env:"COMPANY_DEPARTMENTS"`
}

// Placeholder for database configuration.
type Database struct {
	Schema     []string
	SampleData *db.SampleData
}

type Slack struct {
	ConnectionString string `env:"SLACK_CONNECTION_STRING"`
}

type Email struct {
	ConnectionString string `env:"EMAIL_CONNECTION_STRING"`
}

// Options contains options for creating new configurations.
type Options struct {
	Flags *flags
}

// Option is a function that sets options for new configurations.
type Option func(o *Options)

// New creates a new configuration.
func New(options ...Option) (Configuration, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	cfg := Configuration{
		Services: Services{
			Workflow: Workflow{
				Company: Company{
					Name:        defaultCompanyName,
					Departments: defaultCompanyDepartments,
				},
			},
			Database: Database{
				// Placeholder for database schema.
				Schema:     nil,
				SampleData: nil,
			},
			Slack: Slack{
				ConnectionString: defaultSlackConnectionString,
			},
			Email: Email{
				ConnectionString: defaultEmailConnectionString,
			},
		},
	}

	// Apply flag values if provided.
	if opts.Flags != nil {
		if opts.Flags.companyName != "" {
			cfg.Services.Workflow.Company.Name = opts.Flags.companyName
		}
		if opts.Flags.departments != "" {
			departments := strings.Split(opts.Flags.departments, ",")
			for i, dept := range departments {
				departments[i] = strings.TrimSpace(dept)
			}
			cfg.Services.Workflow.Company.Departments = departments
		}
		if opts.Flags.slack != "" {
			cfg.Services.Slack.ConnectionString = opts.Flags.slack
		}
		if opts.Flags.email != "" {
			cfg.Services.Email.ConnectionString = opts.Flags.email
		}
	}

	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// WithFlags sets the flags for the configuration.
func WithFlags(f *flags) Option {
	return func(o *Options) {
		o.Flags = f
	}
}
