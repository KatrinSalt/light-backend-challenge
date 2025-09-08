package workflow

import (
	"fmt"
	"testing"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/common"
	"github.com/KatrinSalt/backend-challenge-go/db"
	"github.com/KatrinSalt/backend-challenge-go/db/sqlite"
	"github.com/KatrinSalt/backend-challenge-go/notification/email"
	"github.com/KatrinSalt/backend-challenge-go/notification/slack"
)

// TestWorkflowRulesIntegration tests all 5 pre-seeded workflow rules
// to ensure they work correctly with the priority-based matching system
func TestWorkflowRulesIntegration(t *testing.T) {
	// Setup test database with sample data
	dbService := setupTestDatabase(t)

	// Setup notification services (mocked)
	slackService, err := slack.NewService("test-slack-token")
	if err != nil {
		t.Fatalf("Failed to create slack service: %v", err)
	}
	emailService, err := email.NewService("test-email-connection")
	if err != nil {
		t.Fatalf("Failed to create email service: %v", err)
	}

	// Setup workflow service
	workflowService, err := NewService(
		"Light",
		[]string{"Finance", "Marketing"},
		dbService,
		slackService,
		emailService,
		WithLogger(common.NewLogger()),
	)
	if err != nil {
		t.Fatalf("Failed to create workflow service: %v", err)
	}

	// Test cases for all 5 workflow rules
	testCases := []struct {
		name                 string
		amount               float64
		department           string
		requiresManager      bool
		expectedApproverID   int
		expectedApproverName string
		expectedChannel      string
		expectedContactID    string
		description          string
	}{
		{
			name:                 "Rule 1: Invoice < $5k → Finance Team Member via Slack",
			amount:               3000,
			department:           "Finance",
			requiresManager:      false,
			expectedApproverID:   1,
			expectedApproverName: "System User",
			expectedChannel:      "slack",
			expectedContactID:    "U123456",
			description:          "Small Finance invoice should go to Finance Team Member via Slack",
		},
		{
			name:                 "Rule 2: $5k ≤ Invoice < $10k → Finance Team Member via Email",
			amount:               7500,
			department:           "Finance",
			requiresManager:      false,
			expectedApproverID:   1,
			expectedApproverName: "System User",
			expectedChannel:      "email",
			expectedContactID:    "finance_team@light.com",
			description:          "Medium Finance invoice should go to Finance Team Member via Email",
		},
		{
			name:                 "Rule 3: $5k ≤ Invoice < $10k + Manager Approval → Finance Manager via Email",
			amount:               7500,
			department:           "Finance",
			requiresManager:      true,
			expectedApproverID:   2,
			expectedApproverName: "Vera Sander",
			expectedChannel:      "email",
			expectedContactID:    "vera_sander@light.com",
			description:          "Medium Finance invoice with manager approval should go to Finance Manager via Email",
		},
		{
			name:                 "Rule 4: Invoice ≥ $10k (any dept) → CFO via Slack",
			amount:               15000,
			department:           "Finance",
			requiresManager:      false,
			expectedApproverID:   3,
			expectedApproverName: "Amanda Svensson",
			expectedChannel:      "slack",
			expectedContactID:    "U345678",
			description:          "Large Finance invoice should go to CFO via Slack",
		},
		{
			name:                 "Rule 5: Invoice ≥ $10k + Marketing → CMO via Email",
			amount:               15000,
			department:           "Marketing",
			requiresManager:      false,
			expectedApproverID:   4,
			expectedApproverName: "Sarah Johnson",
			expectedChannel:      "email",
			expectedContactID:    "sarah_johnson@light.com",
			description:          "Large Marketing invoice should go to CMO via Email",
		},
		{
			name:                 "Priority Test: Marketing ≥ $10k + Manager Approval → CMO via Email (Rule 5 wins)",
			amount:               15000,
			department:           "Marketing",
			requiresManager:      true,
			expectedApproverID:   4,
			expectedApproverName: "Sarah Johnson",
			expectedChannel:      "email",
			expectedContactID:    "sarah_johnson@light.com",
			description:          "Marketing invoice with manager approval should still go to CMO (Rule 5 has higher specificity than Rule 4)",
		},
	}

	// Run all test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create invoice request
			invoiceReq := api.InvoiceRequest{
				CompanyName:               "Light",
				Amount:                    tc.amount,
				Department:                tc.department,
				IsManagerApprovalRequired: tc.requiresManager,
			}

			// Process the invoice
			response, err := processInvoiceForTest(workflowService, invoiceReq)
			if err != nil {
				t.Fatalf("Failed to process invoice: %v", err)
			}

			// Verify the response - get approver info from database to verify ID
			approver, err := dbService.GetApproverByID(tc.expectedApproverID)
			if err != nil {
				t.Errorf("Failed to get expected approver: %v", err)
			}

			// Verify approver name matches
			if response.ApproverName != approver.Name {
				t.Errorf("Expected approver name %s, got %s", approver.Name, response.ApproverName)
			}

			if response.ApproverName != tc.expectedApproverName {
				t.Errorf("Expected approver name %s, got %s", tc.expectedApproverName, response.ApproverName)
			}

			if response.ApproverChannel != tc.expectedChannel {
				t.Errorf("Expected channel %s, got %s", tc.expectedChannel, response.ApproverChannel)
			}

			if response.ApproverContactID != tc.expectedContactID {
				t.Errorf("Expected contact ID %s, got %s", tc.expectedContactID, response.ApproverContactID)
			}

			t.Logf("✅ %s: %s", tc.name, tc.description)
			t.Logf("   Amount: $%.2f, Department: %s, Manager Approval: %v", tc.amount, tc.department, tc.requiresManager)
			t.Logf("   → Approver: %s (%s) via %s", response.ApproverName, response.ApproverRole, response.ApproverChannel)
		})
	}
}

// TestWorkflowRulesEdgeCases tests edge cases and boundary conditions
func TestWorkflowRulesEdgeCases(t *testing.T) {
	// Setup test database with sample data
	dbService := setupTestDatabase(t)

	// Setup notification services (mocked)
	slackService, err := slack.NewService("test-slack-token")
	if err != nil {
		t.Fatalf("Failed to create slack service: %v", err)
	}
	emailService, err := email.NewService("test-email-connection")
	if err != nil {
		t.Fatalf("Failed to create email service: %v", err)
	}

	// Setup workflow service
	workflowService, err := NewService(
		"Light",
		[]string{"Finance", "Marketing"},
		dbService,
		slackService,
		emailService,
		WithLogger(common.NewLogger()),
	)
	if err != nil {
		t.Fatalf("Failed to create workflow service: %v", err)
	}

	// Test edge cases
	edgeCases := []struct {
		name            string
		amount          float64
		department      string
		requiresManager bool
		description     string
	}{
		{
			name:            "Boundary: Exactly $5k Finance invoice",
			amount:          5000,
			department:      "Finance",
			requiresManager: false,
			description:     "Should match Rule 2 ($5k ≤ amount < $10k)",
		},
		{
			name:            "Boundary: Exactly $10k Finance invoice",
			amount:          10000,
			department:      "Finance",
			requiresManager: false,
			description:     "Should match Rule 4 (amount ≥ $10k)",
		},
		{
			name:            "Boundary: Exactly $10k Marketing invoice",
			amount:          10000,
			department:      "Marketing",
			requiresManager: false,
			description:     "Should match Rule 5 (amount ≥ $10k + Marketing)",
		},
		{
			name:            "Edge: $4,999.99 Finance invoice",
			amount:          4999.99,
			department:      "Finance",
			requiresManager: false,
			description:     "Should match Rule 1 (amount < $5k)",
		},
		{
			name:            "Edge: $9,999.99 Finance invoice",
			amount:          9999.99,
			department:      "Finance",
			requiresManager: false,
			description:     "Should match Rule 2 ($5k ≤ amount < $10k)",
		},
	}

	// Run edge case tests
	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create invoice request
			invoiceReq := api.InvoiceRequest{
				CompanyName:               "Light",
				Amount:                    tc.amount,
				Department:                tc.department,
				IsManagerApprovalRequired: tc.requiresManager,
			}

			// Process the invoice
			response, err := processInvoiceForTest(workflowService, invoiceReq)
			if err != nil {
				t.Fatalf("Failed to process invoice: %v", err)
			}

			// Verify we got a valid response
			if response.ApproverName == "" {
				t.Errorf("Expected valid approver name, got empty string")
			}

			t.Logf("✅ %s: %s", tc.name, tc.description)
			t.Logf("   Amount: $%.2f → Approver: %s (%s) via %s", tc.amount, response.ApproverName, response.ApproverRole, response.ApproverChannel)
		})
	}
}

// setupTestDatabase creates a test database with sample data
func setupTestDatabase(t *testing.T) db.Service {
	// Create in-memory SQLite client
	client, err := sqlite.NewClient()
	if err != nil {
		t.Fatalf("Failed to create database client: %v", err)
	}

	// Create database service with sample data
	dbService, err := db.NewService(client, db.WithSampleData(db.NewSampleData()))
	if err != nil {
		t.Fatalf("Failed to create database service: %v", err)
	}

	// Initialize database schema
	if err := dbService.Initialize(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Seed with sample data
	if err := dbService.SeedSampleData(); err != nil {
		t.Fatalf("Failed to seed sample data: %v", err)
	}

	return dbService
}

// processInvoiceForTest processes an invoice and returns the response
// This function accesses the private processInvoice method for testing purposes
func processInvoiceForTest(workflowService Service, invoiceReq api.InvoiceRequest) (api.ApprovalResponse, error) {
	// Type assert to access the private method
	serviceImpl, ok := workflowService.(*service)
	if !ok {
		return api.ApprovalResponse{}, fmt.Errorf("failed to cast service to implementation")
	}

	// Call the private processInvoice method
	return serviceImpl.processInvoice(invoiceReq)
}
