package internal

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestFFLogsAPIClient_Post(t *testing.T) {
	responseBody := `{"data":{"reportData":{"report":{"code":"EXAMPLE_CODE"}}}}`	
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(bytes.NewBufferString(responseBody)),
				Header:     make(http.Header),
			}, nil
		},
	}
	client := FFLogsAPIClient{
		Client: mockClient,
	}
	token := "test_token"
	reqBody := []byte(`{"query":"query { test }"}`)
	response, err := client.Post(reqBody, token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	t.Logf("Response: %+v", response)
	expectedCode := "EXAMPLE_CODE"
	if response.Data.ReportData.Report.Code != expectedCode {
		t.Fatalf("expected code to be %s, got %s", expectedCode, response.Data.ReportData.Report.Code)
	}
}
