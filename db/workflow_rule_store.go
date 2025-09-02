package db

// WorkflowRuleStore defines the interface for workflow rule operations
type WorkflowRuleStore interface {
	// Will add methods when needed
}

// workflowRuleStore implements WorkflowRuleStore
type workflowRuleStore struct {
	db *client
}

// NewWorkflowRuleStore creates a new workflow rule store
func NewWorkflowRuleStore(db *client) WorkflowRuleStore {
	return &workflowRuleStore{db: db}
}
