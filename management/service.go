package management

import (
	"fmt"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/common"
	"github.com/KatrinSalt/backend-challenge-go/db"
)

// databaseService defines the interface for database operations needed by the management service.
type databaseService interface {
	// Company operations
	GetCompanyByName(name string) (db.Company, error)

	// Workflow Rule operations
	CreateWorkflowRule(rule db.WorkflowRule) (db.WorkflowRule, error)
	GetWorkflowRuleByID(id int) (db.WorkflowRule, error)
	UpdateWorkflowRule(rule db.WorkflowRule) error
	DeleteWorkflowRule(id int) error
	ListWorkflowRules(companyID int) ([]db.WorkflowRule, error)

	// Approver operations
	CreateApprover(approver db.Approver) (db.Approver, error)
	GetApproverByID(id int) (db.Approver, error)
	UpdateApprover(approver db.Approver) error
	DeleteApprover(id int) error
	ListApprovers(companyID int) ([]db.Approver, error)
}

// Service defines the interface for management operations.
type Service interface {
	// Workflow Rule Management
	CreateWorkflowRule(rule api.WorkflowRule) (api.WorkflowRule, error)
	GetWorkflowRuleByID(id int) (api.WorkflowRule, error)
	UpdateWorkflowRule(rule api.WorkflowRule) error
	DeleteWorkflowRule(id int) error
	ListWorkflowRules() ([]api.WorkflowRule, error)

	// Approver Management
	CreateApprover(approver api.Approver) (api.Approver, error)
	GetApproverByID(id int) (api.Approver, error)
	UpdateApprover(approver api.Approver) error
	DeleteApprover(id int) error
	ListApprovers() ([]api.Approver, error)
}

// service implements the management service.
type service struct {
	logger    common.Logger
	dbService databaseService
	company   company
}

// Company represents a company in the management service.
type company struct {
	id   int    `json:"id"`
	name string `json:"name"`
}

// NewService creates a new management service.
func NewService(logger common.Logger, dbService databaseService, companyName string) (Service, error) {
	if companyName == "" {
		return nil, fmt.Errorf("company name is required")
	}

	// Query database to get company by name
	dbCompany, err := dbService.GetCompanyByName(companyName)
	if err != nil {
		return nil, fmt.Errorf("failed to get company %s: %w", companyName, err)
	}

	company := company{
		id:   dbCompany.ID,
		name: dbCompany.Name,
	}

	return &service{
		logger:    logger,
		dbService: dbService,
		company:   company,
	}, nil
}

// CreateWorkflowRule creates a new workflow rule.
func (s *service) CreateWorkflowRule(rule api.WorkflowRule) (api.WorkflowRule, error) {
	if err := rule.Validate(); err != nil {
		return api.WorkflowRule{}, fmt.Errorf("invalid workflow rule: %w", err)
	}

	// Set the company ID from the service
	rule.CompanyID = s.company.id

	// Convert API struct to DB struct
	dbRule := s.apiToDBWorkflowRule(rule)

	// Create in database
	createdRule, err := s.dbService.CreateWorkflowRule(dbRule)
	if err != nil {
		return api.WorkflowRule{}, fmt.Errorf("failed to create workflow rule: %w", err)
	}

	// Convert back to API struct
	return s.dbToAPIWorkflowRule(createdRule), nil
}

// GetWorkflowRuleByID retrieves a workflow rule by its ID.
func (s *service) GetWorkflowRuleByID(id int) (api.WorkflowRule, error) {
	if id <= 0 {
		return api.WorkflowRule{}, fmt.Errorf("invalid workflow rule ID: %d", id)
	}

	dbRule, err := s.dbService.GetWorkflowRuleByID(id)
	if err != nil {
		return api.WorkflowRule{}, fmt.Errorf("failed to get workflow rule: %w", err)
	}

	return s.dbToAPIWorkflowRule(dbRule), nil
}

// UpdateWorkflowRule updates an existing workflow rule.
func (s *service) UpdateWorkflowRule(rule api.WorkflowRule) error {
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("invalid workflow rule: %w", err)
	}

	if rule.ID <= 0 {
		return fmt.Errorf("invalid workflow rule ID: %d", rule.ID)
	}

	// Convert API struct to DB struct
	dbRule := s.apiToDBWorkflowRule(rule)

	// Update in database
	if err := s.dbService.UpdateWorkflowRule(dbRule); err != nil {
		return fmt.Errorf("failed to update workflow rule: %w", err)
	}

	return nil
}

// DeleteWorkflowRule deletes a workflow rule by its ID.
func (s *service) DeleteWorkflowRule(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid workflow rule ID: %d", id)
	}

	if err := s.dbService.DeleteWorkflowRule(id); err != nil {
		return fmt.Errorf("failed to delete workflow rule: %w", err)
	}

	return nil
}

// ListWorkflowRules retrieves all workflow rules for the company.
func (s *service) ListWorkflowRules() ([]api.WorkflowRule, error) {
	dbRules, err := s.dbService.ListWorkflowRules(s.company.id)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflow rules: %w", err)
	}

	// Convert DB structs to API structs
	apiRules := make([]api.WorkflowRule, len(dbRules))
	for i, dbRule := range dbRules {
		apiRules[i] = s.dbToAPIWorkflowRule(dbRule)
	}

	return apiRules, nil
}

// CreateApprover creates a new approver.
func (s *service) CreateApprover(approver api.Approver) (api.Approver, error) {
	if err := approver.Validate(); err != nil {
		return api.Approver{}, fmt.Errorf("invalid approver: %w", err)
	}

	// Convert API struct to DB struct
	dbApprover := s.apiToDBApprover(approver)

	// Create in database
	createdApprover, err := s.dbService.CreateApprover(dbApprover)
	if err != nil {
		return api.Approver{}, fmt.Errorf("failed to create approver: %w", err)
	}

	// Convert back to API struct
	return s.dbToAPIApprover(createdApprover), nil
}

// GetApproverByID retrieves an approver by their ID.
func (s *service) GetApproverByID(id int) (api.Approver, error) {
	if id <= 0 {
		return api.Approver{}, fmt.Errorf("invalid approver ID: %d", id)
	}

	dbApprover, err := s.dbService.GetApproverByID(id)
	if err != nil {
		return api.Approver{}, fmt.Errorf("failed to get approver: %w", err)
	}

	return s.dbToAPIApprover(dbApprover), nil
}

// UpdateApprover updates an existing approver.
func (s *service) UpdateApprover(approver api.Approver) error {
	if err := approver.Validate(); err != nil {
		return fmt.Errorf("invalid approver: %w", err)
	}

	if approver.ID <= 0 {
		return fmt.Errorf("invalid approver ID: %d", approver.ID)
	}

	// Convert API struct to DB struct
	dbApprover := s.apiToDBApprover(approver)

	// Update in database
	if err := s.dbService.UpdateApprover(dbApprover); err != nil {
		return fmt.Errorf("failed to update approver: %w", err)
	}

	return nil
}

// DeleteApprover deletes an approver by their ID.
func (s *service) DeleteApprover(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid approver ID: %d", id)
	}

	if err := s.dbService.DeleteApprover(id); err != nil {
		return fmt.Errorf("failed to delete approver: %w", err)
	}

	return nil
}

// ListApprovers retrieves all approvers for the company.
func (s *service) ListApprovers() ([]api.Approver, error) {
	dbApprovers, err := s.dbService.ListApprovers(s.company.id)
	if err != nil {
		return nil, fmt.Errorf("failed to list approvers: %w", err)
	}

	// Convert DB structs to API structs
	apiApprovers := make([]api.Approver, len(dbApprovers))
	for i, dbApprover := range dbApprovers {
		apiApprovers[i] = s.dbToAPIApprover(dbApprover)
	}

	return apiApprovers, nil
}

// Helper functions for converting between API and DB structs

func (s *service) apiToDBWorkflowRule(rule api.WorkflowRule) db.WorkflowRule {
	dbRule := db.WorkflowRule{
		ID:              rule.ID,
		CompanyID:       rule.CompanyID,
		MinAmount:       rule.MinAmount,
		MaxAmount:       rule.MaxAmount,
		Department:      rule.Department,
		ApproverID:      rule.ApproverID,
		ApprovalChannel: rule.ApprovalChannel,
	}

	// Convert bool to *int for IsManagerApprovalRequired
	if rule.IsManagerApprovalRequired != nil {
		var value int
		if *rule.IsManagerApprovalRequired {
			value = 1
		}
		dbRule.IsManagerApprovalRequired = &value
	}

	return dbRule
}

func (s *service) dbToAPIWorkflowRule(rule db.WorkflowRule) api.WorkflowRule {
	apiRule := api.WorkflowRule{
		ID:              rule.ID,
		CompanyID:       rule.CompanyID,
		MinAmount:       rule.MinAmount,
		MaxAmount:       rule.MaxAmount,
		Department:      rule.Department,
		ApproverID:      rule.ApproverID,
		ApprovalChannel: rule.ApprovalChannel,
	}

	// Convert *int to *bool for IsManagerApprovalRequired
	if rule.IsManagerApprovalRequired != nil {
		value := *rule.IsManagerApprovalRequired == 1
		apiRule.IsManagerApprovalRequired = &value
	}

	return apiRule
}

func (s *service) apiToDBApprover(approver api.Approver) db.Approver {
	return db.Approver{
		ID:        approver.ID,
		CompanyID: s.company.id, // Use company ID from service
		Name:      approver.Name,
		Role:      approver.Role,
		Email:     approver.Email,
		SlackID:   approver.SlackID,
	}
}

func (s *service) dbToAPIApprover(approver db.Approver) api.Approver {
	return api.Approver{
		ID:      approver.ID,
		Name:    approver.Name,
		Role:    approver.Role,
		Email:   approver.Email,
		SlackID: approver.SlackID,
	}
}
