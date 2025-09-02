package main

import (
	"fmt"
	"log"

	"github.com/KatrinSalt/backend-challenge-go/db"
	"github.com/KatrinSalt/backend-challenge-go/db/sqlite"
)

func main() {
	fmt.Println("Invoice Approval Workflow - Database Layer Demo")
	fmt.Println("==============================================")

	// Create SQLite client
	sqlClient, err := sqlite.NewClient()
	if err != nil {
		log.Fatalf("Failed to create SQLite client: %v", err)
	}

	// Create database service
	dbService, err := db.NewDBService(sqlClient)
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}

	fmt.Println("✅ Database service created successfully!")

	// Initialize database schema
	schemaQueries := []string{
		`CREATE TABLE IF NOT EXISTS companies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS approvers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			company_id INTEGER NOT NULL,
			email TEXT NOT NULL,
			slack_id TEXT NOT NULL,
			FOREIGN KEY (company_id) REFERENCES companies (id),
			UNIQUE(company_id, email),
			UNIQUE(company_id, slack_id)
		)`,
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

	if err := dbService.Initialize(schemaQueries); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	fmt.Println("✅ Database schema initialized successfully!")

	// Seed with sample data
	if err := dbService.SeedSampleData(); err != nil {
		log.Fatalf("Failed to seed sample data: %v", err)
	}
	fmt.Println("✅ Sample data seeded successfully!")

	// Demo: Show all populated data
	fmt.Println("\n📊 Database Content Summary:")
	fmt.Println("=============================")

	// Show companies
	fmt.Println("\n🏢 Companies:")
	companies := dbService.GetSampleData().Companies
	for _, company := range companies {
		fmt.Printf("  - ID: %d, Name: %s\n", company.ID, company.Name)
	}

	// Show approvers
	fmt.Println("\n👥 Approvers:")
	approvers := dbService.GetSampleData().Approvers
	for _, approver := range approvers {
		fmt.Printf("  - ID: %d, Company ID: %d, Email: %s, Slack ID: %s\n",
			approver.ID, approver.CompanyID, approver.Email, approver.SlackID)
	}

	// Show workflow rules
	fmt.Println("\n⚙️  Workflow Rules:")
	rules := dbService.GetSampleData().WorkflowRules
	for i, rule := range rules {
		fmt.Printf("  Rule %d:\n", i+1)
		fmt.Printf("    - Company ID: %d\n", rule.CompanyID)
		if rule.MinAmount != nil {
			fmt.Printf("    - Min Amount: $%.2f\n", *rule.MinAmount)
		}
		if rule.MaxAmount != nil {
			fmt.Printf("    - Max Amount: $%.2f\n", *rule.MaxAmount)
		}
		if rule.Department != nil {
			fmt.Printf("    - Department: %s\n", *rule.Department)
		}
		if rule.IsManagerApprovalRequired != nil {
			requires := "Yes"
			if *rule.IsManagerApprovalRequired == 0 {
				requires = "No"
			}
			fmt.Printf("    - Manager Approval Required: %s\n", requires)
		}
		fmt.Printf("    - Approver ID: %d\n", rule.ApproverID)
		channel := "Slack"
		if rule.ApprovalChannel == 1 {
			channel = "Email"
		}
		fmt.Printf("    - Approval Channel: %s\n", channel)
		fmt.Println()
	}

	// Test data retrieval from stores
	fmt.Println("🔍 Testing Data Retrieval:")
	fmt.Println("==========================")

	// Test company retrieval
	fmt.Println("\n🏢 Testing Company Store:")
	company, err := dbService.CompanyStore.GetByID(1)
	if err != nil {
		log.Printf("❌ Failed to fetch company: %v", err)
	} else {
		fmt.Printf("✅ Company retrieved: ID: %d, Name: %s\n", company.ID, company.Name)
	}

	// Test approver retrieval
	fmt.Println("\n👥 Testing Approver Store:")
	approver, err := dbService.ApproverStore.GetByID(1)
	if err != nil {
		log.Printf("❌ Failed to fetch approver: %d: %v", 1, err)
	} else {
		fmt.Printf("✅ Approver retrieved: ID: %d, Email: %s\n", approver.ID, approver.Email)
	}

	// Test workflow rule retrieval
	fmt.Println("\n⚙️  Testing Workflow Rule Store:")
	fmt.Println("✅ Workflow rule store created successfully (GetByCompanyID not implemented yet)")

	fmt.Println("\n🎉 Database layer demo completed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Implement the workflow service logic")
	fmt.Println("2. Create CLI interface for invoice processing")
	fmt.Println("3. Add approval channel services (Slack/Email)")
	fmt.Println("4. Implement invoice approval workflow logic")
}
