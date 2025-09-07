package db

import (
	"fmt"

	"github.com/KatrinSalt/backend-challenge-go/db/sql"
	"github.com/KatrinSalt/backend-challenge-go/db/sqlite"
)

// Service interface for the database service.
type Service interface {
	GetSampleData() *SampleData
	Initialize() error
	SeedSampleData() error
	GetCompanyByName(name string) (Company, error)
	// Workflow Rule Management
	CreateWorkflowRule(rule WorkflowRule) (WorkflowRule, error)
	GetWorkflowRuleByID(id int) (WorkflowRule, error)
	ListWorkflowRules(companyID int) ([]WorkflowRule, error)
	UpdateWorkflowRule(rule WorkflowRule) error
	DeleteWorkflowRule(id int) error
	FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (WorkflowRule, error)
	// Approver Management
	CreateApprover(approver Approver) (Approver, error)
	GetApproverByID(id int) (Approver, error)
	ListApprovers(companyID int) ([]Approver, error)
	UpdateApprover(approver Approver) error
	DeleteApprover(id int) error
}

// Service provides a centralized interface for all database operations.
type service struct {
	client            sql.Client
	schema            []string
	sampleData        *SampleData
	companyStore      CompanyStore
	approverStore     ApproverStore
	workflowRuleStore WorkflowRuleStore
}

// ServiceOptions contains configuration options for the database service.
type ServiceOptions struct {
	Schema            []string
	SampleData        *SampleData
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

// WithSampleData sets the sample data.
func WithSampleData(sampleData *SampleData) ServiceOption {
	return func(o *ServiceOptions) {
		o.SampleData = sampleData
	}
}

// WithSchema sets the schema.
func WithSchema(schema []string) ServiceOption {
	return func(o *ServiceOptions) {
		o.Schema = schema
	}
}

// NewService creates a new database service with all stores.
func NewService(client sql.Client, options ...ServiceOption) (*service, error) {
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

	if len(opts.Schema) == 0 {
		opts.Schema = sqlite.NewDBSchema()
	}

	if opts.SampleData == nil {
		opts.SampleData = NewSampleData()
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

	return &service{
		client:            client,
		schema:            opts.Schema,
		sampleData:        opts.SampleData,
		companyStore:      companyStore,
		approverStore:     approverStore,
		workflowRuleStore: workflowRuleStore,
	}, nil
}

// GetSampleData returns the sample data for the service.
func (s *service) GetSampleData() *SampleData {
	return s.sampleData
}

// Initialize executes a list of SQL queries to initialize the database schema.
func (s *service) Initialize() error {
	for i, query := range s.schema {
		if _, err := s.client.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %d: %w", i+1, err)
		}
	}

	return nil
}

// SeedSampleData populates the initialized database with sample data.
func (s *service) SeedSampleData() error {
	// Check the sample data is available.
	if s.sampleData == nil {
		return fmt.Errorf("no sample data provided")
	}

	// Add companies.
	for _, company := range s.sampleData.Companies {
		_, err := s.companyStore.Create(company)
		if err != nil {
			return err
		}
	}

	// Add approvers.
	for _, approver := range s.sampleData.Approvers {
		_, err := s.approverStore.Create(approver)
		if err != nil {
			return err
		}
	}

	// Add workflow rules.
	for _, rule := range s.sampleData.WorkflowRules {
		_, err := s.workflowRuleStore.Create(rule)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCompanyByName retrieves a company by its name.
func (s *service) GetCompanyByName(name string) (Company, error) {
	return s.companyStore.GetByName(name)
}

// CreateApprover creates a new approver.
func (s *service) CreateApprover(approver Approver) (Approver, error) {
	return s.approverStore.Create(approver)
}

// GetApproverByID retrieves an approver by their ID.
func (s *service) GetApproverByID(id int) (Approver, error) {
	return s.approverStore.GetByID(id)
}

// FindMatchingRule finds a workflow rule that matches the given criteria.
func (s *service) FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (WorkflowRule, error) {
	return s.workflowRuleStore.FindMatchingRule(companyID, amount, department, requiresManager)
}

// CreateWorkflowRule creates a new workflow rule.
func (s *service) CreateWorkflowRule(rule WorkflowRule) (WorkflowRule, error) {
	return s.workflowRuleStore.Create(rule)
}

// GetWorkflowRuleByID retrieves a workflow rule by its ID.
func (s *service) GetWorkflowRuleByID(id int) (WorkflowRule, error) {
	return s.workflowRuleStore.GetByID(id)
}

// UpdateWorkflowRule updates an existing workflow rule.
func (s *service) UpdateWorkflowRule(rule WorkflowRule) error {
	return s.workflowRuleStore.Update(rule)
}

// DeleteWorkflowRule deletes a workflow rule by its ID.
func (s *service) DeleteWorkflowRule(id int) error {
	return s.workflowRuleStore.Delete(id)
}

// UpdateApprover updates an existing approver.
func (s *service) UpdateApprover(approver Approver) error {
	return s.approverStore.Update(approver)
}

// DeleteApprover deletes an approver by their ID.
func (s *service) DeleteApprover(id int) error {
	return s.approverStore.Delete(id)
}

// ListWorkflowRules retrieves all workflow rules for a specific company.
func (s *service) ListWorkflowRules(companyID int) ([]WorkflowRule, error) {
	return s.workflowRuleStore.List(companyID)
}

// ListApprovers retrieves all approvers for a specific company.
func (s *service) ListApprovers(companyID int) ([]Approver, error) {
	return s.approverStore.List(companyID)
}
