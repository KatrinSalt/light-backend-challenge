package sql

import (
	"database/sql"
)

// ErrNoRows is returned by [Row.Scan] when [DB.QueryRow] doesn't
// return a row. In such a case,QueryRow returns a placeholder[*Row]
// value that defers this error until a Scan.
var ErrNoRows = sql.ErrNoRows

// Row is a result row.
type Row interface {
	Scan(dest ...any) error
}

// Rows is a result set.
type Rows interface {
	Scan(dest ...any) error
	Next() bool
	Close() error
	Err() error
}

// Result is the result of a query.
type Result interface {
	RowsAffected() (int64, error)
}

// Client is the interface for the database client.
type Client interface {
	QueryRow(query string, args ...any) Row
	Query(query string, args ...any) (Rows, error)
	Exec(query string, args ...any) (Result, error)
	Transaction() (Tx, error)
	Close() error
}

// Tx is the interface for the database transaction.
type Tx interface {
	QueryRow(query string, args ...any) Row
	Exec(query string, args ...any) (Result, error)
	Commit() error
	Rollback() error
	Prepare(query string) (*sql.Stmt, error)
}
