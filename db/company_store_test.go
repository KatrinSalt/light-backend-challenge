package db

import (
	"errors"
	"testing"

	sqlpkg "github.com/KatrinSalt/backend-challenge-go/db/sql"
	"github.com/google/go-cmp/cmp"
)

func TestNewCompanyStore(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			client  sqlpkg.Client
			options []CompanyStoreOption
		}
		want    *companyStore
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid store with defaults",
			input: struct {
				client  sqlpkg.Client
				options []CompanyStoreOption
			}{
				client:  &mockSQLClient{},
				options: []CompanyStoreOption{},
			},
			want: &companyStore{
				client: &mockSQLClient{},
				table:  "companies",
			},
			wantErr: false,
		},
		{
			name: "valid store with custom table name",
			input: struct {
				client  sqlpkg.Client
				options []CompanyStoreOption
			}{
				client: &mockSQLClient{},
				options: []CompanyStoreOption{
					func(o *CompanyStoreOptions) {
						o.Table = "custom_companies"
					},
				},
			},
			want: &companyStore{
				client: &mockSQLClient{},
				table:  "custom_companies",
			},
			wantErr: false,
		},
		{
			name: "nil client",
			input: struct {
				client  sqlpkg.Client
				options []CompanyStoreOption
			}{
				client:  nil,
				options: []CompanyStoreOption{},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "nil sql client",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewCompanyStore(test.input.client, test.input.options...)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("NewCompanyStore() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("NewCompanyStore() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("NewCompanyStore() unexpected error: %v", gotErr)
				return
			}

			if got == nil {
				t.Errorf("NewCompanyStore() returned nil store")
				return
			}

			// Check that the service was created with a client
			if got.client == nil {
				t.Errorf("NewCompanyStore() client should not be nil")
			}
			if got.table != test.want.table {
				t.Errorf("NewCompanyStore() table mismatch, want %q, got %q", test.want.table, got.table)
			}
		})
	}
}

func TestCompanyStore_Create(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store   *companyStore
			company Company
		}
		want    Company
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful company creation",
			input: struct {
				store   *companyStore
				company Company
			}{
				store: &companyStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								values: []interface{}{1, "Light"},
							},
						},
					},
					table: "companies",
				},
				company: Company{Name: "Light"},
			},
			want: Company{
				ID:   1,
				Name: "Light",
			},
			wantErr: false,
		},
		{
			name: "company already exists",
			input: struct {
				store   *companyStore
				company Company
			}{
				store: &companyStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execErr: errors.New("duplicate key value violates unique constraint"),
						},
					},
					table: "companies",
				},
				company: Company{Name: "Light"},
			},
			want:    Company{},
			wantErr: true,
			errMsg:  "duplicate key value violates unique constraint",
		},
		{
			name: "transaction creation fails",
			input: struct {
				store   *companyStore
				company Company
			}{
				store: &companyStore{
					client: &mockSQLClient{
						txErr: errors.New("database connection failed"),
					},
					table: "companies",
				},
				company: Company{Name: "Light"},
			},
			want:    Company{},
			wantErr: true,
			errMsg:  "database connection failed",
		},
		{
			name: "query row scan fails",
			input: struct {
				store   *companyStore
				company Company
			}{
				store: &companyStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								scanErr: errors.New("scan error"),
							},
						},
					},
					table: "companies",
				},
				company: Company{Name: "Light"},
			},
			want:    Company{},
			wantErr: true,
			errMsg:  "scan error",
		},
		{
			name: "commit fails",
			input: struct {
				store   *companyStore
				company Company
			}{
				store: &companyStore{
					client: &mockSQLClient{
						tx: &mockSQLTx{
							execResult: &mockSQLResult{},
							queryRowResult: &mockSQLRow{
								values: []interface{}{1, "Light"},
							},
							commitErr: errors.New("commit failed"),
						},
					},
					table: "companies",
				},
				company: Company{Name: "Light"},
			},
			want:    Company{},
			wantErr: true,
			errMsg:  "commit failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.store.Create(test.input.company)

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

func TestCompanyStore_GetByID(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store *companyStore
			id    int
		}
		want    Company
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful company retrieval by ID",
			input: struct {
				store *companyStore
				id    int
			}{
				store: &companyStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							values: []interface{}{1, "Light"},
						},
					},
					table: "companies",
				},
				id: 1,
			},
			want: Company{
				ID:   1,
				Name: "Light",
			},
			wantErr: false,
		},
		{
			name: "company not found",
			input: struct {
				store *companyStore
				id    int
			}{
				store: &companyStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: sqlpkg.ErrNoRows,
						},
					},
					table: "companies",
				},
				id: 999,
			},
			want:    Company{},
			wantErr: true,
			errMsg:  "company not found",
		},
		{
			name: "database error",
			input: struct {
				store *companyStore
				id    int
			}{
				store: &companyStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: errors.New("database error"),
						},
					},
					table: "companies",
				},
				id: 1,
			},
			want:    Company{},
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

func TestCompanyStore_GetByName(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			store *companyStore
			name  string
		}
		want    Company
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful company retrieval by name",
			input: struct {
				store *companyStore
				name  string
			}{
				store: &companyStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							values: []interface{}{1, "Light"},
						},
					},
					table: "companies",
				},
				name: "Light",
			},
			want: Company{
				ID:   1,
				Name: "Light",
			},
			wantErr: false,
		},
		{
			name: "company not found by name",
			input: struct {
				store *companyStore
				name  string
			}{
				store: &companyStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: sqlpkg.ErrNoRows,
						},
					},
					table: "companies",
				},
				name: "NonExistent",
			},
			want:    Company{},
			wantErr: true,
			errMsg:  "company not found",
		},
		{
			name: "database error on name lookup",
			input: struct {
				store *companyStore
				name  string
			}{
				store: &companyStore{
					client: &mockSQLClient{
						queryRowResult: &mockSQLRow{
							scanErr: errors.New("database error"),
						},
					},
					table: "companies",
				},
				name: "Light",
			},
			want:    Company{},
			wantErr: true,
			errMsg:  "database error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := test.input.store.GetByName(test.input.name)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("GetByName() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("GetByName() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("GetByName() unexpected error: %v", gotErr)
				return
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetByName() mismatch (-want +got)\n%s", diff)
			}
		})
	}
}
