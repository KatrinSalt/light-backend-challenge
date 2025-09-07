package api

// ApprovalRequest represents an approval request to be sent.
type ApprovalRequest struct {
	Approver Approver       `json:"approver"`
	Invoice  InvoiceDetails `json:"invoice"`
}

func (a *ApprovalRequest) Validate() error {
	if err := a.Approver.validate(); err != nil {
		return err
	}
	return nil
}
