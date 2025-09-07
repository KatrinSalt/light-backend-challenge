package api

import "errors"

var (
	ErrMissingContactID = errors.New("slack_id or email is missing")
)

// Approver represents a person who can approve invoices.
type Approver struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Email   string `json:"email"`
	SlackID string `json:"slack_id"`
}

func (a *Approver) validate() error {
	if a.Email == "" && a.SlackID == "" {
		return ErrMissingContactID
	}
	return nil
}
