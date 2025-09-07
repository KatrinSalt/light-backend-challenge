package workflow

import (
	"errors"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/common"
	"github.com/KatrinSalt/backend-challenge-go/db"
)

const (
	slackApprovalChannel approvalChannel = "slack"
	emailApprovalChannel approvalChannel = "email"
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
	ProcessInvoice(invoice api.InvoiceRequest) (api.ApprovalResponse, error)
	ValidateCompany() error
	GetCompanyName() string
	GetCompanyDepartments() []string
}

type service struct {
	log     common.Logger
	company Company
	db      databaseService
	slack   notificationService
	email   notificationService
}

// Company represents a company in the system.
type Company struct {
	Name        string
	Departments []string
}

// Options holds the configuration for the service.
type Options struct {
	Logger common.Logger
}

// Option is a function that configures the service.
type Option func(*service)

// NewService returns a new service.
func NewService(company Company, db databaseService, slack, email notificationService, options ...Option) (Service, error) {
	if company.Name == "" {
		return nil, errors.New("company name is required to start workflow service")
	}
	if len(company.Departments) == 0 {
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
		company: company,
		db:      db,
		slack:   slack,
		email:   email,
	}

	for _, option := range options {
		option(s)
	}
	if s.log == nil {
		s.log = common.NewLogger()
	}

	return s, nil
}

func (s *service) ValidateCompany() error {
	if _, err := s.getCompanyID(s.company.Name); err != nil {
		return err
	}
	return nil
}

// Start the service.
func (s *service) ProcessInvoice(invoice api.InvoiceRequest) (api.ApprovalResponse, error) {
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

func (s *service) GetCompanyName() string {
	return s.company.Name
}

func (s *service) GetCompanyDepartments() []string {
	return s.company.Departments
}

func (s *service) getCompanyID(companyName string) (int, error) {
	company, err := s.db.GetCompanyByName(companyName)
	if err != nil {
		s.log.Error("failed to find company in the system", "company_name", companyName, "error", err)
		return 0, err
	}
	return company.ID, nil
}

func (s *service) findMatchingRule(q invoiceQuery) (db.WorkflowRule, error) {
	// Find matching rule given the invoice details.
	rule, err := s.db.FindMatchingRule(q.companyID, q.amount, q.department, q.isManagerApprovalRequired)
	if err != nil {
		s.log.Error("failed to find matching workflow rule", "error", err)
		return db.WorkflowRule{}, err
	}
	return rule, nil
}

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
