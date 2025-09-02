package db

import (
	"fmt"
	"log"
)

// CreateTables creates all the necessary tables in the database.
func (client *client) CreateTables() error {
	// TODO: check logging strategy later on.
	log.Println("Initializing database schema...")

	queries := []string{
		// companies table.
		`CREATE TABLE IF NOT EXISTS companies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		// approvers table.
		`CREATE TABLE IF NOT EXISTS approvers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_id INTEGER NOT NULL,
			email TEXT NOT NULL,
			slack_id TEXT NOT NULL,
			FOREIGN KEY (company_id) REFERENCES companies (id),
			UNIQUE(company_id, email),
			UNIQUE(company_id, slack_id)
		)`,
		// workflow rules table.
		`CREATE TABLE IF NOT EXISTS workflow_rules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_id INTEGER NOT NULL,
			min_amount REAL,
			max_amount REAL,
			department TEXT,
			is_manager_approval_required INTEGER DEFAULT 0 CHECK (is_manager_approval_required IN (0, 1)),
			approver_id INTEGER NOT NULL,
			approval_channel INTEGER NOT NULL,
			FOREIGN KEY (company_id) REFERENCES companies (id),
			FOREIGN KEY (approver_id) REFERENCES approvers (id)
		)`,
	}

	for i, query := range queries {
		if _, err := client.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %d: %w", i+1, err)
		}
	}

	// TODO: check logging strategy later on.
	log.Println("Database schema initialized successfully")

	return nil
}

// ResetDatabase drops all tables and recreates them
func (client *client) ResetDatabase() error {
	log.Println("Resetting database...")

	// Drop all tables
	dropQueries := []string{
		"DROP TABLE IF EXISTS workflow_rules",
		"DROP TABLE IF EXISTS approvers", 
		"DROP TABLE IF EXISTS companies",
	}

	for i, query := range dropQueries {
		if _, err := client.Exec(query); err != nil {
			return fmt.Errorf("failed to execute drop query %d: %w", i+1, err)
		}
	}

	// Recreate tables
	if err := client.CreateTables(); err != nil {
		return fmt.Errorf("failed to recreate tables: %w", err)
	}

	log.Println("Database reset successfully")
	return nil
}
