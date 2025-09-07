package db

// Approver represents a person who can approve invoices.
type Approver struct {
	ID        int    `db:"id"`
	CompanyID int    `db:"company_id"`
	Name      string `db:"name"`
	Role      string `db:"role"`
	Email     string `db:"email"`
	SlackID   string `db:"slack_id"`
}
