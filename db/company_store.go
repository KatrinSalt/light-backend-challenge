package db

import (
	"github.com/KatrinSalt/backend-challenge-go/db/models"
)

// CompanyStore defines the interface for company operations
type CompanyStore interface {
	GetByID(id int) (*models.Company, error)
}

// companyStore implements CompanyStore
type companyStore struct {
	db *client
}

// NewCompanyStore creates a new company store
func NewCompanyStore(db *client) CompanyStore {
	return &companyStore{db: db}
}

// GetByID retrieves a company by ID
func (s *companyStore) GetByID(id int) (*models.Company, error) {
	var company models.Company
	query := "SELECT id, name FROM companies WHERE id = ?"

	err := s.db.QueryRow(query, id).Scan(&company.ID, &company.Name)
	if err != nil {
		return nil, err
	}

	return &company, nil
}
