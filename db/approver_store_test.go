package db

import (
	"errors"
	"testing"

	sqlpkg "github.com/KatrinSalt/backend-challenge-go/db/sql"
	"github.com/google/go-cmp/cmp"
)

func TestNewApproverStore(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			client  sqlpkg.Client
			options []ApproverStoreOption
		}
		want    *approverStore
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid store with defaults",
			input: struct {
				client  sqlpkg.Client
				options []ApproverStoreOption
			}{
				client:  &mockSQLClient{},
				options: []ApproverStoreOption{},
			},
			want: &approverStore{
				client: &mockSQLClient{},
				table:  "approvers",
			},
			wantErr: false,
		},
		{
			name: "valid store with custom table name",
			input: struct {
				client  sqlpkg.Client
				options []ApproverStoreOption
			}{
				client: &mockSQLClient{},
				options: []ApproverStoreOption{
					func(o *ApproverStoreOptions) {
						o.Table = "custom_approvers"
					},
				},
			},
			want: &approverStore{
				client: &mockSQLClient{},
				table:  "custom_approvers",
			},
			wantErr: false,
		},
		{
			name: "nil client",
			input: struct {
				client  sqlpkg.Client
				options []ApproverStoreOption
			}{
				client:  nil,
				options: []ApproverStoreOption{},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "nil sql client",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewApproverStore(test.input.client, test.input.options...)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("NewApproverStore() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("NewApproverStore() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("NewApproverStore() unexpected error: %v", gotErr)
				return
			}

			if got == nil {
				t.Errorf("NewApproverStore() returned nil store")
				return
			}

			// Check that the service was created with a client
			if got.client == nil {
				t.Errorf("NewApproverStore() client should not be nil")
			}
			if got.table != test.want.table {
				t.Errorf("NewApproverStore() table mismatch, want %q, got %q", test.want.table, got.table)
			}
		})
	}
}

func TestApproverStore_Create(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store    *approverStore
			approver Approver
		}
		want    Approver
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful approver creation",
			input: struct {
				store    *approverStore
				approver Approver
			}{
				store: &approverStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								values: []interface{}{1, 1, "John Doe", "Manager", "john@example.com", "U123456"},
							},
						},
					},
					table: "approvers",
				},
				approver: Approver{
					CompanyID: 1,
					Name:      "John Doe",
					Role:      "Manager",
					Email:     "john@example.com",
					SlackID:   "U123456",
				},
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
			name: "approver already exists",
			input: struct {
				store    *approverStore
				approver Approver
			}{
				store: &approverStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execErr: errors.New("duplicate key value violates unique constraint"),
						},
					},
					table: "approvers",
				},
				approver: Approver{
					CompanyID: 1,
					Name:      "John Doe",
					Role:      "Manager",
					Email:     "john@example.com",
					SlackID:   "U123456",
				},
			},
			want:    Approver{},
			wantErr: true,
			errMsg:  "duplicate key value violates unique constraint",
		},
		{
			name: "transaction creation fails",
			input: struct {
				store    *approverStore
				approver Approver
			}{
				store: &approverStore{
					client: &mockSQLClient{
						txErr: errors.New("database connection failed"),
					},
					table: "approvers",
				},
				approver: Approver{
					CompanyID: 1,
					Name:      "John Doe",
					Role:      "Manager",
					Email:     "john@example.com",
					SlackID:   "U123456",
				},
			},
			want:    Approver{},
			wantErr: true,
			errMsg:  "database connection failed",
		},
		{
			name: "query row scan fails",
			input: struct {
				store    *approverStore
				approver Approver
			}{
				store: &approverStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								scanErr: errors.New("scan error"),
							},
						},
					},
					table: "approvers",
				},
				approver: Approver{
					CompanyID: 1,
					Name:      "John Doe",
					Role:      "Manager",
					Email:     "john@example.com",
					SlackID:   "U123456",
				},
			},
			want:    Approver{},
			wantErr: true,
			errMsg:  "scan error",
		},
		{
			name: "commit fails",
			input: struct {
				store    *approverStore
				approver Approver
			}{
				store: &approverStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								values: []interface{}{1, 1, "John Doe", "Manager", "john@example.com", "U123456"},
							},
							commitErr: errors.New("commit failed"),
						},
					},
					table: "approvers",
				},
				approver: Approver{
					CompanyID: 1,
					Name:      "John Doe",
					Role:      "Manager",
					Email:     "john@example.com",
					SlackID:   "U123456",
				},
			},
			want:    Approver{},
			wantErr: true,
			errMsg:  "commit failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.store.Create(test.input.approver)

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

func TestApproverStore_GetByID(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store *approverStore
			id    int
		}
		want    Approver
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful approver retrieval by ID",
			input: struct {
				store *approverStore
				id    int
			}{
				store: &approverStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							values: []interface{}{1, 1, "John Doe", "Manager", "john@example.com", "U123456"},
						},
					},
					table: "approvers",
				},
				id: 1,
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
				store *approverStore
				id    int
			}{
				store: &approverStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: sqlpkg.ErrNoRows,
						},
					},
					table: "approvers",
				},
				id: 999,
			},
			want:    Approver{},
			wantErr: true,
			errMsg:  "approver not found",
		},
		{
			name: "database error",
			input: struct {
				store *approverStore
				id    int
			}{
				store: &approverStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: errors.New("database error"),
						},
					},
					table: "approvers",
				},
				id: 1,
			},
			want:    Approver{},
			wantErr: true,
			errMsg:  "database error",
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
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("GetByID() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("GetByID() unexpected error: %v", gotErr)
				return
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetByID() mismatch (-want +got)\n%s", diff)
			}
		})
	}
}
