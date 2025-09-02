package db

// Company represents a company in the system.
type Company struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
