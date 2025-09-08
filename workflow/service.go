package workflow

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/common"
	"github.com/KatrinSalt/backend-challenge-go/db"
)

const (
	slackApprovalChannel approvalChannel = "slack"
	emailApprovalChannel approvalChannel = "email"
)

var (
	// ErrUnsupportedApprovalChannel is returned when an unsupported approval channel is used.
	ErrUnsupportedApprovalChannel = errors.New("unsupported approval channel")
)

// database interface for the database operations.
type databaseService interface {
	GetCompanyByName(name string) (db.Company, error)
	GetApproverByID(id int) (db.Approver, error)
	FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (db.WorkflowRule, error)
}

type notificationService interface {
	SendApprovalRequest(approvalRequest api.ApprovalRequest) (api.ApprovalResponse, error)
}

// Service interface for the workflow service.
type Service interface {
	ValidateCompany() error
	Run() error
}

type service struct {
	log       common.Logger
	company   company
	db        databaseService
	slack     notificationService
	email     notificationService
	reader    *bufio.Reader
	userInput userInput
}

// userInput contains all the user input fields for invoice processing.
type userInput struct {
	amount                    float64
	department                string
	isManagerApprovalRequired bool
}

// company represents a company in the system.
type company struct {
	name        string
	departments []string
}

// Options holds the configuration for the service.
type Options struct {
	Logger common.Logger
}

// Option is a function that configures the service.
type Option func(*service)

// NewService returns a new service.
func NewService(companyName string, companyDepartments []string, db databaseService, slack, email notificationService, options ...Option) (Service, error) {
	if companyName == "" {
		return nil, errors.New("company name is required to start workflow service")
	}
	if len(companyDepartments) == 0 {
		return nil, errors.New("company allowed departments are required to start workflow service")
	}
	if db == nil {
		return nil, errors.New("database service is required to start workflow service")
	}
	if slack == nil {
		return nil, errors.New("slack notification service is required to start workflow service")
	}
	if email == nil {
		return nil, errors.New("email notification service is required to start workflow service")
	}
	s := &service{
		company: company{
			name:        companyName,
			departments: companyDepartments,
		},
		db:     db,
		slack:  slack,
		email:  email,
		reader: bufio.NewReader(os.Stdin),
	}

	for _, option := range options {
		option(s)
	}
	if s.log == nil {
		s.log = common.NewLogger()
	}

	return s, nil
}

// Run starts the CLI service and handles user interaction.
func (s *service) Run() error {
	s.log.Info("Starting CLI service")

	fmt.Println("🧾 Invoice Approval Workflow")
	fmt.Println("==================================")
	fmt.Println()

	// Main application loop - sequential for CLI simplicity.
	for {
		// Get user input.
		if err := s.getUserInput(); err != nil {
			fmt.Printf("❌ Error getting input: %v\n", err)
			fmt.Println()
			continue
		}

		// Display user input for confirmation.
		s.displayUserInput()

		// Process the invoice.
		s.processInvoiceInteractively()

		// Ask if user wants to process another invoice.
		if !s.askToContinue() {
			break
		}
		fmt.Println()
	}

	fmt.Println("The Invoice Approval Workflow is completed!")
	return nil
}

// processInvoiceInteractively processes the invoice interactively.
func (s *service) processInvoiceInteractively() {
	// Convert to invoice request
	invoice := s.toInvoiceRequest()

	// Process the invoice
	fmt.Println("\n🔄 Processing invoice...")
	resp, err := s.processInvoice(invoice)
	if err != nil {
		fmt.Printf("❌ Failed to process invoice: %v\n", err)
	}

	fmt.Println("✅ Invoice processed successfully and sent for approval!")
	fmt.Printf("👔 Approver: %s\n", resp.ApproverName)
	fmt.Printf("👔 Role: %s\n", resp.ApproverRole)
	fmt.Printf("👔 Channel: %s\n", resp.ApproverChannel)
	fmt.Printf("👔 Contact ID: %s\n", resp.ApproverContactID)
	fmt.Println()
}

// getUserInput prompts the user for invoice details and stores them in the service.
func (s *service) getUserInput() error {
	// Get invoice amount.
	amount, err := s.getInvoiceAmount()
	if err != nil {
		return err
	}
	s.userInput.amount = amount

	// Get department.
	department, err := s.getDepartment()
	if err != nil {
		return err
	}
	s.userInput.department = department

	// Get manager approval requirement.
	managerApproval, err := s.getManagerApprovalRequired()
	if err != nil {
		return err
	}
	s.userInput.isManagerApprovalRequired = managerApproval

	return nil
}

// processInvoice processes the invoice.
func (s *service) processInvoice(invoice api.InvoiceRequest) (api.ApprovalResponse, error) {
	// Verify if the company exists in the system.
	companyID, err := s.getCompanyID(invoice.CompanyName)
	if err != nil {
		return api.ApprovalResponse{}, err
	}

	invoiceQ := toInvoiceQuery(companyID, invoice)

	// Find matching rule given the invoice details.
	rule, err := s.findMatchingRule(invoiceQ)
	if err != nil {
		return api.ApprovalResponse{}, err
	}

	// Get approver details.
	approverInfo, err := s.getApproverInfo(rule)
	if err != nil {
		return api.ApprovalResponse{}, err
	}

	// Send approval request.
	resp, err := s.sendApprovalRequest(approverInfo, invoice)
	if err != nil {
		return api.ApprovalResponse{}, err
	}

	return resp, nil

}

// displayUserInput displays the service's userInput in a formatted way.
func (s *service) displayUserInput() {
	fmt.Println("\n📋 Invoice Details:")
	fmt.Println("==================")

	// Display amount (show "Not specified" if 0)
	if s.userInput.amount > 0 {
		fmt.Printf("💰 Amount: $%.2f\n", s.userInput.amount)
	} else {
		fmt.Printf("💰 Amount: Not specified\n")
	}

	// Display department (show "Not specified" if empty)
	if s.userInput.department != "" {
		fmt.Printf("🏢 Department: %s\n", s.userInput.department)
	} else {
		fmt.Printf("🏢 Department: Not specified\n")
	}

	// Display manager approval requirement
	managerApproval := "No"
	if s.userInput.isManagerApprovalRequired {
		managerApproval = "Yes"
	}
	fmt.Printf("👔 Is Manager Approval Required?: %s\n", managerApproval)
}

// getInvoiceAmount prompts the user for invoice amount and validates it.
func (s *service) getInvoiceAmount() (float64, error) {
	fmt.Print("💰 Enter invoice amount (USD) or press Enter to skip: $")

	amountStr, err := s.reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read amount: %v", err)
	}
	amountStr = strings.TrimSpace(amountStr)

	// Allow empty input (skip)
	if amountStr == "" {
		return 0, nil // Return 0 to indicate skipped
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Printf("❌ Error: invalid amount format. Please enter a valid number or press Enter to skip.\n")
		return s.getInvoiceAmount() // Direct recursive call
	}

	if amount <= 0 {
		fmt.Printf("❌ Error: amount must be greater than 0. Please try again or press Enter to skip.\n")
		return s.getInvoiceAmount() // Direct recursive call
	}

	return amount, nil
}

// getDepartment prompts the user for department and validates it.
func (s *service) getDepartment() (string, error) {
	// Get allowed departments from the workflow service
	allowedDepartments := s.getCompanyDepartments()

	// Create a map for case-insensitive comparison
	allowedDeptMap := make(map[string]string)
	for _, dept := range allowedDepartments {
		allowedDeptMap[strings.ToLower(dept)] = dept
	}

	// Create display string for allowed departments
	allowedDeptStr := strings.Join(allowedDepartments, "/")

	fmt.Printf("🏢 Enter department (%s) or press Enter to skip: ", allowedDeptStr)

	department, err := s.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read department: %v", err)
	}
	department = strings.TrimSpace(department)

	// Allow empty input (skip)
	if department == "" {
		return "", nil // Return empty string to indicate skipped
	}

	// Convert to lower case for comparison
	departmentLower := strings.ToLower(department)

	// Check if department is in allowed list (case-insensitive)
	if originalDept, exists := allowedDeptMap[departmentLower]; exists {
		return originalDept, nil // Return the original case from the allowed list
	}

	fmt.Printf("❌ Error: department must be one of: %s. Please try again or press Enter to skip.\n", allowedDeptStr)
	return s.getDepartment() // Direct recursive call
}

// getManagerApprovalRequired prompts the user for manager approval requirement and validates it.
func (s *service) getManagerApprovalRequired() (bool, error) {
	fmt.Print("👔 Does this invoice require manager approval? (y/n) or press Enter to skip: ")

	managerApprovalStr, err := s.reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read manager approval: %v", err)
	}
	managerApprovalStr = strings.TrimSpace(strings.ToLower(managerApprovalStr))

	// Allow empty input (skip) - default to false
	if managerApprovalStr == "" {
		return false, nil // Return false as default when skipped
	}

	switch managerApprovalStr {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		fmt.Printf("❌ Error: please enter 'y' for yes, 'n' for no, or press Enter to skip.\n")
		return s.getManagerApprovalRequired() // Direct recursive call
	}
}

// askToContinue asks the user if they want to process another invoice.
func (s *service) askToContinue() bool {
	fmt.Print("🔄 Process another invoice? (y/n): ")

	response, err := s.reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))

	switch response {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Printf("❌ Error: please enter 'y' for yes or 'n' for no. Please try again.\n")
		return s.askToContinue() // Direct recursive call
	}
}

func (s *service) ValidateCompany() error {
	if _, err := s.getCompanyID(s.company.name); err != nil {
		return err
	}
	return nil
}

// toInvoiceRequest converts the service's userInput to api.InvoiceRequest.
func (s *service) toInvoiceRequest() api.InvoiceRequest {
	return api.InvoiceRequest{
		CompanyName:               s.company.name,
		Amount:                    s.userInput.amount,
		Department:                s.userInput.department,
		IsManagerApprovalRequired: s.userInput.isManagerApprovalRequired,
	}
}

// GetCompanyDepartments returns the company departments.
func (s *service) getCompanyDepartments() []string {
	return s.company.departments
}

// getCompanyID returns the company ID.
func (s *service) getCompanyID(companyName string) (int, error) {
	company, err := s.db.GetCompanyByName(companyName)
	if err != nil {
		s.log.Error("failed to find company in the system", "company_name", companyName, "error", err)
		return 0, err
	}
	return company.ID, nil
}

// findMatchingRule finds the matching rule given the invoice details.
func (s *service) findMatchingRule(q invoiceQuery) (db.WorkflowRule, error) {
	// Find matching rule given the invoice details.
	rule, err := s.db.FindMatchingRule(q.companyID, q.amount, q.department, q.isManagerApprovalRequired)
	if err != nil {
		s.log.Error("failed to find matching workflow rule", "error", err)
		return db.WorkflowRule{}, err
	}
	return rule, nil
}

// getApproverInfo returns the approver information.
func (s *service) getApproverInfo(rule db.WorkflowRule) (approver, error) {
	a, err := s.db.GetApproverByID(rule.ApproverID)
	if err != nil {
		s.log.Error("failed to find approver in the system", "approver_id", rule.ApproverID, "error", err)
		return approver{}, err
	}

	// Determine notification channel.
	var channel approvalChannel

	switch rule.ApprovalChannel {
	case 0: // Slack
		channel = slackApprovalChannel
	case 1: // Email
		channel = emailApprovalChannel
	default:
		// Unsupported approval channel.
		channel = ""
	}

	return approver{
		approver:        toAPIApprover(a),
		approvalChannel: channel,
	}, nil
}

// sendApprovalRequest sends the approval request.
func (s *service) sendApprovalRequest(approver approver, invoiceReq api.InvoiceRequest) (api.ApprovalResponse, error) {
	approvalRequest := toApprovalRequest(approver, invoiceReq)

	// Validate the approval request.
	if err := approvalRequest.Validate(); err != nil {
		return api.ApprovalResponse{}, err
	}

	// Send approval request.
	switch approver.approvalChannel {
	case slackApprovalChannel:
		return s.slack.SendApprovalRequest(approvalRequest)
	case emailApprovalChannel:
		return s.email.SendApprovalRequest(approvalRequest)
	default:
		return api.ApprovalResponse{}, ErrUnsupportedApprovalChannel
	}
}

// toAPIApprover converts the database approver to an API approver.
func toAPIApprover(a db.Approver) api.Approver {
	return api.Approver{
		Name:    a.Name,
		Role:    a.Role,
		Email:   a.Email,
		SlackID: a.SlackID,
	}
}

// toApprovalRequest converts the approver and invoice request to an approval request.
func toApprovalRequest(approver approver, invoiceReq api.InvoiceRequest) api.ApprovalRequest {
	invoice := api.InvoiceDetails{
		Amount: invoiceReq.Amount,
	}

	return api.ApprovalRequest{
		Approver: approver.approver,
		Invoice:  invoice,
	}
}

// toInvoiceQuery converts the invoice request to an invoice query.
func toInvoiceQuery(id int, invoiceReq api.InvoiceRequest) invoiceQuery {
	return invoiceQuery{
		companyID:                 id,
		amount:                    invoiceReq.Amount,
		department:                invoiceReq.Department,
		isManagerApprovalRequired: invoiceReq.IsManagerApprovalRequired,
	}
}

// WithOptions configures the service with the given Options.
func WithOptions(options Options) Option {
	return func(s *service) {
		if options.Logger != nil {
			s.log = options.Logger
		}
	}
}

// WithLogger configures the service with the given logger.
func WithLogger(logger common.Logger) Option {
	return func(s *service) {
		s.log = logger
	}
}
