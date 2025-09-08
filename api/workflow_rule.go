package api

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidApprovalChannel = errors.New("invalid approval channel")
	ErrInvalidAmountRange     = errors.New("invalid amount range")
)

// WorkflowRule represents a rule that determines how invoices are approved.
type WorkflowRule struct {
	ID                        int      `json:"id,omitempty"`
	CompanyID                 int      `json:"company_id,omitempty"`
	MinAmount                 *float64 `json:"min_amount,omitempty"`
	MaxAmount                 *float64 `json:"max_amount,omitempty"`
	Department                *string  `json:"department,omitempty"`
	IsManagerApprovalRequired int      `json:"is_manager_approval_required,omitempty"`
	ApproverID                int      `json:"approver_id"`
	ApprovalChannel           int      `json:"approval_channel"`
}

// Validate validates the workflow rule.
func (w *WorkflowRule) Validate() error {
	// Validate approval channel (0 = slack, 1 = email)
	if w.ApprovalChannel < 0 || w.ApprovalChannel > 1 {
		return ErrInvalidApprovalChannel
	}

	// Validate manager approval required (0 = no, 1 = yes)
	if w.IsManagerApprovalRequired < 0 || w.IsManagerApprovalRequired > 1 {
		return fmt.Errorf("invalid manager approval required value: %d (must be 0 or 1)", w.IsManagerApprovalRequired)
	}

	// Validate amount range if both are provided
	if w.MinAmount != nil && w.MaxAmount != nil {
		if *w.MinAmount > *w.MaxAmount {
			return ErrInvalidAmountRange
		}
	}

	// Validate required fields
	if w.CompanyID <= 0 {
		return errors.New("company_id is required")
	}
	if w.ApproverID <= 0 {
		return errors.New("approver_id is required")
	}

	return nil
}
