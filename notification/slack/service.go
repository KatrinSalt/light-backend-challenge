package slack

import (
	"fmt"

	"github.com/KatrinSalt/backend-challenge-go/api"
	"github.com/KatrinSalt/backend-challenge-go/common"
)

type service struct {
	log    common.Logger
	client string
}

// Options holds the configuration for the service.
type Options struct {
	Logger common.Logger
}

// Option is a function that configures the service.
type Option func(*service)

// NewService returns a new service.
func NewService(connectionString string, options ...Option) (*service, error) {
	if connectionString == "" {
		return nil, fmt.Errorf("slack client connection string is required to start slack service")
	}
	s := service{
		client: connectionString,
	}

	for _, option := range options {
		option(&s)
	}

	if s.log == nil {
		s.log = common.NewLogger()
	}

	return &s, nil
}

// SendApprovalRequest sends an approval request via email.
func (s *service) SendApprovalRequest(approvalRequest api.ApprovalRequest) (api.ApprovalResponse, error) {
	s.log.Info("Sending approval request via slack",
		"approver_name", approvalRequest.Approver.Name,
		"approver_role", approvalRequest.Approver.Role,
		"approver_slack_id", approvalRequest.Approver.SlackID,
		"invoice_amount", approvalRequest.Invoice.Amount,
	)

	resp := api.ApprovalResponse{
		ApproverName:      approvalRequest.Approver.Name,
		ApproverRole:      approvalRequest.Approver.Role,
		ApproverChannel:   "slack",
		ApproverContactID: approvalRequest.Approver.SlackID,
	}

	return resp, nil
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
