package models

// Approver represents a person who can approve invoices
type Approver struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"company_id"`
	Email     string `json:"email"`
	SlackID   string `json:"slack_id"`
}
