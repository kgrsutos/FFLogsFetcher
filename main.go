package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/alecthomas/kong"
)

// CLIContext はCLIオプションを格納します。
type CLIContext struct {
	Query string `arg:"" optional:"" help:"GraphQL query string."`
}

// GraphQLQuery はGraphQLのクエリを含むリクエストボディを定義します。
type GraphQLQuery struct {
	Query string `json:"query"`
}

// GraphQLResponse はGraphQLのレスポンスを定義します。
type GraphQLResponse struct {
	Data struct {
		ReportData struct {
			Report struct {
				Code string `json:"code"`
			} `json:"report"`
		} `json:"reportData"`
	} `json:"data"`
}

func main() {
	// 環境変数からトークンを取得
	token := os.Getenv("FFLOGS_ACCESS_TOKEN")
	if token == "" {
		fmt.Println("FFLOGS_ACCESS_TOKEN is not set in the environment variables")
		os.Exit(1)
	}

	// Kongを用いたCLIオプションのパース
	var cliContext CLIContext
	kong.Parse(&cliContext)

	// クエリが提供されない場合はデフォルトのクエリを使用
	if cliContext.Query == "" {
		cliContext.Query = `{reportData { report(code: "Ddk8NJLzmqgQXRnB") { code } } }`
	}

	// GraphQLクエリを作成
	query := GraphQLQuery{
		Query: cliContext.Query,
	}

	// JSONに変換
	queryBody, err := json.Marshal(query)
	if err != nil {
		fmt.Println("Error marshalling query:", err)
		os.Exit(1)
	}

	// リクエストを作成
	req, err := http.NewRequest("POST", "https://ja.fflogs.com/api/v2/client", bytes.NewBuffer(queryBody))
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

	// レスポンスから取得したデータを出力
	fmt.Println("Code:", graphqlResponse.Data.ReportData.Report.Code)
}
