package db

import (
	"testing"
)

func TestDatabaseStore(t *testing.T) {
	// Create in-memory database store
	store, err := NewInMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create database store: %v", err)
	}
	defer store.Close()

	// Initialize database
	if err := store.Initialize(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Test company store
	t.Run("Company Store", func(t *testing.T) {
		// Get company by ID
		company, err := store.CompanyStore.GetByID(1)
		if err != nil {
			t.Fatalf("Failed to get company: %v", err)
		}
		if company.Name != "Light" {
			t.Errorf("Expected company name 'Light', got '%s'", company.Name)
		}
	})

	// Test database initialization
	t.Run("Database Initialization", func(t *testing.T) {
		// Test that we can get the database connection
		db := store.GetDB()
		if db == nil {
			t.Error("Database connection should not be nil")
		}
	})

	// Test database reset
	t.Run("Database Reset", func(t *testing.T) {
		if err := store.Reset(); err != nil {
			t.Fatalf("Failed to reset database: %v", err)
		}

		// After reset, company should not exist
		_, err := store.CompanyStore.GetByID(1)
		if err == nil {
			t.Error("Company should not exist after reset")
		}

		// Re-initialize
		if err := store.Initialize(); err != nil {
			t.Fatalf("Failed to re-initialize database: %v", err)
		}

		// Company should exist again
		company, err := store.CompanyStore.GetByID(1)
		if err != nil {
			t.Fatalf("Company should exist after re-initialization: %v", err)
		}
		if company.Name != "Light" {
			t.Errorf("Expected company name 'Light', got '%s'", company.Name)
		}
	})
}
