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
