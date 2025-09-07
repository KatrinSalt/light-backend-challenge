package workflow

// import (
// 	"testing"

// 	"github.com/KatrinSalt/backend-challenge-go/api"
// 	"github.com/KatrinSalt/backend-challenge-go/db"
// )

// func TestNewService(t *testing.T) {
// 	var tests = []struct {
// 		name  string
// 		input struct {
// 			company string
// 			db      databaseService
// 			slack   notificationService
// 			email   notificationService
// 			options []Option
// 		}
// 		wantErr bool
// 	}{
// 		{
// 			name: "valid service creation",
// 			input: struct {
// 				company string
// 				db      databaseService
// 				slack   notificationService
// 				email   notificationService
// 				options []Option
// 			}{
// 				company: "Test Company",
// 				db:      &mockDatabaseService{},
// 				slack:   &mockNotificationService{},
// 				email:   &mockNotificationService{},
// 				options: []Option{},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "empty company name",
// 			input: struct {
// 				company string
// 				db      databaseService
// 				slack   notificationService
// 				email   notificationService
// 				options []Option
// 			}{
// 				company: "",
// 				db:      &mockDatabaseService{},
// 				slack:   &mockNotificationService{},
// 				email:   &mockNotificationService{},
// 				options: []Option{},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "nil database service",
// 			input: struct {
// 				company string
// 				db      databaseService
// 				slack   notificationService
// 				email   notificationService
// 				options []Option
// 			}{
// 				company: "Test Company",
// 				db:      nil,
// 				slack:   &mockNotificationService{},
// 				email:   &mockNotificationService{},
// 				options: []Option{},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "nil slack service",
// 			input: struct {
// 				company string
// 				db      databaseService
// 				slack   notificationService
// 				email   notificationService
// 				options []Option
// 			}{
// 				company: "Test Company",
// 				db:      &mockDatabaseService{},
// 				slack:   nil,
// 				email:   &mockNotificationService{},
// 				options: []Option{},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "nil email service",
// 			input: struct {
// 				company string
// 				db      databaseService
// 				slack   notificationService
// 				email   notificationService
// 				options []Option
// 			}{
// 				company: "Test Company",
// 				db:      &mockDatabaseService{},
// 				slack:   &mockNotificationService{},
// 				email:   nil,
// 				options: []Option{},
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			got, err := NewService(test.input.company, test.input.db, test.input.slack, test.input.email, test.input.options...)

// 			if test.wantErr {
// 				if err == nil {
// 					t.Errorf("NewService() expected error but got none")
// 				}
// 				if got != nil {
// 					t.Errorf("NewService() expected nil service on error but got %v", got)
// 				}
// 			} else {
// 				if err != nil {
// 					t.Errorf("NewService() unexpected error: %v", err)
// 				}
// 				if got == nil {
// 					t.Errorf("NewService() expected service but got nil")
// 				}
// 			}
// 		})
// 	}
// }

// // func TestService_Start(t *testing.T) {
// // 	t.Run("start service", func(t *testing.T) {
// // 		logs := []string{}
// // 		srv := &service{
// // 			log: &mockLogger{
// // 				logs: &logs,
// // 			},
// // 			stopCh: make(chan os.Signal),
// // 			errCh:  make(chan error),
// // 		}
// // 		go func() {
// // 			time.Sleep(time.Millisecond * 100)
// // 			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
// // 		}()
// // 		srv.Start()

// // 		want := []string{
// // 			"Service started.",
// // 			"Service stopped.",
// // 			"reason",
// // 			"interrupt",
// // 		}

// // 		if diff := cmp.Diff(want, logs); diff != "" {
// // 			t.Errorf("Start() = unexpected result (-want +got):\n%s\n", diff)
// // 		}
// // 	})
// // }

// type mockLogger struct {
// 	logs *[]string
// }

// func (l *mockLogger) Info(msg string, args ...any) {
// 	messages := []string{msg}
// 	for _, v := range args {
// 		messages = append(messages, v.(string))
// 	}
// 	*l.logs = append(*l.logs, messages...)
// }

// func (l *mockLogger) Error(msg string, args ...any) {
// 	messages := []string{msg}
// 	for _, v := range args {
// 		messages = append(messages, v.(string))
// 	}
// 	*l.logs = append(*l.logs, messages...)
// }

// // Mock implementations for testing
// type mockDatabaseService struct{}

// func (m *mockDatabaseService) GetCompanyByName(name string) (db.Company, error) {
// 	return db.Company{ID: 1, Name: name}, nil
// }

// func (m *mockDatabaseService) GetApproverByID(id int) (db.Approver, error) {
// 	return db.Approver{ID: id, CompanyID: 1, Name: "Test Approver", Role: "Test Role", Email: "test@example.com", SlackID: "U123456"}, nil
// }

// func (m *mockDatabaseService) FindMatchingRule(companyID int, amount float64, department string, requiresManager bool) (db.WorkflowRule, error) {
// 	return db.WorkflowRule{
// 		ID:                        1,
// 		CompanyID:                 companyID,
// 		MinAmount:                 &amount,
// 		MaxAmount:                 nil,
// 		Department:                &department,
// 		IsManagerApprovalRequired: nil,
// 		ApproverID:                1,
// 		ApprovalChannel:           0,
// 	}, nil
// }

// type mockNotificationService struct{}

// func (m *mockNotificationService) SendApprovalRequest(approvalRequest api.ApprovalRequest) error {
// 	return nil
// }
