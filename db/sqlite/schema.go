package sqlite

// dbSchema contains the complete database schema including tables, constraints, and relationships.
var dbSchema = []string{
	// companies table.
	`CREATE TABLE IF NOT EXISTS companies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	)`,
	// approvers table.
	`CREATE TABLE IF NOT EXISTS approvers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		company_id INTEGER NOT NULL,
		email TEXT NOT NULL,
		slack_id TEXT NOT NULL,
		FOREIGN KEY (company_id) REFERENCES companies (id),
		UNIQUE(company_id, email),
		UNIQUE(company_id, slack_id)
	)`,
	// workflow rules table.
	`CREATE TABLE IF NOT EXISTS workflow_rules (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		company_id INTEGER NOT NULL,
		min_amount REAL,
		max_amount REAL,
		department TEXT,
		is_manager_approval_required INTEGER DEFAULT 0 CHECK (is_manager_approval_required IN (0, 1)),
		approver_id INTEGER NOT NULL,
		approval_channel INTEGER NOT NULL,
		FOREIGN KEY (company_id) REFERENCES companies (id),
		FOREIGN KEY (approver_id) REFERENCES approvers (id)
	)`,
}
