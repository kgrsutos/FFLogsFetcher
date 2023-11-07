package internal

import (
	"encoding/json"
)

type fightIDsService struct {
	stmt string
	token string
	apiClient APIClient
}

func NewFightIDsService(token string, apiClient APIClient) *fightIDsService {
	return &fightIDsService{
		stmt: `query($code: String!) {
			reportData {
				report(code: $code) { 
					fights(killType: Kills) {
						id
						name
						kill
					}
				}
			}
		}`,
		token: token,
		apiClient: apiClient,
	}
}

type APIClient interface {
	Post(reqBody []byte, token string) (*FFLogsAPIResponse, error)
}

func (fs *fightIDsService) Run(reportID string) (*FFLogsAPIResponse, error) {
	query := FFLogsAPIRequestQuery{
		Query: fs.stmt,
		Variables: map[string]interface{}{"code": reportID},
	}
	queryBody, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	res, err := fs.apiClient.Post(queryBody, fs.token)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type tableService struct {
	stmt string
	token string
	fightID int
	apiClient APIClient
}

func NewTableService(token string, fightID int, apiClient APIClient) *tableService {
	return &tableService{
		stmt: `query($code: String!, $fightIDs: [Int]) {
			reportData {
				report(code: $code) { 
					table(fightIDs: $fightIDs)
				}
			}
		}`,
		token: token,
		fightID: fightID,
		apiClient: apiClient,
	}
}

func (ts *tableService) Run(reportID string) (*FFLogsAPIResponse, error) {
	query := FFLogsAPIRequestQuery{
		Query: ts.stmt,
		Variables: map[string]interface{}{"code": reportID, "fightIDs": []int{ts.fightID}},
	}
	queryBody, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	res, err := ts.apiClient.Post(queryBody, ts.token)
	if err != nil {
		return nil, err
	}
	return res, nil
}
