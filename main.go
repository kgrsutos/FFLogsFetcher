package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type QueryConfig struct {
	ReportIDs []string `yaml:"reportIDs"`
}

// CLIContext はCLIオプションを格納します。
type CLIContext struct {
	Query string `arg:"" optional:"" help:"GraphQL query string."`
}

// GraphQLQuery はGraphQLのクエリを含むリクエストボディを定義します。
type GraphQLQuery struct {
	Query string `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GraphQLResponse はGraphQLのレスポンスを定義します。
type GraphQLResponse struct {
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

var fightIDQuery = `query($code: String!) {
	reportData {
		report(code: $code) { 
			fights(killType: Kills) {
				id
				name
				kill
			}
		}
	}
}`

var tableQuery = `query($code: String!, $fightIDs: [Int]) {
	reportData {
		report(code: $code) { 
			table(fightIDs: $fightIDs)
		}
	}
}`

func main() {
	// Load Env
	token := os.Getenv("FFLOGS_ACCESS_TOKEN")
	if token == "" {
		fmt.Println("FFLOGS_ACCESS_TOKEN is not set in the environment variables")
		os.Exit(1)
	}
	// Load Query Config
	config, err := LoadQueryConfig("query.yaml")
	if err != nil {
		fmt.Printf("Error loading query config: %s\n", err)
		return
	}
	deathEvents := []DeathEvent{}
	for _, reportID := range config.ReportIDs {
		// fightID Query
		query := GraphQLQuery{
			Query: fightIDQuery,
			Variables: map[string]interface{}{"code": reportID},
		}
		// JSONに変換
		queryBody, err := json.Marshal(query)
		if err != nil {
			fmt.Println("Error marshalling query:", err)
			os.Exit(1)
		}
		graphqlResponse := post(queryBody, token)
		// レスポンスから取得したデータを出力
		for _, res := range graphqlResponse.Data.ReportData.Report.Fights {
			// fightID Query
			q2 := GraphQLQuery{
				Query: tableQuery,
				Variables: map[string]interface{}{"code": reportID, "fightIDs": []int{res.ID}},
			}
			// JSONに変換
			queryBody2, err := json.Marshal(q2)
			if err != nil {
				fmt.Println("Error marshalling query:", err)
				os.Exit(1)
			}
			graphqlResponse = post(queryBody2, token)
			deathEvents = append(deathEvents, graphqlResponse.Data.ReportData.Report.Table.Data.DeathEvents...)
			fmt.Println(deathEvents)
			break
		}
		break
	}

	// ファイルを開く、存在しない場合は新規作成
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close() // 関数終了時にファイルを閉じる

	// マップの各キーと値をイテレートして、ファイルに書き込み
	for _, ev := range deathEvents {
		tmpStr := strings.Replace(ev.Ability.Name, "{", "", -1)
		tmpStr = strings.Replace(tmpStr, "}", "", -1)
		_, err := file.WriteString(fmt.Sprintf("%s\t%v\n", ev.Name, tmpStr))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func post(reqBody []byte, token string) *GraphQLResponse {
	// リクエストを作成
	req, err := http.NewRequest("POST", "https://ja.fflogs.com/api/v2/client", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	// リクエストヘッダを設定
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// HTTPクライアントを作成し、リクエストを実行
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to the server:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// レスポンスを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	// GraphQLレスポンスをデコード
	var graphqlResponse GraphQLResponse
	if err := json.Unmarshal(body, &graphqlResponse); err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		os.Exit(1)
	}
	return &graphqlResponse
}

func LoadQueryConfig(path string) (QueryConfig, error) {
	var config QueryConfig

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
