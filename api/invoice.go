package api

// InvoiceRequest represents an invoice that needs approval.
type InvoiceRequest struct {
	CompanyName               string  `json:"company_name"`
	Amount                    float64 `json:"amount"`
	Department                string  `json:"department,omitempty"`
	IsManagerApprovalRequired bool    `json:"is_manager_approval_required,omitempty"`
}

type InvoiceDetails struct {
	Amount float64 `json:"amount"`
}
