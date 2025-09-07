package cli

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/KatrinSalt/backend-challenge-go/api"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			workflowSvc workflowService
			options     []Option
		}
		want    Service
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid service with defaults",
			input: struct {
				workflowSvc workflowService
				options     []Option
			}{
				workflowSvc: &mockWorkflowService{},
				options:     []Option{},
			},
			want:    &service{workflowService: &mockWorkflowService{}},
			wantErr: false,
		},
		{
			name: "valid service with custom logger",
			input: struct {
				workflowSvc workflowService
				options     []Option
			}{
				workflowSvc: &mockWorkflowService{},
				options: []Option{
					WithLogger(&mockLogger{}),
				},
			},
			want:    &service{workflowService: &mockWorkflowService{}, log: &mockLogger{}},
			wantErr: false,
		},
		{
			name: "valid service with options",
			input: struct {
				workflowSvc workflowService
				options     []Option
			}{
				workflowSvc: &mockWorkflowService{},
				options: []Option{
					WithOptions(Options{
						Logger: &mockLogger{},
					}),
				},
			},
			want:    &service{workflowService: &mockWorkflowService{}, log: &mockLogger{}},
			wantErr: false,
		},
		{
			name: "nil workflow service",
			input: struct {
				workflowSvc workflowService
				options     []Option
			}{
				workflowSvc: nil,
				options:     []Option{},
			},
			want:    nil,
			wantErr: true,
			errMsg:  "workflow service is required to start CLI service",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewService(test.input.workflowSvc, test.input.options...)

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("NewService() expected error but got none")
					return
				}
				if test.errMsg != "" && gotErr.Error() != test.errMsg {
					t.Errorf("NewService() expected error %q but got %q", test.errMsg, gotErr.Error())
				}
				return
			}

			if gotErr != nil {
				t.Errorf("NewService() unexpected error: %v", gotErr)
				return
			}

			if got == nil {
				t.Errorf("NewService() returned nil service")
				return
			}

			// Check that the service was created with a workflow service
			svc := got.(*service)
			if svc.workflowService == nil {
				t.Errorf("NewService() workflowService should not be nil")
			}
			if svc.reader == nil {
				t.Errorf("NewService() reader should not be nil")
			}
		})
	}
}

func TestService_ToInvoiceRequest(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			service   *service
			userInput userInput
		}
		want api.InvoiceRequest
	}{
		{
			name: "complete user input",
			input: struct {
				service   *service
				userInput userInput
			}{
				service: &service{
					workflowService: &mockWorkflowService{
						companyName: "Light",
					},
				},
				userInput: userInput{
					amount:                    1500.50,
					department:                "Finance",
					isManagerApprovalRequired: true,
				},
			},
			want: api.InvoiceRequest{
				CompanyName:               "Light",
				Amount:                    1500.50,
				Department:                "Finance",
				IsManagerApprovalRequired: true,
			},
		},
		{
			name: "minimal user input",
			input: struct {
				service   *service
				userInput userInput
			}{
				service: &service{
					workflowService: &mockWorkflowService{
						companyName: "Light",
					},
				},
				userInput: userInput{
					amount:                    0,
					department:                "",
					isManagerApprovalRequired: false,
				},
			},
			want: api.InvoiceRequest{
				CompanyName:               "Light",
				Amount:                    0,
				Department:                "",
				IsManagerApprovalRequired: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.input.service.userInput = test.input.userInput
			got := test.input.service.toInvoiceRequest()

			if got.CompanyName != test.want.CompanyName {
				t.Errorf("toInvoiceRequest() CompanyName = %q, want %q", got.CompanyName, test.want.CompanyName)
			}
			if got.Amount != test.want.Amount {
				t.Errorf("toInvoiceRequest() Amount = %v, want %v", got.Amount, test.want.Amount)
			}
			if got.Department != test.want.Department {
				t.Errorf("toInvoiceRequest() Department = %q, want %q", got.Department, test.want.Department)
			}
			if got.IsManagerApprovalRequired != test.want.IsManagerApprovalRequired {
				t.Errorf("toInvoiceRequest() IsManagerApprovalRequired = %v, want %v", got.IsManagerApprovalRequired, test.want.IsManagerApprovalRequired)
			}
		})
	}
}

func TestService_DisplayUserInput(t *testing.T) {
	tests := []struct {
		name  string
		input userInput
		want  []string
	}{
		{
			name: "complete input",
			input: userInput{
				amount:                    1500.50,
				department:                "Finance",
				isManagerApprovalRequired: true,
			},
			want: []string{
				"💰 Amount: $1500.50",
				"🏢 Department: Finance",
				"👔 Is Manager Approval Required?: Yes",
			},
		},
		{
			name: "zero amount",
			input: userInput{
				amount:                    0,
				department:                "IT",
				isManagerApprovalRequired: false,
			},
			want: []string{
				"💰 Amount: Not specified",
				"🏢 Department: IT",
				"👔 Is Manager Approval Required?: No",
			},
		},
		{
			name: "empty department",
			input: userInput{
				amount:                    500.0,
				department:                "",
				isManagerApprovalRequired: true,
			},
			want: []string{
				"💰 Amount: $500.00",
				"🏢 Department: Not specified",
				"👔 Is Manager Approval Required?: Yes",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &service{
				userInput: test.input,
			}

			// Capture output.
			output := captureOutput(func() {
				service.displayUserInput()
			})

			// Check that all expected strings are in the output.
			for _, expected := range test.want {
				if !strings.Contains(output, expected) {
					t.Errorf("displayUserInput() output should contain %q, got: %s", expected, output)
				}
			}
		})
	}
}

func TestService_GetInvoiceAmount(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{
			name:    "valid amount",
			input:   "1500.50\n",
			want:    1500.50,
			wantErr: false,
		},
		{
			name:    "valid integer amount",
			input:   "1000\n",
			want:    1000.0,
			wantErr: false,
		},
		{
			name:    "empty input (skip)",
			input:   "\n",
			want:    0,
			wantErr: false,
		},
		{
			name:    "whitespace only (skip)",
			input:   "   \n",
			want:    0,
			wantErr: false,
		},
		{
			name:    "zero amount",
			input:   "0\n",
			want:    100, // After recursion, it will return the valid input
			wantErr: false,
		},
		{
			name:    "negative amount",
			input:   "-100\n",
			want:    100, // After recursion, it will return the valid input
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "abc\n",
			want:    100, // After recursion, it will return the valid input
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// For tests that would cause recursive calls, a valid input needs to be provided
			// after the invalid input to break the recursion.
			testInput := test.input
			if test.name == "zero amount" || test.name == "negative amount" || test.name == "invalid format" {
				testInput = test.input + "100\n" // Provide valid input to break recursion
			}

			service := &service{
				reader: bufio.NewReader(strings.NewReader(testInput)),
			}

			got, gotErr := service.getInvoiceAmount()

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("getInvoiceAmount() expected error but got none")
					return
				}
				return
			}

			if gotErr != nil {
				t.Errorf("getInvoiceAmount() unexpected error: %v", gotErr)
				return
			}

			if got != test.want {
				t.Errorf("getInvoiceAmount() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestService_GetDepartment(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			userInput    string
			allowedDepts []string
		}
		want    string
		wantErr bool
	}{
		{
			name: "valid department (exact case)",
			input: struct {
				userInput    string
				allowedDepts []string
			}{
				userInput:    "Finance\n",
				allowedDepts: []string{"Finance", "IT", "HR"},
			},
			want:    "Finance",
			wantErr: false,
		},
		{
			name: "valid department (lowercase)",
			input: struct {
				userInput    string
				allowedDepts []string
			}{
				userInput:    "finance\n",
				allowedDepts: []string{"Finance", "IT", "HR"},
			},
			want:    "Finance",
			wantErr: false,
		},
		{
			name: "valid department (uppercase)",
			input: struct {
				userInput    string
				allowedDepts []string
			}{
				userInput:    "IT\n",
				allowedDepts: []string{"Finance", "IT", "HR"},
			},
			want:    "IT",
			wantErr: false,
		},
		{
			name: "empty input (skip)",
			input: struct {
				userInput    string
				allowedDepts []string
			}{
				userInput:    "\n",
				allowedDepts: []string{"Finance", "IT", "HR"},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "whitespace only (skip)",
			input: struct {
				userInput    string
				allowedDepts []string
			}{
				userInput:    "   \n",
				allowedDepts: []string{"Finance", "IT", "HR"},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "invalid department",
			input: struct {
				userInput    string
				allowedDepts []string
			}{
				userInput:    "Marketing\nFinance\n", // Provide valid input after invalid to break recursion
				allowedDepts: []string{"Finance", "IT", "HR"},
			},
			want:    "Finance",
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &service{
				reader: bufio.NewReader(strings.NewReader(test.input.userInput)),
				workflowService: &mockWorkflowService{
					departments: test.input.allowedDepts,
				},
			}

			got, gotErr := service.getDepartment()

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("getDepartment() expected error but got none")
					return
				}
				return
			}

			if gotErr != nil {
				t.Errorf("getDepartment() unexpected error: %v", gotErr)
				return
			}

			if got != test.want {
				t.Errorf("getDepartment() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestService_GetManagerApprovalRequired(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    bool
		wantErr bool
	}{
		{
			name:    "yes (y)",
			input:   "y\n",
			want:    true,
			wantErr: false,
		},
		{
			name:    "yes (yes)",
			input:   "yes\n",
			want:    true,
			wantErr: false,
		},
		{
			name:    "yes (Y)",
			input:   "Y\n",
			want:    true,
			wantErr: false,
		},
		{
			name:    "yes (YES)",
			input:   "YES\n",
			want:    true,
			wantErr: false,
		},
		{
			name:    "no (n)",
			input:   "n\n",
			want:    false,
			wantErr: false,
		},
		{
			name:    "no (no)",
			input:   "no\n",
			want:    false,
			wantErr: false,
		},
		{
			name:    "no (N)",
			input:   "N\n",
			want:    false,
			wantErr: false,
		},
		{
			name:    "no (NO)",
			input:   "NO\n",
			want:    false,
			wantErr: false,
		},
		{
			name:    "empty input (skip)",
			input:   "\n",
			want:    false,
			wantErr: false,
		},
		{
			name:    "whitespace only (skip)",
			input:   "   \n",
			want:    false,
			wantErr: false,
		},
		{
			name:    "invalid input",
			input:   "maybe\ny\n", // Provide valid input after invalid to break recursion
			want:    true,
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &service{
				reader: bufio.NewReader(strings.NewReader(test.input)),
			}

			got, gotErr := service.getManagerApprovalRequired()

			if test.wantErr {
				if gotErr == nil {
					t.Errorf("getManagerApprovalRequired() expected error but got none")
					return
				}
				return
			}

			if gotErr != nil {
				t.Errorf("getManagerApprovalRequired() unexpected error: %v", gotErr)
				return
			}

			if got != test.want {
				t.Errorf("getManagerApprovalRequired() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestService_AskToContinue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "yes (y)",
			input: "y\n",
			want:  true,
		},
		{
			name:  "yes (yes)",
			input: "yes\n",
			want:  true,
		},
		{
			name:  "yes (Y)",
			input: "Y\n",
			want:  true,
		},
		{
			name:  "yes (YES)",
			input: "YES\n",
			want:  true,
		},
		{
			name:  "no (n)",
			input: "n\n",
			want:  false,
		},
		{
			name:  "no (no)",
			input: "no\n",
			want:  false,
		},
		{
			name:  "no (N)",
			input: "N\n",
			want:  false,
		},
		{
			name:  "no (NO)",
			input: "NO\n",
			want:  false,
		},
		{
			name:  "invalid input",
			input: "maybe\n",
			want:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := &service{
				reader: bufio.NewReader(strings.NewReader(test.input)),
			}

			got := service.askToContinue()

			if got != test.want {
				t.Errorf("askToContinue() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestService_ProcessInvoice(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			service   *service
			userInput userInput
		}
		wantErr bool
	}{
		{
			name: "successful processing",
			input: struct {
				service   *service
				userInput userInput
			}{
				service: &service{
					workflowService: &mockWorkflowService{
						companyName: "Light",
						response: api.ApprovalResponse{
							ApproverName:      "John Doe",
							ApproverRole:      "Manager",
							ApproverChannel:   "email",
							ApproverContactID: "john@example.com",
						},
					},
				},
				userInput: userInput{
					amount:                    1500.50,
					department:                "Finance",
					isManagerApprovalRequired: true,
				},
			},
			wantErr: false,
		},
		{
			name: "workflow service error",
			input: struct {
				service   *service
				userInput userInput
			}{
				service: &service{
					workflowService: &mockWorkflowService{
						companyName: "Light",
						processErr:  errors.New("workflow processing failed"),
					},
				},
				userInput: userInput{
					amount:                    1500.50,
					department:                "Finance",
					isManagerApprovalRequired: true,
				},
			},
			wantErr: false, // CLI service doesn't return error, just displays it
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.input.service.userInput = test.input.userInput

			// Capture output
			output := captureOutput(func() {
				test.input.service.processInvoice()
			})

			// Check that processing message is displayed
			if !strings.Contains(output, "Processing invoice...") {
				t.Errorf("processInvoice() should display processing message")
			}

			if test.wantErr {
				if !strings.Contains(output, "Failed to process invoice") {
					t.Errorf("processInvoice() should display error message")
				}
			} else {
				if !strings.Contains(output, "Invoice processed successfully") {
					t.Errorf("processInvoice() should display success message")
				}
			}
		})
	}
}

// Helper function to capture output
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Mock implementations for testing
type mockWorkflowService struct {
	companyName string
	departments []string
	response    api.ApprovalResponse
	processErr  error
}

func (m *mockWorkflowService) ProcessInvoice(invoice api.InvoiceRequest) (api.ApprovalResponse, error) {
	if m.processErr != nil {
		return api.ApprovalResponse{}, m.processErr
	}
	return m.response, nil
}

func (m *mockWorkflowService) GetCompanyName() string {
	return m.companyName
}

func (m *mockWorkflowService) GetCompanyDepartments() []string {
	return m.departments
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...interface{})  {}
func (m *mockLogger) Error(msg string, args ...interface{}) {}
