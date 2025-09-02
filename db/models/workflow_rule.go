package models

// WorkflowRule represents a rule that determines who should approve an invoice
type WorkflowRule struct {
	ID                        int      `json:"id"`
	CompanyID                 int      `json:"company_id"`
	MinAmount                 *float64 `json:"min_amount"`
	MaxAmount                 *float64 `json:"max_amount"`
	Department                *string  `json:"department"`
	IsManagerApprovalRequired *int     `json:"is_manager_approval_required"` // 0 = false, 1 = true
	ApproverID                int      `json:"approver_id"`
	ApprovalChannel           int      `json:"approval_channel"` // 0 = Slack, 1 = Email
}
