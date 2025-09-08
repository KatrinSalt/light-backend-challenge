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
	Update(approver Approver) error
	Delete(id int) error
	List(companyID int) ([]Approver, error)
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

	insert := fmt.Sprintf("INSERT INTO %s (company_id, name, role, email, slack_id) VALUES ($1, $2, $3, $4, $5)", s.table)
	if _, err := tx.Exec(insert, approver.CompanyID, approver.Name, approver.Role, approver.Email, approver.SlackID); err != nil {
		if strings.Contains(err.Error(), sql.SQLStateDuplicateKey) {
			return Approver{}, ErrApproverAlreadyExists
		}
		return Approver{}, err
	}

	// Get the created approver with its generated ID.
	var outApprover Approver
	query := fmt.Sprintf("SELECT id, company_id, name, role, email, slack_id FROM %s WHERE company_id = $1 AND email = $2", s.table)
	if err := tx.QueryRow(query, approver.CompanyID, approver.Email).Scan(&outApprover.ID, &outApprover.CompanyID, &outApprover.Name, &outApprover.Role, &outApprover.Email, &outApprover.SlackID); err != nil {
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
	query := fmt.Sprintf("SELECT id, company_id, name, role, email, slack_id FROM %s WHERE id = $1", s.table)
	err := s.client.QueryRow(query, id).Scan(&approver.ID, &approver.CompanyID, &approver.Name, &approver.Role, &approver.Email, &approver.SlackID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Approver{}, ErrApproverNotFound
		}
		return Approver{}, err
	}
	return approver, nil
}

// Update updates an existing approver.
func (s *approverStore) Update(approver Approver) error {
	tx, err := s.client.Transaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the approver exists
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", s.table)
	if err := tx.QueryRow(checkQuery, approver.ID).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check if approver exists: %w", err)
	}

	if !exists {
		return ErrApproverNotFound
	}

	// Update the approver
	updateQuery := fmt.Sprintf(`
		UPDATE %s 
		SET company_id = $2, name = $3, role = $4, email = $5, slack_id = $6 
		WHERE id = $1`, s.table)

	result, err := tx.Exec(updateQuery,
		approver.ID,
		approver.CompanyID,
		approver.Name,
		approver.Role,
		approver.Email,
		approver.SlackID)

	if err != nil {
		return fmt.Errorf("failed to update approver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrApproverNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// List retrieves all approvers for a specific company.
func (s *approverStore) List(companyID int) ([]Approver, error) {
	query := fmt.Sprintf("SELECT id, company_id, name, role, email, slack_id FROM %s WHERE company_id = $1", s.table)

	rows, err := s.client.Query(query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to query approvers by company ID: %w", err)
	}
	defer rows.Close()

	var approvers []Approver
	for rows.Next() {
		var approver Approver
		err := rows.Scan(&approver.ID, &approver.CompanyID, &approver.Name, &approver.Role, &approver.Email, &approver.SlackID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approver: %w", err)
		}
		approvers = append(approvers, approver)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over approver rows: %w", err)
	}

	return approvers, nil
}

// Delete deletes an approver by their ID.
func (s *approverStore) Delete(id int) error {
	tx, err := s.client.Transaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the approver exists
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", s.table)
	if err := tx.QueryRow(checkQuery, id).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check if approver exists: %w", err)
	}

	if !exists {
		return ErrApproverNotFound
	}

	// Delete the approver
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", s.table)
	result, err := tx.Exec(deleteQuery, id)

	if err != nil {
		return fmt.Errorf("failed to delete approver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrApproverNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
