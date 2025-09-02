package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_sql "github.com/KatrinSalt/backend-challenge-go/db/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// SQLDriver is the sqlite driver.
	SQLDriver = "sqlite3"
	// defaultDataSource is the default data source for the database client.
	defaultDataSource = ":memory:"
)

// client is the database client.
type client struct {
	db *sql.DB
}

// ClientOptions contains options for the client.
type ClientOptions struct {
	DataSource string
}

// ClientOption is a function that sets options for the client.
type ClientOption func(o *ClientOptions)

// NewCLient creates a new sqlite database client.
func NewClient(options ...ClientOption) (*client, error) {
	opts := ClientOptions{
		DataSource: defaultDataSource,
	}
	for _, option := range options {
		option(&opts)
	}

	db, err := sql.Open(SQLDriver, opts.DataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// TODO: check logging strategy later on.
	log.Printf("Database connection created: %s", opts.DataSource)

	return &client{db}, nil
}

// QueryRow executes a query that is expected to return at most one row.
func (c client) QueryRow(query string, args ...any) _sql.Row {
	return c.db.QueryRow(query, args...)
}

// Query executes a query that returns rows.
func (c client) Query(query string, args ...any) (_sql.Rows, error) {
	return c.db.Query(query, args...)
}

// Exec executes a query without returning any rows.
func (c client) Exec(query string, args ...any) (_sql.Result, error) {
	return c.db.Exec(query, args...)
}

// Transaction starts a new database transaction.
func (c client) Transaction() (_sql.Tx, error) {
	t, err := c.db.Begin()
	return tx{t}, err
}

// Close the database connection and release any open resources.
func (c client) Close() error {
	err := c.db.Close()
	if err != nil || errors.Is(err, sql.ErrConnDone) {
		return nil
	}
	return err
}

// tx is a transaction it wraps a *sql.Tx.
type tx struct {
	*sql.Tx
}

// QueryRow executes a query that is expected to return at most one row.
func (t tx) QueryRow(query string, args ...any) _sql.Row {
	return t.Tx.QueryRow(query, args...)
}

// Exec executes a query without returning any rows.
func (t tx) Exec(query string, args ...any) (_sql.Result, error) {
	return t.Tx.Exec(query, args...)
}

// Commit commits the transaction.
func (t tx) Commit() error {
	return t.Tx.Commit()
}

// Rollback rolls back the transaction.
func (t tx) Rollback() error {
	return t.Tx.Rollback()
}

// Prepare creates a prepared statement for use within a transaction.
func (t tx) Prepare(query string) (*sql.Stmt, error) {
	return t.Tx.Prepare(query)
}
