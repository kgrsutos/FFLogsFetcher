package internal

type FFLogsAPIRequestQuery struct {
	Query string `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type FFLogsAPIResponse struct {
	Data struct {
		ReportData struct {
			Report struct {
				Code string
				Fights []struct {
					ID int
					Name string
					Kill bool
				}
				Table struct {
					Data struct {
						DeathEvents []DeathEvent
					}
				}
			}
		}
	}
}

type DeathEvent struct {
	Name string
	Ability struct {
		Name string
	}
}

type FileOutput struct {
	ReportName string
	DeathEvent []DeathEvent
}
