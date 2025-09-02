package main

import (
	"fmt"
	"log"

	"github.com/KatrinSalt/backend-challenge-go/db"
)

func main() {
	fmt.Println("Invoice Approval Workflow - Database Layer Demo")
	fmt.Println("==============================================")

	// Create in-memory database store
	store, err := db.NewInMemoryStore()
	if err != nil {
		log.Fatalf("Failed to create database store: %v", err)
	}
	defer store.Close()

	// Initialize database with schema and sample data
	if err := store.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	fmt.Println("✅ Database initialized successfully!")

	// Demo: Get company by ID
	fmt.Println("\n📋 Company Details:")
	company, err := store.CompanyStore.GetByID(1)
	if err != nil {
		log.Printf("Failed to fetch company: %v", err)
	} else {
		fmt.Printf("  - ID: %d, Name: %s\n", company.ID, company.Name)
	}

	fmt.Println("🎉 Database layer demo completed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Implement the workflow service logic")
	fmt.Println("2. Create CLI interface for invoice processing")
	fmt.Println("3. Add approval channel services (Slack/Email)")
}
