package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// defaultDriver is the default driver for the sqlite database client.
	defaultDriver = "sqlite3"
	// defaultDataSource is the default data source for the database client.
	defaultDataSource = ":memory:"
)

// client is the database client.
type client struct {
	*sql.DB
}

// ClientOptions contains options for the client.
type ClientOptions struct {
	DataSource string
}

// ClientOption is a function that sets options for the client.
type ClientOption func(o *ClientOptions)

// NewConnection creates a new database connection
func NewClient(options ...ClientOption) (*client, error) {
	opts := ClientOptions{
		DataSource: defaultDataSource,
	}
	for _, option := range options {
		option(&opts)
	}

	db, err := sql.Open(defaultDriver, opts.DataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// TODO: check logging strategy later on.
	log.Printf("Database connection created: %s", opts.DataSource)

	return &client{db}, nil
}

// Close closes the database connection
func (client *client) Close() error {
	return client.DB.Close()
}

// Begin starts a new transaction
func (client *client) Begin() (*sql.Tx, error) {
	return client.DB.Begin()
}

// Exec executes a query without returning rows
func (client *client) Exec(query string, args ...interface{}) (sql.Result, error) {
	return client.DB.Exec(query, args...)
}
