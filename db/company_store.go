package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/KatrinSalt/backend-challenge-go/db/sql"
)

var (
	ErrCompanyNotFound      = errors.New("company not found")
	ErrCompanyAlreadyExists = errors.New("company already exists")
)

// CompanyStore defines the interface for company operations
type CompanyStore interface {
	Create(company Company) (Company, error)
	GetByID(id int) (Company, error)
	GetByName(name string) (Company, error)
}

// companyStore implements CompanyStore
type companyStore struct {
	client sql.Client
	table  string
}

// CompanyStoreOptions contains options for the company store.
type CompanyStoreOptions struct {
	Table string
}

// CompanyStoreOption is a function that sets options on the company store.
type CompanyStoreOption func(o *CompanyStoreOptions)

// NewCompanyStore creates a new company store
func NewCompanyStore(client sql.Client, options ...CompanyStoreOption) (*companyStore, error) {
	if client == nil {
		return nil, errors.New("nil sql client")
	}

	opts := CompanyStoreOptions{}
	for _, option := range options {
		option(&opts)
	}
	if len(opts.Table) == 0 {
		opts.Table = defaultCompanyTable
	}

	return &companyStore{
		client: client,
		table:  opts.Table,
	}, nil
}

// Create creates a new company.
func (s *companyStore) Create(company Company) (Company, error) {
	tx, err := s.client.Transaction()
	if err != nil {
		return Company{}, err
	}
	defer tx.Rollback()

	insert := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1)", s.table)
	if _, err := tx.Exec(insert, company.Name); err != nil {
		if strings.Contains(err.Error(), sql.SQLStateDuplicateKey) {
			return Company{}, ErrCompanyAlreadyExists
		}
		return Company{}, err
	}

	// Get the created company with its generated ID.
	var outCompany Company
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE name = $1", s.table)
	if err := tx.QueryRow(query, company.Name).Scan(&outCompany.ID, &outCompany.Name); err != nil {
		return Company{}, err
	}

	if err := tx.Commit(); err != nil {
		return Company{}, err
	}

	return outCompany, nil
}

// GetByID retrieves a company by its ID.
func (s *companyStore) GetByID(id int) (Company, error) {
	var company Company
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE id = $1", s.table)
	err := s.client.QueryRow(query, id).Scan(&company.ID, &company.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Company{}, ErrCompanyNotFound
		}
		return Company{}, err
	}
	return company, nil
}

// GetByName retrieves a company by its name.
func (s *companyStore) GetByName(name string) (Company, error) {
	var company Company
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE name = $1", s.table)
	err := s.client.QueryRow(query, name).Scan(&company.ID, &company.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Company{}, ErrCompanyNotFound
		}
		return Company{}, err
	}
	return company, nil
}
