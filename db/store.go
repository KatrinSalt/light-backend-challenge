package db

import (
	"fmt"
	"log"
)

// Store provides a centralized interface for all database operations
type Store struct {
	db                *client
	CompanyStore      CompanyStore
	ApproverStore     ApproverStore
	WorkflowRuleStore WorkflowRuleStore
}

// NewStore creates a new database store with all individual stores
func NewStore(dbPath string) (*Store, error) {
	// Create database connection
	db, err := NewClient(func(o *ClientOptions) {
		o.DataSource = dbPath
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Initialize individual stores
	companyStore := NewCompanyStore(db)
	approverStore := NewApproverStore(db)
	workflowRuleStore := NewWorkflowRuleStore(db)

	store := &Store{
		db:                db,
		CompanyStore:      companyStore,
		ApproverStore:     approverStore,
		WorkflowRuleStore: workflowRuleStore,
	}

	return store, nil
}

// NewInMemoryStore creates a new in-memory database store
func NewInMemoryStore() (*Store, error) {
	return NewStore(":memory:")
}

// Initialize sets up the database schema and seeds it with sample data
func (s *Store) Initialize() error {
	log.Println("Initializing database...")

	// Initialize schema
	if err := s.db.CreateTables(); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Seed with sample data
	if err := s.db.SeedSampleData(); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// Reset drops all tables and recreates them with sample data
func (s *Store) Reset() error {
	return s.db.ResetDatabase()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// GetDB returns the underlying database connection (for advanced operations)
func (s *Store) GetDB() *client {
	return s.db
}
