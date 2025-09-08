package db

import (
	"errors"
	"strings"
	"testing"

	sqlpkg "github.com/KatrinSalt/backend-challenge-go/db/sql"
	"github.com/google/go-cmp/cmp"
)

func TestNewWorkflowRuleStore(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			client  sqlpkg.Client
			options []WorkflowRuleStoreOption
		}
		want    *workflowRuleStore
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid store with defaults",
			input: struct {
				client  sqlpkg.Client
				options []WorkflowRuleStoreOption
			}{
				client:  &mockSQLClient{},
				options: []WorkflowRuleStoreOption{},
			},
			want: &workflowRuleStore{
				client: &mockSQLClient{},
				table:  "workflow_rules",
			},
			wantErr: false,
		},
		{
			name: "valid store with custom table name",
			input: struct {
				client  sqlpkg.Client
				options []WorkflowRuleStoreOption
			}{
				client: &mockSQLClient{},
				options: []WorkflowRuleStoreOption{
					func(o *WorkflowRuleStoreOptions) {
						o.Table = "custom_workflow_rules"
					},
				},
			},
			want: &workflowRuleStore{
				client: &mockSQLClient{},
				table:  "custom_workflow_rules",
			},
			wantErr: false,
		},
		{
			name: "nil client",
			input: struct {
				client  sqlpkg.Client
				options []WorkflowRuleStoreOption
			}{
				client:  nil,
				options: []WorkflowRuleStoreOption{},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "nil sql client",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewWorkflowRuleStore(test.input.client, test.input.options...)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("NewWorkflowRuleStore() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("NewWorkflowRuleStore() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("NewWorkflowRuleStore() unexpected error: %v", gotErr)
				return
			}

			if got == nil {
				t.Errorf("NewWorkflowRuleStore() returned nil store")
				return
			}

			// Check that the service was created with a client
			if got.client == nil {
				t.Errorf("NewWorkflowRuleStore() client should not be nil")
			}
			if got.table != test.want.table {
				t.Errorf("NewWorkflowRuleStore() table mismatch, want %q, got %q", test.want.table, got.table)
			}
		})
	}
}

func TestWorkflowRuleStore_Create(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store        *workflowRuleStore
			workflowRule WorkflowRule
		}
		want    WorkflowRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful workflow rule creation",
			input: struct {
				store        *workflowRuleStore
				workflowRule WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								values: []interface{}{1, 1, floatPtr(1000.0), floatPtr(5000.0), stringPtr("Finance"), intPtr(0), 1, 0},
							},
						},
					},
					table: "workflow_rules",
				},
				workflowRule: WorkflowRule{
					CompanyID:                 1,
					MinAmount:                 floatPtr(1000.0),
					MaxAmount:                 floatPtr(5000.0),
					Department:                stringPtr("Finance"),
					IsManagerApprovalRequired: intPtr(0),
					ApproverID:                1,
					ApprovalChannel:           0,
				},
			},
			want: WorkflowRule{
				ID:                        1,
				CompanyID:                 1,
				MinAmount:                 floatPtr(1000.0),
				MaxAmount:                 floatPtr(5000.0),
				Department:                stringPtr("Finance"),
				IsManagerApprovalRequired: intPtr(0),
				ApproverID:                1,
				ApprovalChannel:           0,
			},
			wantErr: false,
		},
		{
			name: "workflow rule already exists",
			input: struct {
				store        *workflowRuleStore
				workflowRule WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execErr: errors.New("duplicate key value violates unique constraint"),
						},
					},
					table: "workflow_rules",
				},
				workflowRule: WorkflowRule{
					CompanyID:                 1,
					MinAmount:                 floatPtr(1000.0),
					MaxAmount:                 floatPtr(5000.0),
					Department:                stringPtr("Finance"),
					IsManagerApprovalRequired: intPtr(0),
					ApproverID:                1,
					ApprovalChannel:           0,
				},
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "duplicate key value violates unique constraint",
		},
		{
			name: "transaction creation fails",
			input: struct {
				store        *workflowRuleStore
				workflowRule WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						txErr: errors.New("database connection failed"),
					},
					table: "workflow_rules",
				},
				workflowRule: WorkflowRule{
					CompanyID:                 1,
					MinAmount:                 floatPtr(1000.0),
					MaxAmount:                 floatPtr(5000.0),
					Department:                stringPtr("Finance"),
					IsManagerApprovalRequired: intPtr(0),
					ApproverID:                1,
					ApprovalChannel:           0,
				},
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "database connection failed",
		},
		{
			name: "query row scan fails",
			input: struct {
				store        *workflowRuleStore
				workflowRule WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								scanErr: errors.New("scan error"),
							},
						},
					},
					table: "workflow_rules",
				},
				workflowRule: WorkflowRule{
					CompanyID:                 1,
					MinAmount:                 floatPtr(1000.0),
					MaxAmount:                 floatPtr(5000.0),
					Department:                stringPtr("Finance"),
					IsManagerApprovalRequired: intPtr(0),
					ApproverID:                1,
					ApprovalChannel:           0,
				},
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "scan error",
		},
		{
			name: "commit fails",
			input: struct {
				store        *workflowRuleStore
				workflowRule WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								values: []interface{}{1, 1, floatPtr(1000.0), floatPtr(5000.0), stringPtr("Finance"), intPtr(0), 1, 0},
							},
							commitErr: errors.New("commit failed"),
						},
					},
					table: "workflow_rules",
				},
				workflowRule: WorkflowRule{
					CompanyID:                 1,
					MinAmount:                 floatPtr(1000.0),
					MaxAmount:                 floatPtr(5000.0),
					Department:                stringPtr("Finance"),
					IsManagerApprovalRequired: intPtr(0),
					ApproverID:                1,
					ApprovalChannel:           0,
				},
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "commit failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.store.Create(test.input.workflowRule)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("Create() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("Create() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("Create() unexpected error: %v", gotErr)
				return
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Create() mismatch (-want +got)\n%s", diff)
			}
		})
	}
}

func TestWorkflowRuleStore_List(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store     *workflowRuleStore
			companyID int
		}
		want    []WorkflowRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful workflow rules retrieval by company ID",
			input: struct {
				store     *workflowRuleStore
				companyID int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryResult: &mockSQLRows{
							rows: [][]interface{}{
								{1, 1, floatPtr(1000.0), floatPtr(5000.0), stringPtr("Finance"), intPtr(0), 1, 0},
								{2, 1, floatPtr(5000.0), nil, stringPtr("IT"), intPtr(1), 2, 1},
							},
						},
					},
					table: "workflow_rules",
				},
				companyID: 1,
			},
			want: []WorkflowRule{
				{
					ID:                        1,
					CompanyID:                 1,
					MinAmount:                 floatPtr(1000.0),
					MaxAmount:                 floatPtr(5000.0),
					Department:                stringPtr("Finance"),
					IsManagerApprovalRequired: intPtr(0),
					ApproverID:                1,
					ApprovalChannel:           0,
				},
				{
					ID:                        2,
					CompanyID:                 1,
					MinAmount:                 floatPtr(5000.0),
					MaxAmount:                 nil,
					Department:                stringPtr("IT"),
					IsManagerApprovalRequired: intPtr(1),
					ApproverID:                2,
					ApprovalChannel:           1,
				},
			},
			wantErr: false,
		},
		{
			name: "no workflow rules found for company",
			input: struct {
				store     *workflowRuleStore
				companyID int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryResult: &mockSQLRows{
							rows: [][]interface{}{},
						},
					},
					table: "workflow_rules",
				},
				companyID: 999,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "database query error",
			input: struct {
				store     *workflowRuleStore
				companyID int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryErr: errors.New("database error"),
					},
					table: "workflow_rules",
				},
				companyID: 1,
			},
			want:    nil,
			wantErr: true,
			errMsg:  "failed to query workflow rules: database error",
		},
		{
			name: "scan error during row processing",
			input: struct {
				store     *workflowRuleStore
				companyID int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryResult: &mockSQLRows{
							rows: [][]interface{}{
								{1, 1, floatPtr(1000.0), floatPtr(5000.0), stringPtr("Finance"), intPtr(0), 1, 0},
							},
							scanErr: errors.New("scan error"),
						},
					},
					table: "workflow_rules",
				},
				companyID: 1,
			},
			want:    nil,
			wantErr: true,
			errMsg:  "failed to scan workflow rule: scan error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.store.List(test.input.companyID)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("List() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("List() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("List() unexpected error: %v", gotErr)
				return
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("List() mismatch (-want +got)\n%s", diff)
			}
		})
	}
}

func TestWorkflowRuleStore_FindMatchingRule(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store           *workflowRuleStore
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
				store           *workflowRuleStore
				companyID       int
				amount          float64
				department      string
				requiresManager bool
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							values: []interface{}{1, 1, floatPtr(1000.0), floatPtr(5000.0), stringPtr("Finance"), intPtr(0), 1, 0},
						},
					},
					table: "workflow_rules",
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
				Department:                stringPtr("Finance"),
				IsManagerApprovalRequired: intPtr(0),
				ApproverID:                1,
				ApprovalChannel:           0,
			},
			wantErr: false,
		},
		{
			name: "no matching rule found",
			input: struct {
				store           *workflowRuleStore
				companyID       int
				amount          float64
				department      string
				requiresManager bool
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: sqlpkg.ErrNoRows,
						},
					},
					table: "workflow_rules",
				},
				companyID:       1,
				amount:          10000.0,
				department:      "Unknown",
				requiresManager: false,
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "workflow rule not found",
		},
		{
			name: "database error during rule matching",
			input: struct {
				store           *workflowRuleStore
				companyID       int
				amount          float64
				department      string
				requiresManager bool
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: errors.New("database error"),
						},
					},
					table: "workflow_rules",
				},
				companyID:       1,
				amount:          2500.0,
				department:      "Finance",
				requiresManager: false,
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.store.FindMatchingRule(
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
func stringPtr(s string) *string {
	return &s
}

func TestWorkflowRuleStore_GetByID(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store *workflowRuleStore
			id    int
		}
		want    WorkflowRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful workflow rule retrieval by ID",
			input: struct {
				store *workflowRuleStore
				id    int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							values: []interface{}{
								1, 1, floatPtr(100.0), floatPtr(500.0), stringPtr("Finance"), intPtr(1), 1, 0,
							},
						},
					},
					table: "workflow_rules",
				},
				id: 1,
			},
			want: WorkflowRule{
				ID:                        1,
				CompanyID:                 1,
				MinAmount:                 floatPtr(100.0),
				MaxAmount:                 floatPtr(500.0),
				Department:                stringPtr("Finance"),
				IsManagerApprovalRequired: intPtr(1),
				ApproverID:                1,
				ApprovalChannel:           0,
			},
			wantErr: false,
		},
		{
			name: "workflow rule not found",
			input: struct {
				store *workflowRuleStore
				id    int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: sqlpkg.ErrNoRows,
						},
					},
					table: "workflow_rules",
				},
				id: 999,
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "workflow rule not found",
		},
		{
			name: "database error during retrieval",
			input: struct {
				store *workflowRuleStore
				id    int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: errors.New("database connection failed"),
						},
					},
					table: "workflow_rules",
				},
				id: 1,
			},
			want:    WorkflowRule{},
			wantErr: true,
			errMsg:  "failed to get workflow rule by ID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.store.GetByID(test.input.id)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("GetByID() expected error but got none")
					return
				}
				if test.errMsg != "" && !strings.Contains(gotErr.Error(), test.errMsg) {
					t.Errorf("GetByID() expected error containing %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("GetByID() unexpected error: %v", gotErr)
				return
			}

			if got.ID != test.want.ID {
				t.Errorf("GetByID() ID = %v, want %v", got.ID, test.want.ID)
			}
			if got.CompanyID != test.want.CompanyID {
				t.Errorf("GetByID() CompanyID = %v, want %v", got.CompanyID, test.want.CompanyID)
			}
			if got.ApproverID != test.want.ApproverID {
				t.Errorf("GetByID() ApproverID = %v, want %v", got.ApproverID, test.want.ApproverID)
			}
			if got.ApprovalChannel != test.want.ApprovalChannel {
				t.Errorf("GetByID() ApprovalChannel = %v, want %v", got.ApprovalChannel, test.want.ApprovalChannel)
			}
		})
	}
}

func TestWorkflowRuleStore_Update(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store *workflowRuleStore
			rule  WorkflowRule
		}
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful workflow rule update",
			input: struct {
				store *workflowRuleStore
				rule  WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							queryRowResult: &mockSQLRow{
								values: []interface{}{true}, // Rule exists
							},
							execResult: &mockSQLResult{},
						},
					},
					table: "workflow_rules",
				},
				rule: WorkflowRule{
					ID:                        1,
					CompanyID:                 1,
					MinAmount:                 floatPtr(200.0),
					MaxAmount:                 floatPtr(1000.0),
					Department:                stringPtr("IT"),
					IsManagerApprovalRequired: intPtr(0),
					ApproverID:                2,
					ApprovalChannel:           1,
				},
			},
			wantErr: false,
		},
		{
			name: "workflow rule not found for update",
			input: struct {
				store *workflowRuleStore
				rule  WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							queryRowResult: &mockSQLRow{
								values: []interface{}{false}, // Rule doesn't exist
							},
						},
					},
					table: "workflow_rules",
				},
				rule: WorkflowRule{
					ID: 999,
				},
			},
			wantErr: true,
			errMsg:  "workflow rule not found",
		},
		{
			name: "transaction creation fails",
			input: struct {
				store *workflowRuleStore
				rule  WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						txErr: errors.New("transaction failed"),
					},
					table: "workflow_rules",
				},
				rule: WorkflowRule{
					ID: 1,
				},
			},
			wantErr: true,
			errMsg:  "transaction failed",
		},
		{
			name: "commit fails",
			input: struct {
				store *workflowRuleStore
				rule  WorkflowRule
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							queryRowResult: &mockSQLRow{
								values: []interface{}{true}, // Rule exists
							},
							execResult: &mockSQLResult{},
							commitErr:  errors.New("commit failed"),
						},
					},
					table: "workflow_rules",
				},
				rule: WorkflowRule{
					ID: 1,
				},
			},
			wantErr: true,
			errMsg:  "failed to commit transaction",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := test.input.store.Update(test.input.rule)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("Update() expected error but got none")
					return
				}
				if test.errMsg != "" && !strings.Contains(gotErr.Error(), test.errMsg) {
					t.Errorf("Update() expected error containing %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("Update() unexpected error: %v", gotErr)
			}
		})
	}
}

func TestWorkflowRuleStore_Delete(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store *workflowRuleStore
			id    int
		}
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful workflow rule deletion",
			input: struct {
				store *workflowRuleStore
				id    int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							queryRowResult: &mockSQLRow{
								values: []interface{}{true}, // Rule exists
							},
							execResult: &mockSQLResult{},
						},
					},
					table: "workflow_rules",
				},
				id: 1,
			},
			wantErr: false,
		},
		{
			name: "workflow rule not found for deletion",
			input: struct {
				store *workflowRuleStore
				id    int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							queryRowResult: &mockSQLRow{
								values: []interface{}{false}, // Rule doesn't exist
							},
						},
					},
					table: "workflow_rules",
				},
				id: 999,
			},
			wantErr: true,
			errMsg:  "workflow rule not found",
		},
		{
			name: "transaction creation fails",
			input: struct {
				store *workflowRuleStore
				id    int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						txErr: errors.New("transaction failed"),
					},
					table: "workflow_rules",
				},
				id: 1,
			},
			wantErr: true,
			errMsg:  "transaction failed",
		},
		{
			name: "commit fails",
			input: struct {
				store *workflowRuleStore
				id    int
			}{
				store: &workflowRuleStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							queryRowResult: &mockSQLRow{
								values: []interface{}{true}, // Rule exists
							},
							execResult: &mockSQLResult{},
							commitErr:  errors.New("commit failed"),
						},
					},
					table: "workflow_rules",
				},
				id: 1,
			},
			wantErr: true,
			errMsg:  "failed to commit transaction",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := test.input.store.Delete(test.input.id)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("Delete() expected error but got none")
					return
				}
				if test.errMsg != "" && !strings.Contains(gotErr.Error(), test.errMsg) {
					t.Errorf("Delete() expected error containing %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("Delete() unexpected error: %v", gotErr)
			}
		})
	}
}
