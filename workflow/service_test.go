package workflow

import (
	"errors"
	"testing"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/db"
)

func TestNewService(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			company Company
			db      databaseService
			slack   notificationService
			email   notificationService
			options []Option
		}
		wantErr bool
	}{
		{
			name: "valid service creation",
			input: struct {
				company Company
				db      databaseService
				slack   notificationService
				email   notificationService
				options []Option
			}{
				company: Company{Name: "Test Company", Departments: []string{"Engineering", "Sales"}},
				db:      &mockDatabaseService{},
				slack:   &mockNotificationService{},
				email:   &mockNotificationService{},
				options: []Option{},
			},
			wantErr: false,
		},
		{
			name: "empty company name",
			input: struct {
				company Company
				db      databaseService
				slack   notificationService
				email   notificationService
				options []Option
			}{
				company: Company{Name: "", Departments: []string{"Engineering"}},
				db:      &mockDatabaseService{},
				slack:   &mockNotificationService{},
				email:   &mockNotificationService{},
				options: []Option{},
			},
			wantErr: true,
		},
		{
			name: "empty departments",
			input: struct {
				company Company
				db      databaseService
				slack   notificationService
				email   notificationService
				options []Option
			}{
				company: Company{Name: "Test Company", Departments: []string{}},
				db:      &mockDatabaseService{},
				slack:   &mockNotificationService{},
				email:   &mockNotificationService{},
				options: []Option{},
			},
			wantErr: true,
		},
		{
			name: "nil departments",
			input: struct {
				company Company
				db      databaseService
				slack   notificationService
				email   notificationService
				options []Option
			}{
				company: Company{Name: "Test Company", Departments: nil},
				db:      &mockDatabaseService{},
				slack:   &mockNotificationService{},
				email:   &mockNotificationService{},
				options: []Option{},
			},
			wantErr: true,
		},
		{
			name: "nil database service",
			input: struct {
				company Company
				db      databaseService
				slack   notificationService
				email   notificationService
				options []Option
			}{
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db:      nil,
				slack:   &mockNotificationService{},
				email:   &mockNotificationService{},
				options: []Option{},
			},
			wantErr: true,
		},
		{
			name: "nil slack service",
			input: struct {
				company Company
				db      databaseService
				slack   notificationService
				email   notificationService
				options []Option
			}{
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db:      &mockDatabaseService{},
				slack:   nil,
				email:   &mockNotificationService{},
				options: []Option{},
			},
			wantErr: true,
		},
		{
			name: "nil email service",
			input: struct {
				company Company
				db      databaseService
				slack   notificationService
				email   notificationService
				options []Option
			}{
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db:      &mockDatabaseService{},
				slack:   &mockNotificationService{},
				email:   nil,
				options: []Option{},
			},
			wantErr: true,
		},
		{
			name: "with custom logger option",
			input: struct {
				company Company
				db      databaseService
				slack   notificationService
				email   notificationService
				options []Option
			}{
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db:      &mockDatabaseService{},
				slack:   &mockNotificationService{},
				email:   &mockNotificationService{},
				options: []Option{WithLogger(&mockLogger{})},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewService(test.input.company, test.input.db, test.input.slack, test.input.email, test.input.options...)

			if test.wantErr {
				if err == nil {
					t.Errorf("NewService() expected error but got none")
				}
				if got != nil {
					t.Errorf("NewService() expected nil service on error but got %v", got)
				}
			} else {
				if err != nil {
					t.Errorf("NewService() unexpected error: %v", err)
				}
				if got == nil {
					t.Errorf("NewService() expected service but got nil")
				}
			}
		})
	}
}

func TestService_ProcessInvoice(t *testing.T) {
	var tests = []struct {
		name        string
		service     *service
		invoice     api.InvoiceRequest
		wantResp    api.ApprovalResponse
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful slack approval request",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "John Doe",
						Role:      "Manager",
						Email:     "john@example.com",
						SlackID:   "U123456",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           0, // Slack
					},
				},
				slack: &mockNotificationService{
					response: api.ApprovalResponse{
						ApproverName:      "John Doe",
						ApproverRole:      "Manager",
						ApproverChannel:   "slack",
						ApproverContactID: "U123456",
					},
				},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantResp: api.ApprovalResponse{
				ApproverName:      "John Doe",
				ApproverRole:      "Manager",
				ApproverChannel:   "slack",
				ApproverContactID: "U123456",
			},
			wantErr: false,
		},
		{
			name: "successful email approval request",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "Jane Smith",
						Role:      "Director",
						Email:     "jane@example.com",
						SlackID:   "U789012",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           1, // Email
					},
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{
					response: api.ApprovalResponse{
						ApproverName:      "Jane Smith",
						ApproverRole:      "Director",
						ApproverChannel:   "email",
						ApproverContactID: "jane@example.com",
					},
				},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    750.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantResp: api.ApprovalResponse{
				ApproverName:      "Jane Smith",
				ApproverRole:      "Director",
				ApproverChannel:   "email",
				ApproverContactID: "jane@example.com",
			},
			wantErr: false,
		},
		{
			name: "company not found",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					companyErr: errors.New("company not found"),
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Non-existent Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantResp: api.ApprovalResponse{},
			wantErr:  true,
		},
		{
			name: "workflow rule not found",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					ruleErr: errors.New("workflow rule not found"),
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantResp: api.ApprovalResponse{},
			wantErr:  true,
		},
		{
			name: "approver not found",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                999, // Non-existent approver
						ApprovalChannel:           0,
					},
					approverErr: errors.New("approver not found"),
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantResp: api.ApprovalResponse{},
			wantErr:  true,
		},
		{
			name: "unsupported approval channel",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "John Doe",
						Role:      "Manager",
						Email:     "john@example.com",
						SlackID:   "U123456",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           99, // Unsupported channel
					},
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantResp:    api.ApprovalResponse{},
			wantErr:     true,
			expectedErr: ErrUnsupportedApprovalChannel,
		},
		{
			name: "notification service error",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "John Doe",
						Role:      "Manager",
						Email:     "john@example.com",
						SlackID:   "U123456",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           0,
					},
				},
				slack: &mockNotificationService{
					err: errors.New("failed to send slack message"),
				},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantResp: api.ApprovalResponse{},
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.service.ProcessInvoice(test.invoice)

			if test.wantErr {
				if err == nil {
					t.Errorf("ProcessInvoice() expected error but got none")
				}
				if test.expectedErr != nil && !errors.Is(err, test.expectedErr) {
					t.Errorf("ProcessInvoice() expected error %v but got %v", test.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("ProcessInvoice() unexpected error: %v", err)
				}
				if got != test.wantResp {
					t.Errorf("ProcessInvoice() = %v, want %v", got, test.wantResp)
				}
			}
		})
	}
}

func TestService_ValidateCompany(t *testing.T) {
	var tests = []struct {
		name    string
		service *service
		wantErr bool
	}{
		{
			name: "valid company",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
				},
			},
			wantErr: false,
		},
		{
			name: "company not found",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Non-existent Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					companyErr: errors.New("company not found"),
				},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.service.ValidateCompany()

			if test.wantErr {
				if err == nil {
					t.Errorf("ValidateCompany() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCompany() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestService_ProcessInvoice_EdgeCases(t *testing.T) {
	var tests = []struct {
		name        string
		service     *service
		invoice     api.InvoiceRequest
		wantErr     bool
		expectedErr error
	}{
		{
			name: "zero amount invoice",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "John Doe",
						Role:      "Manager",
						Email:     "john@example.com",
						SlackID:   "U123456",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(0.0),
						MaxAmount:                 floatPtr(100.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           0,
					},
				},
				slack: &mockNotificationService{
					response: api.ApprovalResponse{
						ApproverName:      "John Doe",
						ApproverRole:      "Manager",
						ApproverChannel:   "slack",
						ApproverContactID: "U123456",
					},
				},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    0.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantErr: false,
		},
		{
			name: "very large amount invoice",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "John Doe",
						Role:      "Manager",
						Email:     "john@example.com",
						SlackID:   "U123456",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(1000000.0),
						MaxAmount:                 nil, // No upper limit
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(1),
						ApproverID:                1,
						ApprovalChannel:           1,
					},
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{
					response: api.ApprovalResponse{
						ApproverName:      "John Doe",
						ApproverRole:      "Manager",
						ApproverChannel:   "email",
						ApproverContactID: "john@example.com",
					},
				},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    999999999.99,
				Department:                "Engineering",
				IsManagerApprovalRequired: true,
			},
			wantErr: false,
		},
		{
			name: "empty department",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "John Doe",
						Role:      "Manager",
						Email:     "john@example.com",
						SlackID:   "U123456",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                nil, // No department restriction
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           0,
					},
				},
				slack: &mockNotificationService{
					response: api.ApprovalResponse{
						ApproverName:      "John Doe",
						ApproverRole:      "Manager",
						ApproverChannel:   "slack",
						ApproverContactID: "U123456",
					},
				},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "", // Empty department
				IsManagerApprovalRequired: false,
			},
			wantErr: false,
		},
		{
			name: "approver with both email and slack contact - email channel",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "Jane Smith",
						Role:      "Director",
						Email:     "jane@example.com",
						SlackID:   "U789012", // Both contact methods present
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           1, // Email channel
					},
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{
					response: api.ApprovalResponse{
						ApproverName:      "Jane Smith",
						ApproverRole:      "Director",
						ApproverChannel:   "email",
						ApproverContactID: "jane@example.com",
					},
				},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantErr: false,
		},
		{
			name: "approver with both email and slack contact - slack channel",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "Bob Wilson",
						Role:      "Manager",
						Email:     "bob@example.com", // Both contact methods present
						SlackID:   "U789012",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           0, // Slack channel
					},
				},
				slack: &mockNotificationService{
					response: api.ApprovalResponse{
						ApproverName:      "Bob Wilson",
						ApproverRole:      "Manager",
						ApproverChannel:   "slack",
						ApproverContactID: "U789012",
					},
				},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantErr: false,
		},
		{
			name: "approver missing email",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "John Doe",
						Role:      "Manager",
						Email:     "", // No email
						SlackID:   "U123456",
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           0, // Slack channel
					},
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantErr:     true,
			expectedErr: api.ErrMissingEmail,
		},
		{
			name: "approver missing slack ID",
			service: &service{
				log:     &mockLogger{},
				company: Company{Name: "Test Company", Departments: []string{"Engineering"}},
				db: &mockDatabaseService{
					company: db.Company{ID: 1, Name: "Test Company"},
					approver: db.Approver{
						ID:        1,
						CompanyID: 1,
						Name:      "Jane Smith",
						Role:      "Director",
						Email:     "jane@example.com",
						SlackID:   "", // No Slack ID
					},
					rule: db.WorkflowRule{
						ID:                        1,
						CompanyID:                 1,
						MinAmount:                 floatPtr(100.0),
						MaxAmount:                 floatPtr(1000.0),
						Department:                stringPtr("Engineering"),
						IsManagerApprovalRequired: intPtr(0),
						ApproverID:                1,
						ApprovalChannel:           1, // Email channel
					},
				},
				slack: &mockNotificationService{},
				email: &mockNotificationService{},
			},
			invoice: api.InvoiceRequest{
				CompanyName:               "Test Company",
				Amount:                    500.0,
				Department:                "Engineering",
				IsManagerApprovalRequired: false,
			},
			wantErr:     true,
			expectedErr: api.ErrMissingSlackID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.service.ProcessInvoice(test.invoice)

			if test.wantErr {
				if err == nil {
					t.Errorf("ProcessInvoice() expected error but got none")
				}
				if test.expectedErr != nil && !errors.Is(err, test.expectedErr) {
					t.Errorf("ProcessInvoice() expected error %v but got %v", test.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("ProcessInvoice() unexpected error: %v", err)
				}
				if got.ApproverName == "" {
					t.Errorf("ProcessInvoice() expected valid response but got empty approver name")
				}
			}
		})
	}
}

type mockLogger struct{}

func (l *mockLogger) Info(msg string, args ...any) {
	// Mock implementation - just ignore for testing
}

func (l *mockLogger) Error(msg string, args ...any) {
	// Mock implementation - just ignore for testing
}

// Mock implementations for testing
type mockDatabaseService struct {
	company     db.Company
	companyErr  error
	approver    db.Approver
	approverErr error
	rule        db.WorkflowRule
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

// Helper functions for creating pointers
func floatPtr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}
