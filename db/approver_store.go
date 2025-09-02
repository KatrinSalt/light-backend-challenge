package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/KatrinSalt/backend-challenge-go/db/sql"
)

var (
	ErrApproverNotFound      = errors.New("approver not found")
	ErrApproverAlreadyExists = errors.New("approver already exists")
)

// ApproverStore defines the interface for approver operations
type ApproverStore interface {
	Create(approver Approver) (Approver, error)
	GetByID(id int) (Approver, error)
}

// approverStore implements ApproverStore
type approverStore struct {
	client sql.Client
	table  string
}

// ApproverStoreOptions contains options for the approver store.
type ApproverStoreOptions struct {
	Table string
}

// ApproverStoreOption is a function that sets options on the approver store.
type ApproverStoreOption func(o *ApproverStoreOptions)

// NewApproverStore creates a new approver store
func NewApproverStore(client sql.Client, options ...ApproverStoreOption) (*approverStore, error) {
	if client == nil {
		return nil, errors.New("nil sql client")
	}

	opts := ApproverStoreOptions{}
	for _, option := range options {
		option(&opts)
	}
	if len(opts.Table) == 0 {
		opts.Table = defaultApproverTable
	}

	return &approverStore{
		client: client,
		table:  opts.Table,
	}, nil
}

// Create creates a new approver.
func (s *approverStore) Create(approver Approver) (Approver, error) {
	tx, err := s.client.Transaction()
	if err != nil {
		return Approver{}, err
	}
	defer tx.Rollback()

	insert := fmt.Sprintf("INSERT INTO %s (company_id, email, slack_id) VALUES ($1, $2, $3)", s.table)
	if _, err := tx.Exec(insert, approver.CompanyID, approver.Email, approver.SlackID); err != nil {
		if strings.Contains(err.Error(), sql.SQLStateDuplicateKey) {
			return Approver{}, ErrApproverAlreadyExists
		}
		return Approver{}, err
	}

	// Get the created approver with its generated ID.
	var outApprover Approver
	query := fmt.Sprintf("SELECT id, company_id, email, slack_id FROM %s WHERE company_id = $1 AND email = $2", s.table)
	if err := tx.QueryRow(query, approver.CompanyID, approver.Email).Scan(&outApprover.ID, &outApprover.CompanyID, &outApprover.Email, &outApprover.SlackID); err != nil {
		return Approver{}, err
	}

	if err := tx.Commit(); err != nil {
		return Approver{}, err
	}

	return outApprover, nil
}

// GetByID retrieves an approver by their ID.
func (s *approverStore) GetByID(id int) (Approver, error) {
	var approver Approver
	query := fmt.Sprintf("SELECT id, company_id, email, slack_id FROM %s WHERE id = $1", s.table)
	err := s.client.QueryRow(query, id).Scan(&approver.ID, &approver.CompanyID, &approver.Email, &approver.SlackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Approver{}, ErrApproverNotFound
		}
		return Approver{}, err
	}
	return approver, nil
}
