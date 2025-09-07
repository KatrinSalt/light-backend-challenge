package db

import (
	"errors"
	"testing"

	"github.com/KatrinSalt/backend-challenge-go/db/sql"
	"github.com/google/go-cmp/cmp"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			client  sql.Client
			options []ServiceOption
		}
		want    *service
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid service with defaults",
			input: struct {
				client  sql.Client
				options []ServiceOption
			}{
				client:  &mockClient{},
				options: []ServiceOption{},
			},
			want:    &service{client: &mockClient{}},
			wantErr: false,
		},
		{
			name: "valid service with custom table names",
			input: struct {
				client  sql.Client
				options []ServiceOption
			}{
				client: &mockClient{},
				options: []ServiceOption{
					WithCompanyTable("custom_companies"),
					WithApproverTable("custom_approvers"),
					WithWorkflowRuleTable("custom_workflow_rules"),
				},
			},
			want:    &service{client: &mockClient{}},
			wantErr: false,
		},
		{
			name: "valid service with custom sample data",
			input: struct {
				client  sql.Client
				options []ServiceOption
			}{
				client: &mockClient{},
				options: []ServiceOption{
					WithSampleData(&SampleData{
						Companies: []Company{{Name: "Test Company"}},
					}),
				},
			},
			want:    &service{client: &mockClient{}},
			wantErr: false,
		},
		{
			name: "valid service with custom schema",
			input: struct {
				client  sql.Client
				options []ServiceOption
			}{
				client: &mockClient{},
				options: []ServiceOption{
					WithSchema([]string{"CREATE TABLE test (id INTEGER)"}),
				},
			},
			want:    &service{client: &mockClient{}},
			wantErr: false,
		},
		{
			name: "valid service with multiple options",
			input: struct {
				client  sql.Client
				options []ServiceOption
			}{
				client: &mockClient{},
				options: []ServiceOption{
					WithCompanyTable("custom_companies"),
					WithSampleData(&SampleData{Companies: []Company{{Name: "Test"}}}),
					WithSchema([]string{"CREATE TABLE test (id INTEGER)"}),
				},
			},
			want:    &service{client: &mockClient{}},
			wantErr: false,
		},
		{
			name: "nil client",
			input: struct {
				client  sql.Client
				options []ServiceOption
			}{
				client:  nil,
				options: []ServiceOption{},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "nil sql client",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewService(test.input.client, test.input.options...)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("NewService() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("NewService() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("NewService() unexpected error: %v", gotErr)
				return
			}

			if got == nil {
				t.Errorf("NewService() returned nil service")
				return
			}

			// Check that the service was created with the client
			if got.client != test.input.client {
				t.Errorf("NewService() client mismatch, want %v, got %v", test.input.client, got.client)
			}

			// Check that stores were created
			if got.companyStore == nil {
				t.Errorf("NewService() companyStore is nil")
			}
			if got.approverStore == nil {
				t.Errorf("NewService() approverStore is nil")
			}
			if got.workflowRuleStore == nil {
				t.Errorf("NewService() workflowRuleStore is nil")
			}

			// Check that sample data was initialized
			if got.sampleData == nil {
				t.Errorf("NewService() sampleData is nil")
			}
		})
	}
}

func TestService_Initialize(t *testing.T) {
	tests := []struct {
		name    string
		input   *service
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful initialization with valid queries",
			input: &service{
				client: &mockClient{},
				schema: []string{
					"CREATE TABLE companies (id INTEGER PRIMARY KEY)",
					"CREATE TABLE approvers (id INTEGER PRIMARY KEY)",
				},
			},
			wantErr: false,
		},
		{
			name: "successful initialization with empty schema",
			input: &service{
				client: &mockClient{},
				schema: []string{},
			},
			wantErr: false,
		},
		{
			name: "initialization fails on first query",
			input: &service{
				client: &mockClient{execErr: errors.New("database error")},
				schema: []string{"CREATE TABLE test (id INTEGER)"},
			},
			wantErr: true,
			errMsg:  "failed to execute query 1: database error",
		},
		{
			name: "initialization fails on second query",
			input: &service{
				client: &mockClient{execErr: errors.New("table exists"), failOnQuery: 2},
				schema: []string{
					"CREATE TABLE companies (id INTEGER PRIMARY KEY)",
					"CREATE TABLE approvers (id INTEGER PRIMARY KEY)",
				},
			},
			wantErr: true,
			errMsg:  "failed to execute query 2: table exists",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := test.input.Initialize()

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("Initialize() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("Initialize() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("Initialize() unexpected error: %v", gotErr)
			}
		})
	}
}

func TestService_SeedSampleData(t *testing.T) {
	tests := []struct {
		name    string
		input   *service
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful seeding",
			input: &service{
				companyStore:      &mockCompanyStore{},
				approverStore:     &mockApproverStore{},
				workflowRuleStore: &mockWorkflowRuleStore{},
				sampleData:        NewSampleData(),
			},
			wantErr: false,
		},
		{
			name: "nil sample data",
			input: &service{
				companyStore:      &mockCompanyStore{},
				approverStore:     &mockApproverStore{},
				workflowRuleStore: &mockWorkflowRuleStore{},
				sampleData:        nil,
			},
			wantErr: true,
			errMsg:  "no sample data provided",
		},
		{
			name: "company store creation error",
			input: &service{
				companyStore: &mockCompanyStore{
					createErr: errors.New("company creation failed"),
				},
				approverStore:     &mockApproverStore{},
				workflowRuleStore: &mockWorkflowRuleStore{},
				sampleData:        NewSampleData(),
			},
			wantErr: true,
			errMsg:  "company creation failed",
		},
		{
			name: "approver store creation error",
			input: &service{
				companyStore: &mockCompanyStore{},
				approverStore: &mockApproverStore{
					createErr: errors.New("approver creation failed"),
				},
				workflowRuleStore: &mockWorkflowRuleStore{},
				sampleData:        NewSampleData(),
			},
			wantErr: true,
			errMsg:  "approver creation failed",
		},
		{
			name: "workflow rule store creation error",
			input: &service{
				companyStore:  &mockCompanyStore{},
				approverStore: &mockApproverStore{},
				workflowRuleStore: &mockWorkflowRuleStore{
					createErr: errors.New("workflow rule creation failed"),
				},
				sampleData: NewSampleData(),
			},
			wantErr: true,
			errMsg:  "workflow rule creation failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := test.input.SeedSampleData()

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("SeedSampleData() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("SeedSampleData() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("SeedSampleData() unexpected error: %v", gotErr)
			}
		})
	}
}

func TestService_GetSampleData(t *testing.T) {
	tests := []struct {
		name  string
		input *service
		want  *SampleData
	}{
		{
			name: "returns sample data",
			input: &service{
				sampleData: NewSampleData(),
			},
			want: NewSampleData(),
		},
		{
			name: "returns nil when no sample data",
			input: &service{
				sampleData: nil,
			},
			want: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.GetSampleData()

			if test.want == nil {
				if got != nil {
					t.Errorf("GetSampleData() expected nil but got %v", got)
				}
				return
			}

			if got == nil {
				t.Errorf("GetSampleData() expected sample data but got nil")
				return
			}

			// Compare the sample data structure
			if len(got.Companies) != len(test.want.Companies) {
				t.Errorf("GetSampleData() companies count mismatch, want %d, got %d",
					len(test.want.Companies), len(got.Companies))
			}
			if len(got.Approvers) != len(test.want.Approvers) {
				t.Errorf("GetSampleData() approvers count mismatch, want %d, got %d",
					len(test.want.Approvers), len(got.Approvers))
			}
			if len(got.WorkflowRules) != len(test.want.WorkflowRules) {
				t.Errorf("GetSampleData() workflow rules count mismatch, want %d, got %d",
					len(test.want.WorkflowRules), len(got.WorkflowRules))
			}
		})
	}
}

func TestService_GetCompanyByName(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			service     *service
			companyName string
		}
		want    Company
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful company retrieval",
			input: struct {
				service     *service
				companyName string
			}{
				service: &service{
					companyStore: &mockCompanyStore{
						company: Company{ID: 1, Name: "Light"},
					},
				},
				companyName: "Light",
			},
			want:    Company{ID: 1, Name: "Light"},
			wantErr: false,
		},
		{
			name: "company not found",
			input: struct {
				service     *service
				companyName string
			}{
				service: &service{
					companyStore: &mockCompanyStore{
						getByNameErr: errors.New("company not found"),
					},
				},
				companyName: "NonExistent",
			},
			want:    Company{},
			wantErr: true,
			errMsg:  "company not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.service.GetCompanyByName(test.input.companyName)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("GetCompanyByName() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("GetCompanyByName() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("GetCompanyByName() unexpected error: %v", gotErr)
				return
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetCompanyByName() mismatch (-want +got)\n%s", diff)
			}
		})
	}
}

func TestService_GetApproverByID(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			service    *service
			approverID int
		}
		want    Approver
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful approver retrieval",
			input: struct {
				service    *service
				approverID int
			}{
				service: &service{
					approverStore: &mockApproverStore{
						approver: Approver{
							ID:        1,
							CompanyID: 1,
							Name:      "John Doe",
							Role:      "Manager",
							Email:     "john@example.com",
							SlackID:   "U123456",
						},
					},
				},
				approverID: 1,
			},
			want: Approver{
				ID:        1,
				CompanyID: 1,
				Name:      "John Doe",
				Role:      "Manager",
				Email:     "john@example.com",
				SlackID:   "U123456",
			},
			wantErr: false,
		},
		{
			name: "approver not found",
			input: struct {
				service    *service
				approverID int
			}{
				service: &service{
					approverStore: &mockApproverStore{
						getByIDErr: errors.New("approver not found"),
					},
				},
				approverID: 999,
			},
			want:    Approver{},
			wantErr: true,
			errMsg:  "approver not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.service.GetApproverByID(test.input.approverID)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("GetApproverByID() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("GetApproverByID() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("GetApproverByID() unexpected error: %v", gotErr)
				return
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetApproverByID() mismatch (-want +got)\n%s", diff)
			}
		})
	}
}

func TestService_FindMatchingRule(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			service         *service
			companyID       int
			amount          float64
			department      string
			requiresManager bool
		}
		want    WorkflowRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful rule matching",
			input: struct {
				service         *service
				companyID       int
				amount          float64
				department      string
				requiresManager bool
			}{
				service: &service{
					workflowRuleStore: &mockWorkflowRuleStore{
						rule: WorkflowRule{
							ID:                        1,
							CompanyID:                 1,
							MinAmount:                 floatPtr(1000.0),
							MaxAmount:                 floatPtr(5000.0),
							Department:                nil,
							IsManagerApprovalRequired: intPtr(0),
							ApproverID:                1,
							ApprovalChannel:           0,
						},
					},
				},
				companyID:       1,
				amount:          2500.0,
				department:      "Finance",
				requiresManager: false,
			},
			want: WorkflowRule{
				ID:                        1,
				CompanyID:                 1,
				MinAmount:                 floatPtr(1000.0),
				MaxAmount:                 floatPtr(5000.0),
				Department:                nil,
				IsManagerApprovalRequired: intPtr(0),
				ApproverID:                1,
				ApprovalChannel:           0,
			},
			wantErr: false,
		},
		{
			name: "no matching rule found",
			input: struct {
				service         *service
				companyID       int
				amount          float64
				department      string
				requiresManager bool
			}{
				service: &service{
					workflowRuleStore: &mockWorkflowRuleStore{
						findMatchingRuleErr: errors.New("no matching rule found"),
					},
				},
				companyID:       1,
				amount:          10000.0,
				department:      "Unknown",
				requiresManager: false,
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "no matching rule found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.service.FindMatchingRule(
				test.input.companyID,
				test.input.amount,
				test.input.department,
				test.input.requiresManager,
			)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("FindMatchingRule() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("FindMatchingRule() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("FindMatchingRule() unexpected error: %v", gotErr)
				return
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("FindMatchingRule() mismatch (-want +got)\n%s", diff)
			}
		})
	}
}

// Helper functions for creating pointers
func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

// Mock implementations for testing
type mockClient struct {
	execErr     error
	execCount   int
	failOnQuery int // Which query number should fail (1-based)
}

func (m *mockClient) QueryRow(query string, args ...any) sql.Row {
	return nil
}

func (m *mockClient) Query(query string, args ...any) (sql.Rows, error) {
	return nil, nil
}

func (m *mockClient) Transaction() (sql.Tx, error) {
	return nil, nil
}

func (m *mockClient) Begin() (sql.Tx, error) {
	return nil, nil
}

func (m *mockClient) Close() error {
	return nil
}

func (m *mockClient) Exec(query string, args ...any) (sql.Result, error) {
	m.execCount++
	if m.execErr != nil && (m.failOnQuery == 0 || m.execCount == m.failOnQuery) {
		return nil, m.execErr
	}
	return nil, nil
}

type mockCompanyStore struct {
	createErr    error
	getByNameErr error
	company      Company
}

func (m *mockCompanyStore) Create(company Company) (Company, error) {
	if m.createErr != nil {
		return Company{}, m.createErr
	}
	return m.company, nil
}

func (m *mockCompanyStore) GetByID(id int) (Company, error) {
	return Company{}, nil
}

func (m *mockCompanyStore) GetByName(name string) (Company, error) {
	if m.getByNameErr != nil {
		return Company{}, m.getByNameErr
	}
	return m.company, nil
}

type mockApproverStore struct {
	createErr  error
	getByIDErr error
	approver   Approver
}

func (m *mockApproverStore) Create(approver Approver) (Approver, error) {
	if m.createErr != nil {
		return Approver{}, m.createErr
	}
	return m.approver, nil
}

func (m *mockApproverStore) GetByID(id int) (Approver, error) {
	if m.getByIDErr != nil {
		return Approver{}, m.getByIDErr
	}
	return m.approver, nil
}

type mockWorkflowRuleStore struct {
	createErr           error
	findMatchingRuleErr error
	rule                WorkflowRule
}

func (m *mockWorkflowRuleStore) Create(rule WorkflowRule) (WorkflowRule, error) {
	if m.createErr != nil {
		return WorkflowRule{}, m.createErr
	}
	return m.rule, nil
}

func (m *mockWorkflowRuleStore) GetByCompanyID(companyID int) ([]WorkflowRule, error) {
	return nil, nil
}

func (m *mockWorkflowRuleStore) FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (WorkflowRule, error) {
	if m.findMatchingRuleErr != nil {
		return WorkflowRule{}, m.findMatchingRuleErr
	}
	return m.rule, nil
}
