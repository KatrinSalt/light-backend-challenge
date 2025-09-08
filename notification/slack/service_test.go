package slack

import (
	"testing"

	"github.com/KatrinSalt/backend-challenge-go/api"
)

func TestNewService(t *testing.T) {
	var tests = []struct {
		name       string
		connection string
		options    []Option
		wantErr    bool
	}{
		{
			name:       "valid service creation",
			connection: "slack://test",
			options:    []Option{},
			wantErr:    false,
		},
		{
			name:       "valid service with custom logger",
			connection: "slack://test",
			options:    []Option{WithLogger(&mockLogger{})},
			wantErr:    false,
		},
		{
			name:       "valid service with options",
			connection: "slack://test",
			options:    []Option{WithOptions(Options{Logger: &mockLogger{}})},
			wantErr:    false,
		},
		{
			name:       "empty connection string",
			connection: "",
			options:    []Option{},
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewService(test.connection, test.options...)

			if test.wantErr {
				if err == nil {
					t.Errorf("NewService() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("NewService() unexpected error: %v", err)
				}
				if got == nil {
					t.Errorf("NewService() expected service but got nil")
				}
				if got != nil && got.client != test.connection {
					t.Errorf("NewService() client = %v, want %v", got.client, test.connection)
				}
			}
		})
	}
}

func TestService_SendApprovalRequest(t *testing.T) {
	var tests = []struct {
		name            string
		service         *service
		approvalRequest api.ApprovalRequest
		wantResp        api.ApprovalResponse
		wantErr         bool
	}{
		{
			name: "successful slack approval request",
			service: &service{
				log:    &mockLogger{},
				client: "slack://test",
			},
			approvalRequest: api.ApprovalRequest{
				Approver: api.Approver{
					Name:    "John Doe",
					Role:    "Manager",
					Email:   "john@example.com",
					SlackID: "U123456",
				},
				Invoice: api.InvoiceDetails{
					Amount: 500.0,
				},
			},
			wantResp: api.ApprovalResponse{
				ApproverName:      "John Doe",
				ApproverRole:      "Manager",
				ApproverChannel:   "slack",
				ApproverContactID: "U123456",
			},
			wantErr: false,
		},
		{
			name: "approval request with empty approver name",
			service: &service{
				log:    &mockLogger{},
				client: "slack://test",
			},
			approvalRequest: api.ApprovalRequest{
				Approver: api.Approver{
					Name:    "", // Empty name
					Role:    "Manager",
					Email:   "test@example.com",
					SlackID: "U123456",
				},
				Invoice: api.InvoiceDetails{
					Amount: 500.0,
				},
			},
			wantResp: api.ApprovalResponse{
				ApproverName:      "", // Should preserve empty name
				ApproverRole:      "Manager",
				ApproverChannel:   "slack",
				ApproverContactID: "U123456",
			},
			wantErr: false,
		},
		{
			name: "approval request with empty role",
			service: &service{
				log:    &mockLogger{},
				client: "slack://test",
			},
			approvalRequest: api.ApprovalRequest{
				Approver: api.Approver{
					Name:    "Test User",
					Role:    "", // Empty role
					Email:   "test@example.com",
					SlackID: "U123456",
				},
				Invoice: api.InvoiceDetails{
					Amount: 500.0,
				},
			},
			wantResp: api.ApprovalResponse{
				ApproverName:      "Test User",
				ApproverRole:      "", // Should preserve empty role
				ApproverChannel:   "slack",
				ApproverContactID: "U123456",
			},
			wantErr: false,
		},
		{
			name: "approval request with special characters in name",
			service: &service{
				log:    &mockLogger{},
				client: "slack://test",
			},
			approvalRequest: api.ApprovalRequest{
				Approver: api.Approver{
					Name:    "José María O'Connor-Smith",
					Role:    "Senior Manager",
					Email:   "jose@example.com",
					SlackID: "U123456",
				},
				Invoice: api.InvoiceDetails{
					Amount: 500.0,
				},
			},
			wantResp: api.ApprovalResponse{
				ApproverName:      "José María O'Connor-Smith",
				ApproverRole:      "Senior Manager",
				ApproverChannel:   "slack",
				ApproverContactID: "U123456",
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.service.SendApprovalRequest(test.approvalRequest)

			if test.wantErr {
				if err == nil {
					t.Errorf("SendApprovalRequest() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("SendApprovalRequest() unexpected error: %v", err)
				}
				if got != test.wantResp {
					t.Errorf("SendApprovalRequest() = %v, want %v", got, test.wantResp)
				}
			}
		})
	}
}

// Mock implementations for testing
type mockLogger struct{}

func (l *mockLogger) Info(msg string, args ...any) {
	// Mock implementation - just ignore for testing
}

func (l *mockLogger) Error(msg string, args ...any) {
	// Mock implementation - just ignore for testing
}
