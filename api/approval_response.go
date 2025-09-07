package api

// ApprovalResponse represents a response to an approval request.
type ApprovalResponse struct {
	ApproverName      string `json:"approver_name"`
	ApproverRole      string `json:"approver_role"`
	ApproverChannel   string `json:"approver_channel"`
	ApproverContactID string `json:"approver_contact_id"`
}
