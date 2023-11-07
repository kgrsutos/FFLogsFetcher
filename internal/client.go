package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type FFLogsAPIClient struct {
	Client HTTPClient
}

func (cl *FFLogsAPIClient) Post(reqBody []byte, token string) (*FFLogsAPIResponse, error) {
	req, err := http.NewRequest("POST", "https://ja.fflogs.com/api/v2/client", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := cl.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var fflogsRes FFLogsAPIResponse
	if err := json.Unmarshal(body, &fflogsRes); err != nil {
		return nil, err
	}
	return &fflogsRes, nil
}
