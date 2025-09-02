package db

import (
	"github.com/KatrinSalt/backend-challenge-go/db/models"
)

// Repository interfaces - kept for backward compatibility if needed
type CompanyRepository interface {
	GetByID(id int) (*models.Company, error)
}

type ApproverRepository interface {
	// Will add methods when needed
}

type WorkflowRuleRepository interface {
	// Will add methods when needed
}
