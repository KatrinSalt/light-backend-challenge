package db

// ApproverStore defines the interface for approver operations
type ApproverStore interface {
	// Will add methods when needed
}

// approverStore implements ApproverStore
type approverStore struct {
	db *client
}

// NewApproverStore creates a new approver store
func NewApproverStore(db *client) ApproverStore {
	return &approverStore{db: db}
}
