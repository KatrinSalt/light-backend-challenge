package db

import (
	"database/sql"
	"errors"

	sqlpkg "github.com/KatrinSalt/backend-challenge-go/db/sql"
)

// Mock implementations for testing
type mockSQLClient struct {
	tx             sqlpkg.Tx
	txErr          error
	queryRowResult sqlpkg.Row
	queryResult    sqlpkg.Rows
	queryErr       error
}

func (m *mockSQLClient) QueryRow(query string, args ...any) sqlpkg.Row {
	return m.queryRowResult
}

func (m *mockSQLClient) Query(query string, args ...any) (sqlpkg.Rows, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}
	return m.queryResult, nil
}

func (m *mockSQLClient) Transaction() (sqlpkg.Tx, error) {
	if m.txErr != nil {
		return nil, m.txErr
	}
	return m.tx, nil
}

func (m *mockSQLClient) Close() error {
	return nil
}

func (m *mockSQLClient) Exec(query string, args ...any) (sqlpkg.Result, error) {
	return nil, nil
}

type mockSQLTx struct {
	execResult     sqlpkg.Result
	execErr        error
	queryRowResult sqlpkg.Row
	commitErr      error
}

func (m *mockSQLTx) QueryRow(query string, args ...any) sqlpkg.Row {
	return m.queryRowResult
}

func (m *mockSQLTx) Exec(query string, args ...any) (sqlpkg.Result, error) {
	if m.execErr != nil {
		return nil, m.execErr
	}
	return m.execResult, nil
}

func (m *mockSQLTx) Commit() error {
	return m.commitErr
}

func (m *mockSQLTx) Rollback() error {
	return nil
}

func (m *mockSQLTx) Prepare(query string) (*sql.Stmt, error) {
	return nil, nil
}

type mockSQLResult struct{}

func (m *mockSQLResult) RowsAffected() (int64, error) {
	return 1, nil
}

type mockSQLRow struct {
	values  []interface{}
	scanErr error
}

func (m *mockSQLRow) Scan(dest ...interface{}) error {
	if m.scanErr != nil {
		return m.scanErr
	}
	if len(m.values) != len(dest) {
		return errors.New("value count mismatch")
	}
	for i, val := range m.values {
		if dest[i] == nil {
			continue
		}
		// Simple type assertion for testing
		switch d := dest[i].(type) {
		case *int:
			if v, ok := val.(int); ok {
				*d = v
			}
		case *string:
			if v, ok := val.(string); ok {
				*d = v
			}
		case **float64:
			if v, ok := val.(*float64); ok {
				*d = v
			}
		case **string:
			if v, ok := val.(*string); ok {
				*d = v
			}
		case **int:
			if v, ok := val.(*int); ok {
				*d = v
			}
		case *bool:
			if v, ok := val.(bool); ok {
				*d = v
			}
		}
	}
	return nil
}

type mockSQLRows struct {
	rows    [][]interface{}
	scanErr error
	index   int
}

func (m *mockSQLRows) Next() bool {
	return m.index < len(m.rows)
}

func (m *mockSQLRows) Scan(dest ...interface{}) error {
	if m.scanErr != nil {
		return m.scanErr
	}
	if m.index >= len(m.rows) {
		return errors.New("no more rows")
	}

	row := m.rows[m.index]
	m.index++

	if len(row) != len(dest) {
		return errors.New("value count mismatch")
	}

	for i, val := range row {
		if dest[i] == nil {
			continue
		}
		// Simple type assertion for testing
		switch d := dest[i].(type) {
		case *int:
			if v, ok := val.(int); ok {
				*d = v
			}
		case *string:
			if v, ok := val.(string); ok {
				*d = v
			}
		case **float64:
			if v, ok := val.(*float64); ok {
				*d = v
			}
		case **string:
			if v, ok := val.(*string); ok {
				*d = v
			}
		case **int:
			if v, ok := val.(*int); ok {
				*d = v
			}
		case *bool:
			if v, ok := val.(bool); ok {
				*d = v
			}
		}
	}
	return nil
}

func (m *mockSQLRows) Close() error {
	return nil
}

func (m *mockSQLRows) Err() error {
	return nil
}
