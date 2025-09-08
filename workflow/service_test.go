package workflow

import (
	"bufio"
	"strings"
	"testing"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/db"
)

func TestNewService(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			companyName        string
			companyDepartments []string
			db                 databaseService
			slack              notificationService
			email              notificationService
			options            []Option
		}
		wantErr bool
	}{
		{
			name: "valid service creation",
			input: struct {
				companyName        string
				companyDepartments []string
				db                 databaseService
				slack              notificationService
				email              notificationService
				options            []Option
			}{
				companyName:        "Test Company",
				companyDepartments: []string{"Engineering", "Sales"},
				db:                 &mockDatabaseService{},
				slack:              &mockNotificationService{},
				email:              &mockNotificationService{},
				options:            []Option{},
			},
			wantErr: false,
		},
		{
			name: "empty company name",
			input: struct {
				companyName        string
				companyDepartments []string
				db                 databaseService
				slack              notificationService
				email              notificationService
				options            []Option
			}{
				companyName:        "",
				companyDepartments: []string{"Engineering"},
				db:                 &mockDatabaseService{},
				slack:              &mockNotificationService{},
				email:              &mockNotificationService{},
				options:            []Option{},
			},
			wantErr: true,
		},
		{
			name: "empty departments",
			input: struct {
				companyName        string
				companyDepartments []string
				db                 databaseService
				slack              notificationService
				email              notificationService
				options            []Option
			}{
				companyName:        "Test Company",
				companyDepartments: []string{},
				db:                 &mockDatabaseService{},
				slack:              &mockNotificationService{},
				email:              &mockNotificationService{},
				options:            []Option{},
			},
			wantErr: true,
		},
		{
			name: "nil departments",
			input: struct {
				companyName        string
				companyDepartments []string
				db                 databaseService
				slack              notificationService
				email              notificationService
				options            []Option
			}{
				companyName:        "Test Company",
				companyDepartments: nil,
				db:                 &mockDatabaseService{},
				slack:              &mockNotificationService{},
				email:              &mockNotificationService{},
				options:            []Option{},
			},
			wantErr: true,
		},
		{
			name: "nil database service",
			input: struct {
				companyName        string
				companyDepartments []string
				db                 databaseService
				slack              notificationService
				email              notificationService
				options            []Option
			}{
				companyName:        "Test Company",
				companyDepartments: []string{"Engineering"},
				db:                 nil,
				slack:              &mockNotificationService{},
				email:              &mockNotificationService{},
				options:            []Option{},
			},
			wantErr: true,
		},
		{
			name: "nil slack service",
			input: struct {
				companyName        string
				companyDepartments []string
				db                 databaseService
				slack              notificationService
				email              notificationService
				options            []Option
			}{
				companyName:        "Test Company",
				companyDepartments: []string{"Engineering"},
				db:                 &mockDatabaseService{},
				slack:              nil,
				email:              &mockNotificationService{},
				options:            []Option{},
			},
			wantErr: true,
		},
		{
			name: "nil email service",
			input: struct {
				companyName        string
				companyDepartments []string
				db                 databaseService
				slack              notificationService
				email              notificationService
				options            []Option
			}{
				companyName:        "Test Company",
				companyDepartments: []string{"Engineering"},
				db:                 &mockDatabaseService{},
				slack:              &mockNotificationService{},
				email:              nil,
				options:            []Option{},
			},
			wantErr: true,
		},
		{
			name: "with custom logger option",
			input: struct {
				companyName        string
				companyDepartments []string
				db                 databaseService
				slack              notificationService
				email              notificationService
				options            []Option
			}{
				companyName:        "Test Company",
				companyDepartments: []string{"Engineering"},
				db:                 &mockDatabaseService{},
				slack:              &mockNotificationService{},
				email:              &mockNotificationService{},
				options:            []Option{WithLogger(&mockLogger{})},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewService(test.input.companyName, test.input.companyDepartments, test.input.db, test.input.slack, test.input.email, test.input.options...)

			if test.wantErr {
				if err == nil {
					t.Errorf("NewService() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("NewService() unexpected error: %v", err)
				return
			}
			if got == nil {
				t.Errorf("NewService() returned nil service")
			}
		})
	}
}

// Old tests removed - interface has changed to only have Run() method
// New tests for interactive functionality will be added below

func TestService_Run(t *testing.T) {
	// This test would require mocking stdin/stdout which is complex
	// For now, we'll test the individual methods that Run() calls
	t.Skip("Run() method requires interactive testing - testing individual methods instead")
}

func TestService_getUserInput(t *testing.T) {
	tests := []struct {
		name           string
		service        *service
		input          string
		expectedAmount float64
		expectedDept   string
		expectedMgr    bool
		wantErr        bool
	}{
		{
			name: "valid input",
			service: &service{
				company: company{
					name:        "Test Company",
					departments: []string{"Engineering", "Sales"},
				},
				reader: bufio.NewReader(strings.NewReader("100\nEngineering\ny\n")),
			},
			input:          "100\nEngineering\ny\n",
			expectedAmount: 100,
			expectedDept:   "Engineering",
			expectedMgr:    true,
			wantErr:        false,
		},
		{
			name: "empty input with defaults",
			service: &service{
				company: company{
					name:        "Test Company",
					departments: []string{"Engineering", "Sales"},
				},
				reader: bufio.NewReader(strings.NewReader("\n\n\n")),
			},
			input:          "\n\n\n",
			expectedAmount: 0,
			expectedDept:   "",
			expectedMgr:    false,
			wantErr:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.service.getUserInput()
			if test.wantErr {
				if err == nil {
					t.Errorf("getUserInput() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("getUserInput() unexpected error: %v", err)
				return
			}

			if test.service.userInput.amount != test.expectedAmount {
				t.Errorf("getUserInput() amount = %v, want %v", test.service.userInput.amount, test.expectedAmount)
			}
			if test.service.userInput.department != test.expectedDept {
				t.Errorf("getUserInput() department = %v, want %v", test.service.userInput.department, test.expectedDept)
			}
			if test.service.userInput.isManagerApprovalRequired != test.expectedMgr {
				t.Errorf("getUserInput() manager approval = %v, want %v", test.service.userInput.isManagerApprovalRequired, test.expectedMgr)
			}
		})
	}
}

func TestService_getInvoiceAmount(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "valid amount",
			input:    "100.50\n",
			expected: 100.50,
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    "\n",
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "invalid amount",
			input:    "abc\n",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "negative amount",
			input:    "-50\n",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &service{
				reader: bufio.NewReader(strings.NewReader(test.input)),
			}

			got, err := service.getInvoiceAmount()
			if test.wantErr {
				if err == nil {
					t.Errorf("getInvoiceAmount() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("getInvoiceAmount() unexpected error: %v", err)
				return
			}
			if got != test.expected {
				t.Errorf("getInvoiceAmount() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestService_getDepartment(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid department",
			input:    "Engineering\n",
			expected: "Engineering",
			wantErr:  false,
		},
		{
			name:     "case insensitive department",
			input:    "engineering\n",
			expected: "Engineering",
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    "\n",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "invalid department",
			input:    "InvalidDept\n",
			expected: "",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &service{
				company: company{
					name:        "Test Company",
					departments: []string{"Engineering", "Sales"},
				},
				reader: bufio.NewReader(strings.NewReader(test.input)),
			}

			got, err := service.getDepartment()
			if test.wantErr {
				if err == nil {
					t.Errorf("getDepartment() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("getDepartment() unexpected error: %v", err)
				return
			}
			if got != test.expected {
				t.Errorf("getDepartment() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestService_getManagerApprovalRequired(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		wantErr  bool
	}{
		{
			name:     "yes input",
			input:    "y\n",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "yes input uppercase",
			input:    "YES\n",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "no input",
			input:    "n\n",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    "\n",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "invalid input",
			input:    "maybe\n",
			expected: false,
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &service{
				reader: bufio.NewReader(strings.NewReader(test.input)),
			}

			got, err := service.getManagerApprovalRequired()
			if test.wantErr {
				if err == nil {
					t.Errorf("getManagerApprovalRequired() expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("getManagerApprovalRequired() unexpected error: %v", err)
				return
			}
			if got != test.expected {
				t.Errorf("getManagerApprovalRequired() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestService_askToContinue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "yes input",
			input:    "y\n",
			expected: true,
		},
		{
			name:     "no input",
			input:    "n\n",
			expected: false,
		},
		{
			name:     "yes input uppercase",
			input:    "YES\n",
			expected: true,
		},
		{
			name:     "no input uppercase",
			input:    "NO\n",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &service{
				reader: bufio.NewReader(strings.NewReader(test.input)),
			}

			got := service.askToContinue()
			if got != test.expected {
				t.Errorf("askToContinue() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestService_toInvoiceRequest(t *testing.T) {
	service := &service{
		company: company{
			name:        "Test Company",
			departments: []string{"Engineering", "Sales"},
		},
		userInput: userInput{
			amount:                    100.50,
			department:                "Engineering",
			isManagerApprovalRequired: true,
		},
	}

	expected := api.InvoiceRequest{
		CompanyName:               "Test Company",
		Amount:                    100.50,
		Department:                "Engineering",
		IsManagerApprovalRequired: true,
	}

	got := service.toInvoiceRequest()
	if got != expected {
		t.Errorf("toInvoiceRequest() = %v, want %v", got, expected)
	}
}

func TestService_getCompanyDepartments(t *testing.T) {
	service := &service{
		company: company{
			name:        "Test Company",
			departments: []string{"Engineering", "Sales", "Marketing"},
		},
	}

	expected := []string{"Engineering", "Sales", "Marketing"}
	got := service.getCompanyDepartments()

	if len(got) != len(expected) {
		t.Errorf("getCompanyDepartments() length = %v, want %v", len(got), len(expected))
		return
	}

	for i, dept := range got {
		if dept != expected[i] {
			t.Errorf("getCompanyDepartments()[%d] = %v, want %v", i, dept, expected[i])
		}
	}
}

// Mock implementations for testing
type mockDatabaseService struct {
	company     db.Company
	approver    db.Approver
	rule        db.WorkflowRule
	companyErr  error
	approverErr error
	ruleErr     error
}

func (m *mockDatabaseService) GetCompanyByName(name string) (db.Company, error) {
	if m.companyErr != nil {
		return db.Company{}, m.companyErr
	}
	return m.company, nil
}

func (m *mockDatabaseService) GetApproverByID(id int) (db.Approver, error) {
	if m.approverErr != nil {
		return db.Approver{}, m.approverErr
	}
	return m.approver, nil
}

func (m *mockDatabaseService) FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (db.WorkflowRule, error) {
	if m.ruleErr != nil {
		return db.WorkflowRule{}, m.ruleErr
	}
	return m.rule, nil
}

type mockNotificationService struct {
	response api.ApprovalResponse
	err      error
}

func (m *mockNotificationService) SendApprovalRequest(approvalRequest api.ApprovalRequest) (api.ApprovalResponse, error) {
	if m.err != nil {
		return api.ApprovalResponse{}, m.err
	}
	return m.response, nil
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...any)  {}
func (m *mockLogger) Error(msg string, args ...any) {}

// Helper functions for creating pointers
func floatPtr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}
