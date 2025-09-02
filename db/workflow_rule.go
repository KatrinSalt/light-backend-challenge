package db

// WorkflowRule represents a rule that determines how invoices are approved.
type WorkflowRule struct {
	ID                        int      `db:"id"`
	CompanyID                 int      `db:"company_id"`
	MinAmount                 *float64 `db:"min_amount"`
	MaxAmount                 *float64 `db:"max_amount"`
	Department                *string  `db:"department"`
	IsManagerApprovalRequired *int     `db:"is_manager_approval_required"`
	ApproverID                int      `db:"approver_id"`
	ApprovalChannel           int      `db:"approval_channel"`
}
