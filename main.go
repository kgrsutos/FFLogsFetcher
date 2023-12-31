package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/shigeo-kogure/FFLogsFetcher/internal"
)

func main() {
	// Setup Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	// Load Env
	token := os.Getenv("FFLOGS_ACCESS_TOKEN")
	if token == "" {
		log.Error().Msg("FFLOGS_ACCESS_TOKEN is not set in the environment variables")
		os.Exit(1)
	}
	// Load GraphQL Config
	config, err := internal.LoadRequestConfig("query.yaml")
	if err != nil {
		log.Error().Stack().Err(err).Msg("Error loading query config")
		return
	}
	fflogsAPICl := internal.FFLogsAPIClient{Client: &http.Client{}}
	deathEvents := []internal.DeathEventOutput{}
	for _, reportID := range config.ReportIDs {
		// Get Fight IDs
		fightIDsService := internal.NewFightIDsService(token, &fflogsAPICl)
		fightIDsResponse, err := fightIDsService.Run(reportID)
		if err != nil {
			log.Error().Stack().Err(err).Msg("failed to get fight ids")
			os.Exit(1)
		}
		for _, fightIDRes := range fightIDsResponse.Data.ReportData.Report.Fights {
			// Get Table of Fight ID
			tableService := internal.NewTableService(token, fightIDRes.ID, &fflogsAPICl)
			tableResponse, err := tableService.Run(reportID)
			if err != nil {
				log.Error().Stack().Err(err).Msg("failed to get table info with fight id")
				os.Exit(1)
			}
			for _, ev := range tableResponse.Data.ReportData.Report.Table.Data.DeathEvents {
				deathEvents = append(deathEvents, internal.DeathEventOutput{
					ReportName: fightIDRes.Name,
					PlayerName: ev.Name,
					AbilityName: ev.Ability.Name,
				})
			}
		}
	}
	// Output file
	if err := internal.OutputDeathEvent("output.txt", deathEvents); err != nil {
		log.Error().Stack().Err(err).Msg("failed to output file")
		os.Exit(1)
	}
}
