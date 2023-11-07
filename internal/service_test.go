package internal

import (
	"testing"
)

type MockAPIClient struct {
	PostFunc func(reqBody []byte, token string) (*FFLogsAPIResponse, error)
}

func (m *MockAPIClient) Post(reqBody []byte, token string) (*FFLogsAPIResponse, error) {
	return m.PostFunc(reqBody, token)
}

func TestFightIDsService_Run(t *testing.T) {
	mockAPI := &MockAPIClient{
		PostFunc: func(reqBody []byte, token string) (*FFLogsAPIResponse, error) {
			expectedToken := "test_token"
			if token != expectedToken {
				t.Errorf("expected token to be %s, got %s", expectedToken, token)
			}
			return &FFLogsAPIResponse{}, nil
		},
	}
	token := "test_token"
	fs := NewFightIDsService(token, mockAPI)
	reportID := "test_report_id"
	_, err := fs.Run(reportID)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}

func TestTableService_Run(t *testing.T) {
	mockAPI := &MockAPIClient{
		PostFunc: func(reqBody []byte, token string) (*FFLogsAPIResponse, error) {
			expectedToken := "test_token"
			if token != expectedToken {
				t.Errorf("expected token to be %s, got %s", expectedToken, token)
			}
			return &FFLogsAPIResponse{}, nil
		},
	}
	token := "test_token"
	fightID := 1
	ts := NewTableService(token, fightID, mockAPI)
	reportID := "test_report_id"
	_, err := ts.Run(reportID)
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
}
