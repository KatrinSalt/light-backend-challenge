package db

// import (
// 	"errors"
// 	"testing"

// 	"github.com/KatrinSalt/backend-challenge-go/db/sql"
// 	"github.com/google/go-cmp/cmp"
// )

// func TestNewDBService(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		input struct {
// 			client  sql.Client
// 			options []ServiceOption
// 		}
// 		want    *service
// 		wantErr error
// 	}{
// 		{
// 			name: "new service with defaults",
// 			input: struct {
// 				client  sql.Client
// 				options []ServiceOption
// 			}{
// 				client:  &mockClient{},
// 				options: []ServiceOption{},
// 			},
// 			want: &service{
// 				client:            &mockClient{},
// 				CompanyStore:      nil, // Will be set by constructor
// 				ApproverStore:     nil, // Will be set by constructor
// 				WorkflowRuleStore: nil, // Will be set by constructor
// 				sampleData:        NewSampleData(),
// 			},
// 		},
// 		{
// 			name: "new service with custom table names",
// 			input: struct {
// 				client  sql.Client
// 				options []ServiceOption
// 			}{
// 				client: &mockClient{},
// 				options: []ServiceOption{
// 					WithCompanyTable("custom_companies"),
// 					WithApproverTable("custom_approvers"),
// 					WithWorkflowRuleTable("custom_workflow_rules"),
// 				},
// 			},
// 			want: &service{
// 				client:            &mockClient{},
// 				CompanyStore:      nil, // Will be set by constructor
// 				ApproverStore:     nil, // Will be set by constructor
// 				WorkflowRuleStore: nil, // Will be set by constructor
// 				sampleData:        NewSampleData(),
// 			},
// 		},
// 		{
// 			name: "new service with nil client",
// 			input: struct {
// 				client  sql.Client
// 				options []ServiceOption
// 			}{
// 				client:  nil,
// 				options: []ServiceOption{},
// 			},
// 			want:    nil,
// 			wantErr: errors.New("nil sql client"),
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			got, gotErr := NewService(test.input.client, test.input.options...)

// 			if test.wantErr != nil {
// 				if gotErr == nil {
// 					t.Errorf("NewDBService() = unexpected result, wanted error, got nil")
// 					return
// 				}
// 				if diff := cmp.Diff(test.wantErr.Error(), gotErr.Error()); diff != "" {
// 					t.Errorf("NewDBService() = unexpected error (-want +got)\n%s\n", diff)
// 				}
// 				return
// 			}

// 			if got == nil {
// 				t.Errorf("NewDBService() = unexpected result, got nil")
// 				return
// 			}

// 			// Check that the service was created with the client
// 			if got.client != test.input.client {
// 				t.Errorf("NewDBService() = unexpected client, want %v, got %v", test.input.client, got.client)
// 			}

// 			// Check that stores were created
// 			if got.CompanyStore == nil {
// 				t.Errorf("NewDBService() = CompanyStore is nil")
// 			}
// 			if got.ApproverStore == nil {
// 				t.Errorf("NewDBService() = ApproverStore is nil")
// 			}
// 			if got.WorkflowRuleStore == nil {
// 				t.Errorf("NewDBService() = WorkflowRuleStore is nil")
// 			}

// 			// Check that sample data was initialized
// 			if got.sampleData == nil {
// 				t.Errorf("NewDBService() = sampleData is nil")
// 			}
// 		})
// 	}
// }

// func TestDBService_Initialize(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		input struct {
// 			service *service
// 			queries []string
// 		}
// 		wantErr error
// 	}{
// 		{
// 			name: "initialize with valid queries",
// 			input: struct {
// 				service *service
// 				queries []string
// 			}{
// 				service: &service{
// 					client: &mockClient{},
// 				},
// 				queries: []string{
// 					"CREATE TABLE companies (id INTEGER PRIMARY KEY)",
// 					"CREATE TABLE approvers (id INTEGER PRIMARY KEY)",
// 				},
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "initialize with empty queries",
// 			input: struct {
// 				service *service
// 				queries []string
// 			}{
// 				service: &service{
// 					client: &mockClient{},
// 				},
// 				queries: []string{},
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "initialize with client error",
// 			input: struct {
// 				service *service
// 				queries []string
// 			}{
// 				service: &service{
// 					client: &mockClient{execErr: errors.New("database error")},
// 				},
// 				queries: []string{"CREATE TABLE test (id INTEGER)"},
// 			},
// 			wantErr: errors.New("failed to execute query 1: database error"),
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			gotErr := test.input.service.Initialize(test.input.queries)

// 			if test.wantErr != nil {
// 				if gotErr == nil {
// 					t.Errorf("Initialize() = unexpected result, wanted error, got nil")
// 					return
// 				}
// 				if diff := cmp.Diff(test.wantErr.Error(), gotErr.Error()); diff != "" {
// 					t.Errorf("Initialize() = unexpected error (-want +got)\n%s\n", diff)
// 				}
// 				return
// 			}

// 			if gotErr != nil {
// 				t.Errorf("Initialize() = unexpected error, got %v", gotErr)
// 			}
// 		})
// 	}
// }

// func TestDBService_SeedSampleData(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		input struct {
// 			service *service
// 		}
// 		wantErr error
// 	}{
// 		{
// 			name: "seed sample data successfully",
// 			input: struct {
// 				service *service
// 			}{
// 				service: &service{
// 					CompanyStore:      &mockCompanyStore{},
// 					ApproverStore:     &mockApproverStore{},
// 					WorkflowRuleStore: &mockWorkflowRuleStore{},
// 					sampleData:        NewSampleData(),
// 				},
// 			},
// 			wantErr: nil,
// 		},
// 		{
// 			name: "seed sample data with nil sample data",
// 			input: struct {
// 				service *service
// 			}{
// 				service: &service{
// 					CompanyStore:      &mockCompanyStore{},
// 					ApproverStore:     &mockApproverStore{},
// 					WorkflowRuleStore: &mockWorkflowRuleStore{},
// 					sampleData:        nil,
// 				},
// 			},
// 			wantErr: errors.New("no sample data provided"),
// 		},
// 		{
// 			name: "seed sample data with company store error",
// 			input: struct {
// 				service *service
// 			}{
// 				service: &service{
// 					CompanyStore: &mockCompanyStore{
// 						createErr: errors.New("company creation failed"),
// 					},
// 					ApproverStore:     &mockApproverStore{},
// 					WorkflowRuleStore: &mockWorkflowRuleStore{},
// 					sampleData:        NewSampleData(),
// 				},
// 			},
// 			wantErr: errors.New("company creation failed"),
// 		},
// 		{
// 			name: "seed sample data with approver store error",
// 			input: struct {
// 				service *service
// 			}{
// 				service: &service{
// 					CompanyStore: &mockCompanyStore{},
// 					ApproverStore: &mockApproverStore{
// 						createErr: errors.New("approver creation failed"),
// 					},
// 					WorkflowRuleStore: &mockWorkflowRuleStore{},
// 					sampleData:        NewSampleData(),
// 				},
// 			},
// 			wantErr: errors.New("approver creation failed"),
// 		},
// 		{
// 			name: "seed sample data with workflow rule store error",
// 			input: struct {
// 				service *service
// 			}{
// 				service: &service{
// 					CompanyStore:  &mockCompanyStore{},
// 					ApproverStore: &mockApproverStore{},
// 					WorkflowRuleStore: &mockWorkflowRuleStore{
// 						createErr: errors.New("workflow rule creation failed"),
// 					},
// 					sampleData: NewSampleData(),
// 				},
// 			},
// 			wantErr: errors.New("workflow rule creation failed"),
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			gotErr := test.input.service.SeedSampleData()

// 			if test.wantErr != nil {
// 				if gotErr == nil {
// 					t.Errorf("SeedSampleData() = unexpected result, wanted error, got nil")
// 					return
// 				}
// 				if diff := cmp.Diff(test.wantErr.Error(), gotErr.Error()); diff != "" {
// 					t.Errorf("SeedSampleData() = unexpected error (-want +got)\n%s\n", diff)
// 				}
// 				return
// 			}

// 			if gotErr != nil {
// 				t.Errorf("SeedSampleData() = unexpected error, got %v", gotErr)
// 			}
// 		})
// 	}
// }

// // Mock implementations for testing
// type mockClient struct {
// 	execErr error
// }

// func (m *mockClient) QueryRow(query string, args ...any) sql.Row {
// 	return nil
// }

// func (m *mockClient) Query(query string, args ...any) (sql.Rows, error) {
// 	return nil, nil
// }

// func (m *mockClient) Transaction() (sql.Tx, error) {
// 	return nil, nil
// }

// func (m *mockClient) Begin() (sql.Tx, error) {
// 	return nil, nil
// }

// func (m *mockClient) Close() error {
// 	return nil
// }

// func (m *mockClient) Exec(query string, args ...any) (sql.Result, error) {
// 	if m.execErr != nil {
// 		return nil, m.execErr
// 	}
// 	return nil, nil
// }

// type mockCompanyStore struct {
// 	createErr error
// 	company   Company
// }

// func (m *mockCompanyStore) Create(company Company) (Company, error) {
// 	if m.createErr != nil {
// 		return Company{}, m.createErr
// 	}
// 	return m.company, nil
// }

// func (m *mockCompanyStore) GetByID(id int) (Company, error) {
// 	return Company{}, nil
// }

// func (m *mockCompanyStore) GetByName(name string) (Company, error) {
// 	return Company{}, nil
// }

// type mockApproverStore struct {
// 	createErr error
// 	approver  Approver
// }

// func (m *mockApproverStore) Create(approver Approver) (Approver, error) {
// 	if m.createErr != nil {
// 		return Approver{}, m.createErr
// 	}
// 	return m.approver, nil
// }

// func (m *mockApproverStore) GetByID(id int) (Approver, error) {
// 	return Approver{}, nil
// }

// type mockWorkflowRuleStore struct {
// 	createErr error
// 	rule      WorkflowRule
// }

// func (m *mockWorkflowRuleStore) Create(rule WorkflowRule) (WorkflowRule, error) {
// 	if m.createErr != nil {
// 		return WorkflowRule{}, m.createErr
// 	}
// 	return m.rule, nil
// }

// func (m *mockWorkflowRuleStore) GetByCompanyID(companyID int) ([]WorkflowRule, error) {
// 	return nil, nil
// }

// func (m *mockWorkflowRuleStore) FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (WorkflowRule, error) {
// 	return m.rule, nil
// }
