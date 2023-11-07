package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
)

var query struct {
	ReportData struct {
		Report struct {
			Code string
		} `graphql:"report(code: $code)"`
	}
}

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("FFLOGS_ACCESS_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient("https://ja.fflogs.com/api/v2/client", httpClient)
	variables := map[string]interface{}{
		"code": "Ddk8NJLzmqgQXRnB",
	}
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(query.ReportData.Report.Code)
}
