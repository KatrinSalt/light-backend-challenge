package management

import (
	"errors"
	"strings"
	"testing"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/common"
	"github.com/KatrinSalt/backend-challenge-go/db"
	"github.com/google/go-cmp/cmp"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			logger      common.Logger
			dbService   databaseService
			companyName string
		}
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful service creation",
			input: struct {
				logger      common.Logger
				dbService   databaseService
				companyName string
			}{
				logger: &mockLogger{},
				dbService: &mockDBService{
					getCompanyByNameResult: db.Company{
						ID:   1,
						Name: "Test Company",
					},
				},
				companyName: "Test Company",
			},
			wantErr: false,
		},
		{
			name: "empty company name",
			input: struct {
				logger      common.Logger
				dbService   databaseService
				companyName string
			}{
				logger:      &mockLogger{},
				dbService:   &mockDBService{},
				companyName: "",
			},
			wantErr: true,
			errMsg:  "company name is required",
		},
		{
			name: "company not found",
			input: struct {
				logger      common.Logger
				dbService   databaseService
				companyName string
			}{
				logger: &mockLogger{},
				dbService: &mockDBService{
					getCompanyByNameErr: errors.New("company not found"),
				},
				companyName: "Non-existent Company",
			},
			wantErr: true,
			errMsg:  "failed to get company",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewService(test.input.logger, test.input.dbService, test.input.companyName)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("NewService() expected error but got none")
					return
				}
				if test.errMsg != "" && !strings.Contains(gotErr.Error(), test.errMsg) {
					t.Errorf("NewService() expected error containing %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("NewService() unexpected error: %v", gotErr)
				return
			}

			if got == nil {
				t.Errorf("NewService() expected service but got nil")
			}
		})
	}
}

func TestService_CreateWorkflowRule(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			service *service
			rule    api.WorkflowRule
		}
		want    api.WorkflowRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful workflow rule creation",
			input: struct {
				service *service
				rule    api.WorkflowRule
			}{
				service: &service{
					logger: &mockLogger{},
					dbService: &mockDBService{
						createWorkflowRuleResult: db.WorkflowRule{
							ID:                        1,
							CompanyID:                 1,
							MinAmount:                 floatPtr(100.0),
							MaxAmount:                 floatPtr(500.0),
							Department:                stringPtr("Finance"),
							IsManagerApprovalRequired: intPtr(1),
							ApproverID:                1,
							ApprovalChannel:           0,
						},
					},
					company: company{
						id:   1,
						name: "Test Company",
					},
				},
				rule: api.WorkflowRule{
					CompanyID:                 1,
					MinAmount:                 floatPtr(100.0),
					MaxAmount:                 floatPtr(500.0),
					Department:                stringPtr("Finance"),
					IsManagerApprovalRequired: 1,
					ApproverID:                1,
					ApprovalChannel:           0,
				},
			},
			want: api.WorkflowRule{
				ID:                        1,
				CompanyID:                 1,
				MinAmount:                 floatPtr(100.0),
				MaxAmount:                 floatPtr(500.0),
				Department:                stringPtr("Finance"),
				IsManagerApprovalRequired: 1,
				ApproverID:                1,
				ApprovalChannel:           0,
			},
			wantErr: false,
		},
		{
			name: "invalid workflow rule - invalid approval channel",
			input: struct {
				service *service
				rule    api.WorkflowRule
			}{
				service: &service{
					logger:    &mockLogger{},
					dbService: &mockDBService{},
					company: company{
						id:   1,
						name: "Test Company",
					},
				},
				rule: api.WorkflowRule{
					CompanyID:       1,
					ApproverID:      1,
					ApprovalChannel: 2, // Invalid channel
				},
			},
			want:    api.WorkflowRule{},
			wantErr: true,
			errMsg:  "invalid workflow rule",
		},
		{
			name: "invalid workflow rule - invalid amount range",
			input: struct {
				service *service
				rule    api.WorkflowRule
			}{
				service: &service{
					logger:    &mockLogger{},
					dbService: &mockDBService{},
					company: company{
						id:   1,
						name: "Test Company",
					},
				},
				rule: api.WorkflowRule{
					CompanyID:       1,
					MinAmount:       floatPtr(500.0),
					MaxAmount:       floatPtr(100.0), // Min > Max
					ApproverID:      1,
					ApprovalChannel: 0,
				},
			},
			want:    api.WorkflowRule{},
			wantErr: true,
			errMsg:  "invalid workflow rule",
		},
		{
			name: "database error during creation",
			input: struct {
				service *service
				rule    api.WorkflowRule
			}{
				service: &service{
					logger: &mockLogger{},
					dbService: &mockDBService{
						createWorkflowRuleErr: errors.New("database error"),
					},
					company: company{
						id:   1,
						name: "Test Company",
					},
				},
				rule: api.WorkflowRule{
					CompanyID:       1,
					ApproverID:      1,
					ApprovalChannel: 0,
				},
			},
			want:    api.WorkflowRule{},
			wantErr: true,
			errMsg:  "failed to create workflow rule",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.service.CreateWorkflowRule(test.input.rule)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("CreateWorkflowRule() expected error but got none")
					return
				}
				if test.errMsg != "" && !strings.Contains(gotErr.Error(), test.errMsg) {
					t.Errorf("CreateWorkflowRule() expected error containing %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("CreateWorkflowRule() unexpected error: %v", gotErr)
				return
			}

			if !cmp.Equal(got, test.want) {
				t.Errorf("CreateWorkflowRule() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestService_CreateApprover(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			service  *service
			approver api.Approver
		}
		want    api.Approver
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful approver creation",
			input: struct {
				service  *service
				approver api.Approver
			}{
				service: &service{
					logger: &mockLogger{},
					dbService: &mockDBService{
						createApproverResult: db.Approver{
							ID:        1,
							CompanyID: 1,
							Name:      "John Doe",
							Role:      "Manager",
							Email:     "john@example.com",
							SlackID:   "U123456",
						},
					},
					company: company{
						id:   1,
						name: "Test Company",
					},
				},
				approver: api.Approver{
					Name:    "John Doe",
					Role:    "Manager",
					Email:   "john@example.com",
					SlackID: "U123456",
				},
			},
			want: api.Approver{
				ID:      1,
				Name:    "John Doe",
				Role:    "Manager",
				Email:   "john@example.com",
				SlackID: "U123456",
			},
			wantErr: false,
		},
		{
			name: "invalid approver - missing email",
			input: struct {
				service  *service
				approver api.Approver
			}{
				service: &service{
					logger:    &mockLogger{},
					dbService: &mockDBService{},
					company: company{
						id:   1,
						name: "Test Company",
					},
				},
				approver: api.Approver{
					Name:    "John Doe",
					Role:    "Manager",
					SlackID: "U123456",
					// Missing Email
				},
			},
			want:    api.Approver{},
			wantErr: true,
			errMsg:  "invalid approver",
		},
		{
			name: "database error during creation",
			input: struct {
				service  *service
				approver api.Approver
			}{
				service: &service{
					logger: &mockLogger{},
					dbService: &mockDBService{
						createApproverErr: errors.New("database error"),
					},
					company: company{
						id:   1,
						name: "Test Company",
					},
				},
				approver: api.Approver{
					Name:    "John Doe",
					Role:    "Manager",
					Email:   "john@example.com",
					SlackID: "U123456",
				},
			},
			want:    api.Approver{},
			wantErr: true,
			errMsg:  "failed to create approver",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.service.CreateApprover(test.input.approver)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("CreateApprover() expected error but got none")
					return
				}
				if test.errMsg != "" && !strings.Contains(gotErr.Error(), test.errMsg) {
					t.Errorf("CreateApprover() expected error containing %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("CreateApprover() unexpected error: %v", gotErr)
				return
			}

			if !cmp.Equal(got, test.want) {
				t.Errorf("CreateApprover() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestService_ListWorkflowRules(t *testing.T) {
	tests := []struct {
		name    string
		input   *service
		want    []api.WorkflowRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful workflow rules listing",
			input: &service{
				logger: &mockLogger{},
				dbService: &mockDBService{
					listWorkflowRulesResult: []db.WorkflowRule{
						{
							ID:                        1,
							CompanyID:                 1,
							MinAmount:                 floatPtr(100.0),
							MaxAmount:                 floatPtr(500.0),
							Department:                stringPtr("Finance"),
							IsManagerApprovalRequired: intPtr(1),
							ApproverID:                1,
							ApprovalChannel:           0,
						},
					},
				},
				company: company{
					id:   1,
					name: "Test Company",
				},
			},
			want: []api.WorkflowRule{
				{
					ID:                        1,
					CompanyID:                 1,
					MinAmount:                 floatPtr(100.0),
					MaxAmount:                 floatPtr(500.0),
					Department:                stringPtr("Finance"),
					IsManagerApprovalRequired: 1,
					ApproverID:                1,
					ApprovalChannel:           0,
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.ListWorkflowRules()

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("ListWorkflowRules() expected error but got none")
					return
				}
				if test.errMsg != "" && !strings.Contains(gotErr.Error(), test.errMsg) {
					t.Errorf("ListWorkflowRules() expected error containing %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("ListWorkflowRules() unexpected error: %v", gotErr)
				return
			}

			if len(got) != len(test.want) {
				t.Errorf("ListWorkflowRules() returned %d rules, want %d", len(got), len(test.want))
				return
			}

			if len(got) > 0 && got[0].ID != test.want[0].ID {
				t.Errorf("ListWorkflowRules() first rule ID = %v, want %v", got[0].ID, test.want[0].ID)
			}
		})
	}
}

func TestService_ListApprovers(t *testing.T) {
	tests := []struct {
		name    string
		input   *service
		want    []api.Approver
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful approvers listing",
			input: &service{
				logger: &mockLogger{},
				dbService: &mockDBService{
					listApproversResult: []db.Approver{
						{
							ID:        1,
							CompanyID: 1,
							Name:      "John Doe",
							Role:      "Manager",
							Email:     "john@example.com",
							SlackID:   "U123456",
						},
					},
				},
				company: company{
					id:   1,
					name: "Test Company",
				},
			},
			want: []api.Approver{
				{
					ID:      1,
					Name:    "John Doe",
					Role:    "Manager",
					Email:   "john@example.com",
					SlackID: "U123456",
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.ListApprovers()

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("ListApprovers() expected error but got none")
					return
				}
				if test.errMsg != "" && !strings.Contains(gotErr.Error(), test.errMsg) {
					t.Errorf("ListApprovers() expected error containing %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("ListApprovers() unexpected error: %v", gotErr)
				return
			}

			if len(got) != len(test.want) {
				t.Errorf("ListApprovers() returned %d approvers, want %d", len(got), len(test.want))
				return
			}

			if len(got) > 0 && got[0].ID != test.want[0].ID {
				t.Errorf("ListApprovers() first approver ID = %v, want %v", got[0].ID, test.want[0].ID)
			}
		})
	}
}

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

// Mock implementations for testing
type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...interface{})  {}
func (m *mockLogger) Error(msg string, args ...interface{}) {}

type mockDBService struct {
	// Company methods
	getCompanyByNameResult db.Company
	getCompanyByNameErr    error

	// Workflow Rule methods
	createWorkflowRuleResult db.WorkflowRule
	createWorkflowRuleErr    error
	getWorkflowRuleResult    db.WorkflowRule
	getWorkflowRuleErr       error
	updateWorkflowRuleErr    error
	deleteWorkflowRuleErr    error
	listWorkflowRulesResult  []db.WorkflowRule
	listWorkflowRulesErr     error

	// Approver methods
	createApproverResult db.Approver
	createApproverErr    error
	getApproverResult    db.Approver
	getApproverErr       error
	updateApproverErr    error
	deleteApproverErr    error
	listApproversResult  []db.Approver
	listApproversErr     error
}

// mockDBService implements management.databaseService interface

func (m *mockDBService) GetSampleData() *db.SampleData {
	return nil
}

func (m *mockDBService) Initialize() error {
	return nil
}

func (m *mockDBService) SeedSampleData() error {
	return nil
}

func (m *mockDBService) GetCompanyByName(name string) (db.Company, error) {
	return m.getCompanyByNameResult, m.getCompanyByNameErr
}

func (m *mockDBService) FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (db.WorkflowRule, error) {
	return db.WorkflowRule{}, nil
}

func (m *mockDBService) CreateWorkflowRule(rule db.WorkflowRule) (db.WorkflowRule, error) {
	if m.createWorkflowRuleErr != nil {
		return db.WorkflowRule{}, m.createWorkflowRuleErr
	}
	return m.createWorkflowRuleResult, nil
}

func (m *mockDBService) GetWorkflowRuleByID(id int) (db.WorkflowRule, error) {
	if m.getWorkflowRuleErr != nil {
		return db.WorkflowRule{}, m.getWorkflowRuleErr
	}
	return m.getWorkflowRuleResult, nil
}

func (m *mockDBService) UpdateWorkflowRule(rule db.WorkflowRule) error {
	return m.updateWorkflowRuleErr
}

func (m *mockDBService) DeleteWorkflowRule(id int) error {
	return m.deleteWorkflowRuleErr
}

func (m *mockDBService) ListWorkflowRules(companyID int) ([]db.WorkflowRule, error) {
	if m.listWorkflowRulesErr != nil {
		return nil, m.listWorkflowRulesErr
	}
	return m.listWorkflowRulesResult, nil
}

func (m *mockDBService) CreateApprover(approver db.Approver) (db.Approver, error) {
	if m.createApproverErr != nil {
		return db.Approver{}, m.createApproverErr
	}
	return m.createApproverResult, nil
}

func (m *mockDBService) GetApproverByID(id int) (db.Approver, error) {
	if m.getApproverErr != nil {
		return db.Approver{}, m.getApproverErr
	}
	return m.getApproverResult, nil
}

func (m *mockDBService) UpdateApprover(approver db.Approver) error {
	return m.updateApproverErr
}

func (m *mockDBService) DeleteApprover(id int) error {
	return m.deleteApproverErr
}

func (m *mockDBService) ListApprovers(companyID int) ([]db.Approver, error) {
	if m.listApproversErr != nil {
		return nil, m.listApproversErr
	}
	return m.listApproversResult, nil
}
