package db

import (
	"fmt"
	"log"

	"github.com/KatrinSalt/backend-challenge-go/db/sql"
)

// dbService provides a centralized interface for all database operations.
type dbService struct {
	client            sql.Client
	sampleData        *SampleData
	CompanyStore      CompanyStore
	ApproverStore     ApproverStore
	WorkflowRuleStore WorkflowRuleStore
}

// ServiceOptions contains configuration options for the database service.
type ServiceOptions struct {
	CompanyTable      string
	ApproverTable     string
	WorkflowRuleTable string
}

// ServiceOption is a function that sets options on the database service.
type ServiceOption func(o *ServiceOptions)

// WithCompanyTable sets the company table name.
func WithCompanyTable(table string) ServiceOption {
	return func(o *ServiceOptions) {
		o.CompanyTable = table
	}
}

// WithApproverTable sets the approver table name.
func WithApproverTable(table string) ServiceOption {
	return func(o *ServiceOptions) {
		o.ApproverTable = table
	}
}

// WithWorkflowRuleTable sets the workflow rule table name.
func WithWorkflowRuleTable(table string) ServiceOption {
	return func(o *ServiceOptions) {
		o.WorkflowRuleTable = table
	}
}

// NewDBService creates a new database service with all stores.
func NewDBService(client sql.Client, options ...ServiceOption) (*dbService, error) {
	if client == nil {
		return nil, fmt.Errorf("nil sql client")
	}

	opts := &ServiceOptions{
		CompanyTable:      defaultCompanyTable,
		ApproverTable:     defaultApproverTable,
		WorkflowRuleTable: defaultWorkflowRuleTable,
	}

	for _, option := range options {
		option(opts)
	}

	// Create company store.
	companyStore, err := NewCompanyStore(client, func(o *CompanyStoreOptions) {
		o.Table = opts.CompanyTable
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create company store: %w", err)
	}

	// Create approver store.
	approverStore, err := NewApproverStore(client, func(o *ApproverStoreOptions) {
		o.Table = opts.ApproverTable
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create approver store: %w", err)
	}

	// Create workflow rule store.
	workflowRuleStore, err := NewWorkflowRuleStore(client, func(o *WorkflowRuleStoreOptions) {
		o.Table = opts.WorkflowRuleTable
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow rule store: %w", err)
	}

	return &dbService{
		client:            client,
		CompanyStore:      companyStore,
		ApproverStore:     approverStore,
		WorkflowRuleStore: workflowRuleStore,
		sampleData:        NewSampleData(),
	}, nil
}

// GetSampleData returns the sample data for the service.
func (s *dbService) GetSampleData() *SampleData {
	return s.sampleData
}

// Initialize executes a list of SQL queries to initialize the database schema.
func (s *dbService) Initialize(queries []string) error {
	// TODO: check logging strategy later on.
	log.Println("Initializing database schema...")

	for i, query := range queries {
		if _, err := s.client.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %d: %w", i+1, err)
		}
	}

	// TODO: check logging strategy later on.
	log.Println("Database schema initialized successfully")

	return nil
}

// SeedSampleData populates the initialized database with sample data.
func (s *dbService) SeedSampleData() error {
	log.Println("Seeding database with sample data...")

	// Check the sample data is available.
	if s.sampleData == nil {
		return fmt.Errorf("no sample data provided")
	}

	// Add companies.
	for _, company := range s.sampleData.Companies {
		_, err := s.CompanyStore.Create(company)
		if err != nil {
			return err
		}
	}

	// Add approvers.
	for _, approver := range s.sampleData.Approvers {
		_, err := s.ApproverStore.Create(approver)
		if err != nil {
			return err
		}
	}

	// Add workflow rules.
	for _, rule := range s.sampleData.WorkflowRules {
		_, err := s.WorkflowRuleStore.Create(rule)
		if err != nil {
			return err
		}
	}

	log.Println("Database seeded successfully")
	return nil
}
