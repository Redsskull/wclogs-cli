package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/display"
	"wclogs-cli/models"
	"wclogs-cli/output"
)

// executeTableCommand is the shared handler with player filtering support
func executeTableCommand(tableType string, reportCode string, fightID int, topN int, noColor bool, verbose bool, outputPath string, playerName string) error {
	// Get table info from types.go
	info, exists := tableTypes[tableType]
	if !exists {
		return fmt.Errorf("unsupported table type: %s", tableType)
	}

	if verbose {
		if playerName != "" {
			color.HiBlue("🔍 Fetching %s for player '%s' in report %s, fight %d", info.Description, playerName, reportCode, fightID)
		} else {
			color.HiBlue("🔍 Fetching %s for report %s, fight %d", info.Description, reportCode, fightID)
		}
	}

	// Auth logic
	if verbose {
		color.HiBlue("🔐 Loading configuration...")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// API client setup
	if verbose {
		color.HiBlue("🔐 Setting up authentication...")
	}
	authClient := auth.NewClient(cfg.ClientID, cfg.ClientSecret)

	if verbose {
		color.HiBlue("📡 Creating API client...")
	}
	apiClient := api.NewClient(authClient)

	// Validation
	if verbose {
		color.HiBlue("✅ Validating parameters...")
	}
	if err := api.ValidateQueryVariables(reportCode, fightID); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// If player filtering is requested, validate the player name exists
	if playerName != "" {
		if verbose {
			color.HiBlue("👥 Validating player name...")
		}

		// Get masterData to validate player name exists
		masterRequest := api.NewMasterDataRequest(reportCode)
		masterResponse, err := apiClient.Query(masterRequest.Query, masterRequest.Variables)
		if err != nil {
			return fmt.Errorf("failed to fetch player data: %w", err)
		}

		if masterResponse.Data == nil || masterResponse.Data.ReportData == nil ||
			masterResponse.Data.ReportData.Report == nil ||
			masterResponse.Data.ReportData.Report.MasterData == nil {
			return fmt.Errorf("no player data found for report: %s", reportCode)
		}

		playerLookup := models.NewPlayerLookup(masterResponse.Data.ReportData.Report.MasterData.Actors)
		if err := playerLookup.ValidatePlayerName(playerName); err != nil {
			return fmt.Errorf("player validation failed: %w", err)
		}

		if verbose {
			color.HiGreen("✅ Player '%s' found in report", playerName)
		}
	}

	// Query execution
	if verbose {
		color.HiBlue("🚀 Executing GraphQL query for %s...", info.Description)
	}

	// Use our generic request builder
	request := api.NewTableRequest(reportCode, fightID, info.DataType)
	response, err := apiClient.Query(request.Query, request.Variables)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Pass the APIs response
	if response.Data == nil || response.Data.ReportData == nil || response.Data.ReportData.Report == nil {
		return fmt.Errorf("no report data found for code: %s", reportCode)
	}

	rawTable := response.Data.ReportData.Report.Table
	if len(rawTable) == 0 {
		return fmt.Errorf("no %s data found for fight %d in report %s", info.Description, fightID, reportCode)
	}

	if verbose {
		color.HiBlue("✅ Got table data! Parsing...")
	}

	// Process the data
	tableData, err := models.ParseTableData(rawTable)
	if err != nil {
		return fmt.Errorf("failed to parse table data: %w", err)
	}

	players := models.GetPlayersFromTable(tableData)

	if verbose {
		color.HiBlue("📊 Found %d players in the table", len(players))
	}

	// Apply player filtering if requested
	if playerName != "" {
		filteredPlayers := filterPlayersByName(players, playerName)
		if len(filteredPlayers) == 0 {
			return fmt.Errorf("player '%s' not found in %s data for fight %d", playerName, info.Description, fightID)
		}
		players = filteredPlayers

		if verbose {
			color.HiGreen("🎯 Filtered to %d player(s) matching '%s'", len(players), playerName)
		}
	}

	// Choose the output method
	if outputPath != "" {
		// File output - use new output system
		var total int64
		for _, player := range players {
			total += int64(player.Total)
		}

		outputData := &output.OutputData{
			Players:    players,
			ReportCode: reportCode,
			FightID:    fightID,
			Title:      info.Title,
			Total:      total,
		}

		return output.HandleOutput(outputData, outputPath, topN, noColor, verbose)
	} else {
		// Terminal output - use existing beautiful display
		options := display.DefaultTableOptions()
		options.TopN = topN
		options.UseColors = !noColor

		// Display with custom title for this data type
		if playerName != "" {
			fmt.Printf("\n%s %s for %s %s\n", info.Emoji, info.Title, color.HiYellowString(playerName), info.Emoji)
		} else {
			fmt.Printf("\n%s %s %s\n", info.Emoji, info.Title, info.Emoji)
		}

		display.DisplayTable(players, tableType, options)

		return nil
	}
}

// filterPlayersByName filters players list to match the specified player name
func filterPlayersByName(players []*models.Player, targetName string) []*models.Player {
	var filtered []*models.Player

	for _, player := range players {
		// Case-insensitive name matching
		if strings.EqualFold(player.Name, targetName) {
			filtered = append(filtered, player)
		}
	}

	return filtered
}
