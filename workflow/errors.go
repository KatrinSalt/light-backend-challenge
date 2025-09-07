package workflow

import "errors"

var (
	// ErrCompanyNotFound is returned when a company is not found in the system.
	ErrCompanyNotFound = errors.New("company not found")
	// ErrWorkflowRuleNotFound is returned when a workflow rule is not found in the system.
	ErrWorkflowRuleNotFound = errors.New("workflow rule not found")
	// ErrApproverNotFound is returned when an approver is not found in the system.
	ErrApproverNotFound = errors.New("approver not found")
	// ErrUnsupportedApprovalChannel is returned when an unsupported approval channel is used.
	ErrUnsupportedApprovalChannel = errors.New("unsupported approval channel")
)
