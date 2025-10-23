package cmd

import (
	"fmt"

	"github.com/fatih/color"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/models"
	"wclogs-cli/output"
)

// executePlayersCommand handles the players command
func executePlayersCommand(reportCode string, verbose bool, outputPath string) error {
	if verbose {
		color.HiBlue("🔍 Fetching player list for report %s", reportCode)
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
	if reportCode == "" {
		return fmt.Errorf("report code cannot be empty")
	}

	if len(reportCode) < 6 {
		return fmt.Errorf("report code '%s' is too short (must be at least 6 characters)", reportCode)
	}

	// Query execution
	if verbose {
		color.HiBlue("🚀 Executing masterData GraphQL query...")
	}

	request := api.NewMasterDataRequest(reportCode)
	response, err := apiClient.Query(request.Query, request.Variables)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Validate response structure
	if response.Data == nil || response.Data.ReportData == nil || response.Data.ReportData.Report == nil {
		return fmt.Errorf("no report data found for code: %s", reportCode)
	}

	if response.Data.ReportData.Report.MasterData == nil {
		return fmt.Errorf("no master data found for report %s", reportCode)
	}

	masterData := response.Data.ReportData.Report.MasterData
	if len(masterData.Actors) == 0 {
		return fmt.Errorf("no players found in report %s", reportCode)
	}

	if verbose {
		color.HiBlue("✅ Found %d players in the report!", len(masterData.Actors))
	}

	// Create player lookup
	playerLookup := models.NewPlayerLookup(masterData.Actors)

	// Handle output
	if outputPath != "" {
		// File output
		return outputPlayersToFile(playerLookup, reportCode, outputPath, verbose)
	} else {
		// Terminal output
		displayPlayersInTerminal(playerLookup, reportCode)
		return nil
	}
}

// displayPlayersInTerminal shows the player list in a beautiful terminal format
func displayPlayersInTerminal(playerLookup *models.PlayerLookup, reportCode string) {
	players := playerLookup.GetAllPlayers()

	// Header
	fmt.Printf("\n👥 %s 👥\n", color.HiCyanString("PLAYERS IN REPORT %s", reportCode))
	fmt.Printf("%s\n\n", color.HiBlackString("Found %d players:", len(players)))

	// Table headers
	color.HiWhite("%-3s %-20s %-12s %-20s", "#", "NAME", "CLASS", "SERVER")
	color.HiBlack("%-3s %-20s %-12s %-20s", "---", "--------------------", "------------", "--------------------")

	// Player list with class colors
	for i, player := range players {
		classColor := getClassColor(player.Class)
		fmt.Printf("%-3d %-20s %s %-20s\n",
			i+1,
			player.Name,
			classColor.Sprintf("%-12s", player.Class),
			player.Server)
	}

	fmt.Printf("\n%s\n", color.HiGreenString("✅ Use these exact names with --player flag"))
	fmt.Printf("%s\n", color.HiYellowString("Example: wclogs damage %s 5 --player \"%s\"", reportCode, players[0].Name))
}

// outputPlayersToFile saves the player list to a file
func outputPlayersToFile(playerLookup *models.PlayerLookup, reportCode string, outputPath string, verbose bool) error {
	players := playerLookup.GetAllPlayers()

	// Create output data structure similar to table data
	outputData := &output.PlayersOutputData{
		Players:    players,
		ReportCode: reportCode,
		Count:      len(players),
	}

	return output.HandlePlayersOutput(outputData, outputPath, verbose)
}

// getClassColor returns the appropriate color function for each class
func getClassColor(class string) *color.Color {
	switch class {
	case "Death Knight":
		return color.New(color.FgHiRed)
	case "Demon Hunter":
		return color.New(color.FgHiMagenta)
	case "Druid":
		return color.New(color.FgYellow)
	case "Hunter":
		return color.New(color.FgGreen)
	case "Mage":
		return color.New(color.FgCyan)
	case "Monk":
		return color.New(color.FgHiGreen)
	case "Paladin":
		return color.New(color.FgHiYellow)
	case "Priest":
		return color.New(color.FgWhite)
	case "Rogue":
		return color.New(color.FgYellow)
	case "Shaman":
		return color.New(color.FgBlue)
	case "Warlock":
		return color.New(color.FgMagenta)
	case "Warrior":
		return color.New(color.FgRed)
	case "Evoker":
		return color.New(color.FgHiCyan)
	default:
		return color.New(color.FgWhite)
	}
}
