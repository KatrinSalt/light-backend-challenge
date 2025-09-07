package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/KatrinSalt/backend-challenge-go/db/sql"
)

var (
	ErrWorkflowRuleNotFound      = errors.New("workflow rule not found")
	ErrWorkflowRuleAlreadyExists = errors.New("workflow rule already exists")
)

// WorkflowRuleStore defines the interface for workflow rule operations
type WorkflowRuleStore interface {
	Create(workflowRule WorkflowRule) (WorkflowRule, error)
	GetByCompanyID(companyID int) ([]WorkflowRule, error)
	FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (WorkflowRule, error)
}

// workflowRuleStore implements WorkflowRuleStore
type workflowRuleStore struct {
	client sql.Client
	table  string
}

// WorkflowRuleStoreOptions contains options for the workflow rule store.
type WorkflowRuleStoreOptions struct {
	Table string
}

// WorkflowRuleStoreOption is a function that sets options on the workflow rule store.
type WorkflowRuleStoreOption func(o *WorkflowRuleStoreOptions)

// NewWorkflowRuleStore creates a new workflow rule store
func NewWorkflowRuleStore(client sql.Client, options ...WorkflowRuleStoreOption) (*workflowRuleStore, error) {
	if client == nil {
		return nil, errors.New("nil sql client")
	}

	opts := WorkflowRuleStoreOptions{}
	for _, option := range options {
		option(&opts)
	}
	if len(opts.Table) == 0 {
		opts.Table = defaultWorkflowRuleTable
	}

	return &workflowRuleStore{
		client: client,
		table:  opts.Table,
	}, nil
}

// Create creates a new workflow rule.
func (s *workflowRuleStore) Create(workflowRule WorkflowRule) (WorkflowRule, error) {
	tx, err := s.client.Transaction()
	if err != nil {
		return WorkflowRule{}, err
	}
	defer tx.Rollback()

	insert := fmt.Sprintf("INSERT INTO %s (company_id, min_amount, max_amount, department, is_manager_approval_required, approver_id, approval_channel) VALUES ($1, $2, $3, $4, $5, $6, $7)", s.table)
	if _, err := tx.Exec(insert, workflowRule.CompanyID, workflowRule.MinAmount, workflowRule.MaxAmount, workflowRule.Department, workflowRule.IsManagerApprovalRequired, workflowRule.ApproverID, workflowRule.ApprovalChannel); err != nil {
		if strings.Contains(err.Error(), sql.SQLStateDuplicateKey) {
			return WorkflowRule{}, ErrWorkflowRuleAlreadyExists
		}
		return WorkflowRule{}, err
	}

	// Get the created workflow rule with its generated ID.
	var outWorkflowRule WorkflowRule
	query := fmt.Sprintf("SELECT id, company_id, min_amount, max_amount, department, is_manager_approval_required, approver_id, approval_channel FROM %s WHERE company_id = $1 AND approver_id = $2 ORDER BY id DESC LIMIT 1", s.table)
	if err := tx.QueryRow(query, workflowRule.CompanyID, workflowRule.ApproverID).Scan(&outWorkflowRule.ID, &outWorkflowRule.CompanyID, &outWorkflowRule.MinAmount, &outWorkflowRule.MaxAmount, &outWorkflowRule.Department, &outWorkflowRule.IsManagerApprovalRequired, &outWorkflowRule.ApproverID, &outWorkflowRule.ApprovalChannel); err != nil {
		return WorkflowRule{}, err
	}

	if err := tx.Commit(); err != nil {
		return WorkflowRule{}, err
	}

	return outWorkflowRule, nil
}

// GetByCompanyID retrieves all workflow rules for a specific company.
func (s *workflowRuleStore) GetByCompanyID(companyID int) ([]WorkflowRule, error) {
	query := fmt.Sprintf("SELECT id, company_id, min_amount, max_amount, department, is_manager_approval_required, approver_id, approval_channel FROM %s WHERE company_id = $1", s.table)

	rows, err := s.client.Query(query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workflow rules: %w", err)
	}
	defer rows.Close()

	var rules []WorkflowRule
	for rows.Next() {
		var rule WorkflowRule
		err := rows.Scan(
			&rule.ID,
			&rule.CompanyID,
			&rule.MinAmount,
			&rule.MaxAmount,
			&rule.Department,
			&rule.IsManagerApprovalRequired,
			&rule.ApproverID,
			&rule.ApprovalChannel,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow rule: %w", err)
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workflow rules: %w", err)
	}

	return rules, nil
}

func (s *workflowRuleStore) FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (WorkflowRule, error) {
	query := `
		SELECT id, company_id, min_amount, max_amount, department, 
		       is_manager_approval_required, approver_id, approval_channel 
		FROM workflow_rules 
		WHERE company_id = $1 
			AND (
				-- Amount logic: inclusive lower bound, exclusive upper bound
				(min_amount IS NULL OR $2 >= min_amount) AND
				(max_amount IS NULL OR $2 < max_amount)
			)
			AND (department IS NULL OR department = $3)
			AND (is_manager_approval_required IS NULL OR is_manager_approval_required = $4)
		ORDER BY 
			(CASE WHEN min_amount IS NOT NULL THEN 1 ELSE 0 END +
			 CASE WHEN max_amount IS NOT NULL THEN 1 ELSE 0 END +
			 CASE WHEN department IS NOT NULL THEN 1 ELSE 0 END +
			 CASE WHEN is_manager_approval_required IS NOT NULL THEN 1 ELSE 0 END) DESC,
			id
		LIMIT 1`

	var rule WorkflowRule

	// Convert bool to int: false -> 0, true -> 1
	managerApprovalInt := 0
	if requiresManager {
		managerApprovalInt = 1
	}

	err := s.client.QueryRow(query, companyID, amount, department, managerApprovalInt).Scan(
		&rule.ID, &rule.CompanyID, &rule.MinAmount, &rule.MaxAmount,
		&rule.Department, &rule.IsManagerApprovalRequired, &rule.ApproverID, &rule.ApprovalChannel)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WorkflowRule{}, ErrWorkflowRuleNotFound
		}
		return WorkflowRule{}, err
	}

	return rule, nil
}
