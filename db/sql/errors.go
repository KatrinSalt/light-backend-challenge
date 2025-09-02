package sql

var (
	// SQLStateForeignKey is returned when the inputted value violates foreign key constraints.
	SQLStateForeignKey = "SQLSTATE 23503"
	// SQLStateDuplicateKey is returned when the input value violates constraint
	// set for duplicate keys.
	SQLStateDuplicateKey = "SQLSTATE 23505"
)
