package config

import (
	"flag"
	"fmt"
	"os"
)

// flags holds the command line flags for the config package.
type flags struct {
	companyName string
	departments string
	slack       string
	email       string
}

// ParseFlags parses the command line flags and returns a flags struct.
func ParseFlags(args []string) (*flags, error) {
	var (
		f flags
	)

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Command-line configuration for the workflow service:\n\n")
		fs.PrintDefaults()
	}

	fs.StringVar(&f.companyName, "company", "", "A company name for which the workflow service will be configured.")
	fs.StringVar(&f.departments, "departments", "", "Comma-separated list of departments (e.g., 'Finance,Marketing').")
	fs.StringVar(&f.slack, "slack-connection-string", "", "A connection string for the slack service.")
	fs.StringVar(&f.email, "email-connection-string", "", "A connection string for the email service.")

	if err := fs.Parse(args); err != nil {
		return &f, err
	}

	return &f, nil
}
