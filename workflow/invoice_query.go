package workflow

// invoiceQuery represents the details of an invoice needed to find matching workflow rule.
type invoiceQuery struct {
	companyID                 int
	amount                    float64
	department                string
	isManagerApprovalRequired bool
}
