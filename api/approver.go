package api

import "errors"

var (
	ErrMissingEmail   = errors.New("email is missing")
	ErrMissingSlackID = errors.New("slack_id is missing")
)

// Approver represents a person who can approve invoices.
type Approver struct {
	ID        int    `json:"id,omitempty"`
	CompanyID int    `json:"company_id,omitempty"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	SlackID   string `json:"slack_id"`
}

func (a *Approver) Validate() error {
	if a.Email == "" {
		return ErrMissingEmail
	}
	if a.SlackID == "" {
		return ErrMissingSlackID
	}

	return nil
}
