package db

// SampleData contains all the sample data for seeding the database.
type SampleData struct {
	Companies     []Company
	Approvers     []Approver
	WorkflowRules []WorkflowRule
}

// NewSampleData creates and returns sample data based on Figure 1 from the README.md.
func NewSampleData() *SampleData {
	return &SampleData{
		Companies:     getSampleCompanies(),
		Approvers:     getSampleApprovers(),
		WorkflowRules: getSampleWorkflowRules(),
	}
}

// Default values related to the sample company.
const (
	companyID           = 1
	marketingDepartment = "Marketing"
)

// getSampleCompanies returns sample company data.
func getSampleCompanies() []Company {
	return []Company{
		{
			Name: "Light",
		},
	}
}

const (
	approverID1 = 1
	approverID2 = 2
	approverID3 = 3
	approverID4 = 4
)

// getSampleApprovers returns sample approver data.
func getSampleApprovers() []Approver {
	return []Approver{
		// finance team member.
		{
			ID:        approverID1,
			CompanyID: companyID,
			Email:     "finance_team@light.com",
			SlackID:   "U123456",
		},
		// finance department manager.
		{
			ID:        approverID2,
			CompanyID: companyID,
			Email:     "vera_sander@light.com",
			SlackID:   "U789012",
		},
		// Chief Financial Officer (CFO).
		{
			ID:        approverID3,
			CompanyID: companyID,
			Email:     "amanda_svensson@light.com",
			SlackID:   "U345678",
		},
		// Chief Marketing Officer (CMO).
		{
			ID:        approverID4,
			CompanyID: companyID,
			Email:     "sarah_johnson@light.com",
			SlackID:   "U456789",
		},
	}
}

// ApprovalChannel represents the channel for sending approval requests.
type ApprovalChannel int

const (
	slack ApprovalChannel = iota // 0
	email                        // 1
)

// String returns the string representation of the approval channel.
func (ac ApprovalChannel) String() string {
	switch ac {
	case slack:
		return "slack"
	case email:
		return "email"
	default:
		return ""
	}
}

// getSampleWorkflowRules returns sample workflow rules based on Figure 1
func getSampleWorkflowRules() []WorkflowRule {
	// Default values for workflow rules.
	// Rule 1: Send an approval request to any finance team member via Slack when invoice < $5k.
	rule1MaxAmount := 5000.0
	rule1ApproverID := approverID1
	rule1ApprovalChannel := int(slack)

	// Rule 2: Send an approval request to any finance team member via Email when  $5k <= invoice < $10k,
	// and if manager approval is not required.
	rule2MinAmount := 5000.0
	rule2MaxAmount := 10000.0
	rule2ApproverID := approverID1
	rule2ApprovalChannel := int(email)

	// Rule 3: Send an approval request to finance department manager via Email when  $5k <= invoice < $10k,
	// and if manager approval is required.
	rule3MinAmount := 5000.0
	rule3MaxAmount := 10000.0
	rule3ApproverID := approverID2
	rule3ApprovalChannel := int(email)
	rule3IsManagerApprovalRequired := 1

	// Rule 4: Send an approval request to CFO via Slack when invoice >= $10k,
	// and if invoice is not related to marketing.
	rule4MinAmount := 10000.0
	rule4ApproverID := approverID3
	rule4ApprovalChannel := int(slack)

	// Rule 5: Send an approval request to CMO via Email when invoice >= $10k,
	// and if invoice is related to marketing.
	rule5MinAmount := 10000.0
	rule5ApproverID := approverID4
	rule5Department := marketingDepartment
	rule5ApprovalChannel := int(email)

	return []WorkflowRule{
		// Rule 1: Send an approval request to any finance team member via Slack when invoice < $5k.
		{
			CompanyID:                 companyID,
			MinAmount:                 nil,
			MaxAmount:                 &rule1MaxAmount,
			Department:                nil,
			IsManagerApprovalRequired: nil,                  // Defaults to false
			ApproverID:                rule1ApproverID,      // Finance team member
			ApprovalChannel:           rule1ApprovalChannel, // Slack
		},
		// Rule 2: Send an approval request to any finance team member via Email when  $5k <= invoice < $10k,
		// and if manager approval is not required.
		{
			CompanyID:                 companyID,
			MinAmount:                 &rule2MinAmount,
			MaxAmount:                 &rule2MaxAmount,
			Department:                nil,
			IsManagerApprovalRequired: nil,                  // Defaults to false
			ApproverID:                rule2ApproverID,      // Finance team member
			ApprovalChannel:           rule2ApprovalChannel, // Email
		},
		// Rule 3: Send an approval request to finance department manager via Email when  $5k <= invoice < $10k,
		// and if manager approval is required.
		{
			CompanyID:                 companyID,
			MinAmount:                 &rule3MinAmount,
			MaxAmount:                 &rule3MaxAmount,
			Department:                nil,
			IsManagerApprovalRequired: &rule3IsManagerApprovalRequired, // true
			ApproverID:                rule3ApproverID,                 // Finance department manager
			ApprovalChannel:           rule3ApprovalChannel,            // Email
		},
		// Rule 4: Send an approval request to CFO via Slack when invoice >= $10k,
		// and if invoice is not related to marketing.
		{
			CompanyID:                 companyID,
			MinAmount:                 &rule4MinAmount,
			MaxAmount:                 nil,
			Department:                nil,
			IsManagerApprovalRequired: nil,                  // Defaults to false
			ApproverID:                rule4ApproverID,      // CFO
			ApprovalChannel:           rule4ApprovalChannel, // Slack
		},
		// Rule 5: Send an approval request to CMO via Email when invoice >= $10k,
		// and if invoice is related to marketing.
		{
			CompanyID:                 companyID,
			MinAmount:                 &rule5MinAmount,
			MaxAmount:                 nil,
			Department:                &rule5Department,     // Marketing department
			IsManagerApprovalRequired: nil,                  // Defaults to false
			ApproverID:                rule5ApproverID,      // CMO
			ApprovalChannel:           rule5ApprovalChannel, // Email
		},
	}
}

// // SeedSampleData populates the database with the provided workflow example (Fig.1).
// func SeedSampleData(client sql.Client) error {
// 	log.Println("Seeding database with sample data...")

// 	// Get sample data.
// 	sampleData := NewSampleData()

// 	// Start a transaction.
// 	tx, err := client.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	// Rollback the transaction if an error occurs.
// 	defer tx.Rollback()

// 	// Insert companies data.
// 	stmt, err := tx.Prepare("INSERT OR IGNORE INTO companies (name) VALUES (?)")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	for _, company := range sampleData.Companies {
// 		if _, err := stmt.Exec(company.Name); err != nil {
// 			return err
// 		}
// 	}

// 	// Insert approvers data.
// 	approverStmt, err := tx.Prepare("INSERT OR IGNORE INTO approvers (id, company_id, email, slack_id) VALUES (?, ?, ?, ?)")
// 	if err != nil {
// 		return err
// 	}
// 	defer approverStmt.Close()

// 	for _, approver := range sampleData.Approvers {
// 		if _, err := approverStmt.Exec(approver.ID, approver.CompanyID, approver.Email, approver.SlackID); err != nil {
// 			return err
// 		}
// 	}

// 	// Insert workflow rules data.
// 	ruleStmt, err := tx.Prepare("INSERT OR IGNORE INTO workflow_rules (id, company_id, min_amount, max_amount, department, is_manager_approval_required, approver_id, approval_channel) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
// 	if err != nil {
// 		return err
// 	}
// 	defer ruleStmt.Close()

// 	for _, rule := range sampleData.WorkflowRules {
// 		if _, err := ruleStmt.Exec(rule.ID, rule.CompanyID, rule.MinAmount, rule.MaxAmount, rule.Department, rule.IsManagerApprovalRequired, rule.ApproverID, rule.ApprovalChannel); err != nil {
// 			return err
// 		}
// 	}

// 	// Commit the transaction
// 	if err := tx.Commit(); err != nil {
// 		return err
// 	}

// 	// TODO: check logging strategy later on.
// 	log.Println("Database seeded successfully")

// 	return nil
// }
