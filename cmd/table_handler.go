package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/display"
	"wclogs-cli/models"
	"wclogs-cli/output"
	"wclogs-cli/services"
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

// executeDeathAnalysis provides detailed death analysis using Events API
func executeDeathAnalysis(reportCode string, fightIDStr string, playerName string, verbose bool) error {
	fightID, err := strconv.Atoi(fightIDStr)
	if err != nil {
		return fmt.Errorf("fight-id must be a number, got: %s", fightIDStr)
	}

	if verbose {
		color.HiBlue("💀 Starting comprehensive death analysis for report %s, fight %d", reportCode, fightID)
	}

	// Setup API client
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	authClient := auth.NewClient(cfg.ClientID, cfg.ClientSecret)
	apiClient := api.NewClient(authClient)

	// Create lookup service for ability and actor names
	lookupService := services.NewLookupService(apiClient)

	// Get fight information first to calculate correct survival times
	if verbose {
		color.HiBlue("⚔️  Fetching fight information...")
	}

	fightRequest := api.NewFightInfoRequest(reportCode)
	fightResponse, err := apiClient.Query(fightRequest.Query, fightRequest.Variables)
	if err != nil {
		return fmt.Errorf("failed to fetch fight data: %w", err)
	}

	if fightResponse.Data == nil || fightResponse.Data.ReportData == nil ||
		fightResponse.Data.ReportData.Report == nil {
		return fmt.Errorf("no fight data found")
	}

	var currentFight *models.Fight
	for _, fight := range fightResponse.Data.ReportData.Report.Fights {
		if fight.ID == fightID {
			currentFight = &fight
			break
		}
	}

	if currentFight == nil {
		return fmt.Errorf("fight %d not found in report", fightID)
	}

	if verbose {
		fightDuration := time.Duration((currentFight.EndTime - currentFight.StartTime) * int64(time.Millisecond))
		color.HiGreen("✅ Fight found: %s (Duration: %s, Kill: %t)",
			currentFight.Name, fightDuration.String(), currentFight.Kill)
	}

	// Load all actors (players, NPCs, pets) for name lookups
	if verbose {
		color.HiBlue("👥 Loading actors and game data...")
	}

	err = lookupService.LoadActorsFromReport(reportCode)
	if err != nil {
		return fmt.Errorf("failed to load actors: %w", err)
	}

	playerLookup := lookupService.GetPlayerLookup()

	// Get death events
	if verbose {
		color.HiBlue("💀 Fetching death events...")
	}

	var targetPlayerID *int
	if playerName != "" {
		// Find the specific player in the lookup
		found := false
		for id, name := range playerLookup {
			if strings.EqualFold(name, playerName) {
				targetPlayerID = &id
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("player '%s' not found", playerName)
		}
	}

	request := api.NewDeathEventsRequest(reportCode, fightID, targetPlayerID)
	response, err := apiClient.Query(request.Query, request.Variables)
	if err != nil {
		return fmt.Errorf("failed to fetch death events: %w", err)
	}

	if response.Data == nil || response.Data.ReportData == nil ||
		response.Data.ReportData.Report == nil ||
		response.Data.ReportData.Report.Events == nil {
		color.HiYellow("⚠️  No death events found - everyone survived! 🎉")
		return nil
	}

	// Parse the death events JSON
	events, err := models.ParseEventsJSON(response.Data.ReportData.Report.Events.Data)
	if err != nil {
		return fmt.Errorf("failed to parse death events: %w", err)
	}

	if len(events) == 0 {
		color.HiGreen("🎉 No deaths in this fight - perfect execution!")
		return nil
	}

	// Preload ability names for all death events to reduce API calls
	var abilityIDs []int
	for _, event := range events {
		if event.Type == "death" {
			if event.KillingAbilityGameID != nil {
				abilityIDs = append(abilityIDs, *event.KillingAbilityGameID)
			}
		}
	}
	if len(abilityIDs) > 0 {
		if verbose {
			color.HiBlue("🔍 Loading ability names...")
		}
		lookupService.PreloadAbilities(abilityIDs)
	}

	// Display death analysis - summary by default, detailed with flags
	if playerName != "" {
		// Single player detailed analysis
		displayPlayerDeathAnalysis(events, playerLookup, currentFight, lookupService, apiClient, reportCode, fightID, playerName, verbose)
	} else {
		// Fight summary for all deaths
		displayDeathSummary(events, playerLookup, currentFight, lookupService, verbose)
	}

	return nil
}

// displayDeathSummary shows a concise overview of all deaths in the fight
func displayDeathSummary(events []*models.Event, playerLookup map[int]string, fight *models.Fight, lookupService *services.LookupService, verbose bool) {
	color.HiRed("\n💀 DEATH ANALYSIS SUMMARY 💀\n")

	fightDuration := time.Duration((fight.EndTime - fight.StartTime) * int64(time.Millisecond))
	fmt.Printf("Fight: %s (Duration: %s)\n",
		color.HiYellowString(fight.Name),
		color.HiWhiteString(fightDuration.String()))

	result := color.HiGreenString("SUCCESS ✅")
	if !fight.Kill {
		result = color.HiRedString("WIPE ❌") + fmt.Sprintf(" (%.1f%%)", fight.FightPercentage)
	}
	fmt.Printf("Result: %s\n", result)
	fmt.Printf("Deaths: %s\n\n", color.HiRedString("%d", len(events)))

	if len(events) == 0 {
		color.HiGreen("🎉 Perfect execution - no deaths!\n")
		return
	}

	// Group deaths by timing and ability
	deathsByTime := make(map[string][]string)
	abilityCount := make(map[int]int)
	fightStartTime := float64(fight.StartTime)

	for _, event := range events {
		if event.Type != "death" {
			continue
		}

		playerName := "Unknown"
		if event.TargetID != nil {
			if name, exists := playerLookup[*event.TargetID]; exists {
				playerName = name
			} else {
				playerName = fmt.Sprintf("Player-%d", *event.TargetID)
			}
		}

		survivalTime := time.Duration((event.Timestamp - fightStartTime) * float64(time.Millisecond))
		timeKey := fmt.Sprintf("%.0fs", survivalTime.Seconds())
		deathsByTime[timeKey] = append(deathsByTime[timeKey], playerName)

		if event.KillingAbilityGameID != nil {
			abilityCount[*event.KillingAbilityGameID]++
		}
	}

	// Display death timeline
	fmt.Printf("📅 DEATH TIMELINE:\n")
	for timeKey, players := range deathsByTime {
		if len(players) == 1 {
			fmt.Printf("  • %s: %s\n",
				color.HiWhiteString(timeKey),
				color.HiYellowString(players[0]))
		} else {
			fmt.Printf("  • %s: %s (%d players)\n",
				color.HiWhiteString(timeKey),
				color.HiYellowString(strings.Join(players, ", ")),
				len(players))
		}
	}

	// Display top killing abilities
	if len(abilityCount) > 0 {
		fmt.Printf("\n⚔️  TOP KILLING ABILITIES:\n")
		type abilityDeath struct {
			id    int
			count int
		}
		var sortedAbilities []abilityDeath
		for id, count := range abilityCount {
			sortedAbilities = append(sortedAbilities, abilityDeath{id, count})
		}
		// Simple sort by count (descending)
		for i := 0; i < len(sortedAbilities)-1; i++ {
			for j := i + 1; j < len(sortedAbilities); j++ {
				if sortedAbilities[j].count > sortedAbilities[i].count {
					sortedAbilities[i], sortedAbilities[j] = sortedAbilities[j], sortedAbilities[i]
				}
			}
		}

		for _, ability := range sortedAbilities {
			abilityName := lookupService.GetAbilityName(ability.id)
			fmt.Printf("  • %s: %s\n",
				color.HiYellowString(abilityName),
				color.HiRedString("%d deaths", ability.count))
		}
	}

	color.HiCyan("\n💡 TIP: Use --player \"PlayerName\" for detailed death analysis of a specific player")
	fmt.Println()
}

// displayPlayerDeathAnalysis shows detailed analysis for a specific player
func displayPlayerDeathAnalysis(events []*models.Event, playerLookup map[int]string, fight *models.Fight, lookupService *services.LookupService, apiClient *api.Client, reportCode string, fightID int, targetPlayerName string, verbose bool) {
	color.HiRed("\n💀 DETAILED DEATH ANALYSIS: %s 💀\n", color.HiYellowString(targetPlayerName))

	fightDuration := time.Duration((fight.EndTime - fight.StartTime) * int64(time.Millisecond))
	fmt.Printf("Fight: %s (Duration: %s)\n",
		color.HiYellowString(fight.Name),
		color.HiWhiteString(fightDuration.String()))

	// Find deaths for this specific player
	var playerDeaths []*models.Event
	var targetPlayerID int
	for _, event := range events {
		if event.Type == "death" && event.TargetID != nil {
			if name, exists := playerLookup[*event.TargetID]; exists && strings.EqualFold(name, targetPlayerName) {
				playerDeaths = append(playerDeaths, event)
				targetPlayerID = *event.TargetID
			}
		}
	}

	if verbose {
		fmt.Printf("🔍 Debug: Found player ID %d for '%s'\n", targetPlayerID, targetPlayerName)
		fmt.Printf("🔍 Debug: Fight start time: %d ms\n", fight.StartTime)
	}

	if len(playerDeaths) == 0 {
		color.HiGreen("🎉 %s survived the entire fight!\n", targetPlayerName)
		return
	}

	fmt.Printf("Deaths: %s\n\n", color.HiRedString("%d", len(playerDeaths)))

	fightStartTime := float64(fight.StartTime)

	for i, event := range playerDeaths {
		survivalTime := time.Duration((event.Timestamp - fightStartTime) * float64(time.Millisecond))

		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("%s Death #%d\n", color.HiRedString("💀"), i+1)
		fmt.Printf("  ⏱️  Survival Time: %s\n", color.HiWhiteString(survivalTime.String()))

		if verbose && event.TargetID != nil {
			fmt.Printf("  🔍 Debug: Death event targetID: %d\n", *event.TargetID)
		}

		// Get readable ability and source names
		abilityName, sourceName := lookupService.FormatKillingInfo(event.KillerID, event.KillingAbilityGameID)

		fmt.Printf("  ⚔️  Killed by: %s from %s\n",
			color.HiRedString(abilityName),
			color.HiMagentaString(sourceName))

		// Detailed timeline analysis - use 5-second focused window
		timeWindow := 5.0 * 1000 // 5 seconds in milliseconds
		startTime := event.Timestamp - timeWindow
		if startTime < fightStartTime {
			startTime = fightStartTime
		}

		if verbose {
			fmt.Printf("  🕐 Death at: %.1fs into fight\n", (event.Timestamp-fightStartTime)/1000.0)
			fmt.Printf("  📊 Analyzing 5s around death...\n")
		}

		// Show damage timeline leading to death - use the actual targetID from this specific death event
		actualPlayerID := targetPlayerID
		if event.TargetID != nil {
			actualPlayerID = *event.TargetID
		}

		fmt.Printf("  📈 Events Around Death:\n")
		displayDamageTimeline(apiClient, reportCode, fightID, actualPlayerID, startTime, event.Timestamp, lookupService, verbose)

		// Get healing summary (not full timeline)
		fmt.Printf("  💚 Healing Analysis:\n")
		healingTotal := getHealingSummary(apiClient, reportCode, fightID, actualPlayerID, startTime, event.Timestamp)
		if healingTotal > 0 {
			fmt.Printf("    • Total healing: %s (healers tried hard!)\n",
				color.HiGreenString("%d", healingTotal))
		} else {
			fmt.Printf("    • %s\n", color.HiYellowString("No significant healing - may have been unavoidable"))
		}

		// Get defensive abilities summary
		fmt.Printf("  🛡️  Defensive Analysis:\n")
		defensiveCount := getDefensiveSummary(apiClient, reportCode, fightID, actualPlayerID, startTime, event.Timestamp)
		if defensiveCount > 0 {
			fmt.Printf("    • Used %s defensive abilities\n", color.HiBlueString("%d", defensiveCount))
		} else {
			fmt.Printf("    • %s\n", color.HiYellowString("No defensives used - could have helped survive"))
		}

		fmt.Println()
	}

	// Player-specific insights
	color.HiBlue("📊 INSIGHTS:")
	if len(playerDeaths) > 1 {
		fmt.Printf("• %s died %d times - focus on mechanics and survival\n", targetPlayerName, len(playerDeaths))
	}
	if len(playerDeaths) == 1 {
		survivalPct := (playerDeaths[0].Timestamp - fightStartTime) / float64(fight.EndTime-fight.StartTime) * 100
		fmt.Printf("• %s survived %.1f%% of the fight\n", targetPlayerName, survivalPct)
	}
}

// getHealingSummary returns total healing received in the time window
func getHealingSummary(apiClient *api.Client, reportCode string, fightID, playerID int, startTime, endTime float64) int {
	request := api.NewHealingReceivedRequest(reportCode, fightID, playerID, startTime, endTime)
	response, err := apiClient.Query(request.Query, request.Variables)
	if err != nil {
		return 0
	}

	if response.Data == nil || response.Data.ReportData == nil ||
		response.Data.ReportData.Report == nil ||
		response.Data.ReportData.Report.Events == nil {
		return 0
	}

	events, err := models.ParseEventsJSON(response.Data.ReportData.Report.Events.Data)
	if err != nil {
		return 0
	}

	totalHealing := 0
	for _, event := range events {
		if event.Type == "heal" && event.Amount != nil {
			totalHealing += *event.Amount
		}
	}
	return totalHealing
}

// getDefensiveSummary returns count of defensive abilities used in the time window
func getDefensiveSummary(apiClient *api.Client, reportCode string, fightID, playerID int, startTime, endTime float64) int {
	request := api.NewDefensiveAbilitiesRequest(reportCode, fightID, playerID, startTime, endTime)
	response, err := apiClient.Query(request.Query, request.Variables)
	if err != nil {
		return 0
	}

	if response.Data == nil || response.Data.ReportData == nil ||
		response.Data.ReportData.Report == nil ||
		response.Data.ReportData.Report.Events == nil {
		return 0
	}

	events, err := models.ParseEventsJSON(response.Data.ReportData.Report.Events.Data)
	if err != nil {
		return 0
	}

	defensiveCount := 0
	for _, event := range events {
		if event.Type == "cast" || event.Type == "begincast" {
			defensiveCount++
		}
	}
	return defensiveCount
}

// displayDamageTimeline shows all events around death to find damage sources
func displayDamageTimeline(apiClient *api.Client, reportCode string, fightID, playerID int, startTime, endTime float64, lookupService *services.LookupService, verbose bool) {
	// Use shorter 5-second window around death for all events
	deathTime := endTime
	windowStart := deathTime - 5000 // 5 seconds before death
	windowEnd := deathTime + 1000   // 1 second after death

	if verbose {
		fmt.Printf("    🔍 Debug: Querying ALL events for player %d\n", playerID)
		fmt.Printf("    🔍 Debug: 5-second window: %.1fs to %.1fs\n",
			windowStart/1000.0, windowEnd/1000.0)
	}

	// Query all events targeting this player around death time
	request := &api.GraphQLRequest{
		Query: `
			query AllEventsAroundDeath($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
				reportData {
					report(code: $code) {
						events(
							fightIDs: [$fightID],
							targetID: $playerID,
							startTime: $startTime,
							endTime: $endTime,
							limit: 100
						) {
							data
						}
					}
				}
			}`,
		Variables: map[string]any{
			"code":      reportCode,
			"fightID":   fightID,
			"playerID":  playerID,
			"startTime": windowStart,
			"endTime":   windowEnd,
		},
	}
	response, err := apiClient.Query(request.Query, request.Variables)
	if err != nil {
		fmt.Printf("    ❌ Failed to fetch damage data: %v\n", err)
		return
	}

	// Save raw response to debug file if needed
	if verbose {
		filename := fmt.Sprintf("events_debug_%s_%d_%d.json", reportCode, fightID, playerID)
		if jsonData, err := json.MarshalIndent(response, "", "  "); err == nil {
			if err := os.WriteFile(filename, jsonData, 0644); err == nil {
				fmt.Printf("    🔍 Debug: Saved raw response to %s\n", filename)
			}
		}
	}

	if response.Data == nil || response.Data.ReportData == nil ||
		response.Data.ReportData.Report == nil ||
		response.Data.ReportData.Report.Events == nil {
		fmt.Printf("    📊 No damage events found\n")
		return
	}

	events, err := models.ParseEventsJSON(response.Data.ReportData.Report.Events.Data)
	if err != nil {
		fmt.Printf("    ❌ Failed to parse damage events: %v\n", err)
		return
	}

	if verbose {
		fmt.Printf("    🔍 Debug: Found %d total events\n", len(events))
		if len(events) > 0 {
			fmt.Printf("    🔍 Debug: First event type: %s, timestamp: %.0f\n", events[0].Type, events[0].Timestamp)
			if len(events) > 1 {
				fmt.Printf("    🔍 Debug: Last event type: %s, timestamp: %.0f\n", events[len(events)-1].Type, events[len(events)-1].Timestamp)
			}
		}
	}

	if len(events) == 0 {
		fmt.Printf("    📊 No events found in 5-second death window\n")
		if verbose {
			fmt.Printf("    💡 This might indicate instant-death mechanics\n")
		}
		return
	}

	// Show all events around death time for context
	fmt.Printf("    📊 Events in 5-second death window:\n")

	totalDamage := 0
	damageCount := 0

	for _, event := range events {
		timeFromDeath := (deathTime - event.Timestamp) / 1000.0
		timeLabel := fmt.Sprintf("%.1fs", timeFromDeath)
		if timeFromDeath < 0 {
			timeLabel = fmt.Sprintf("+%.1fs", -timeFromDeath)
		}

		switch event.Type {
		case "damage":
			if event.Amount != nil {
				damageCount++
				totalDamage += *event.Amount

				abilityName := "Unknown"
				if event.AbilityID != nil {
					abilityName = lookupService.GetAbilityName(*event.AbilityID)
				}

				sourceName := "Unknown"
				if event.SourceID != nil {
					sourceName = lookupService.GetActorName(*event.SourceID)
				}

				fmt.Printf("    • -%s: %s damage from %s (%s)\n",
					timeLabel,
					color.HiRedString("%d", *event.Amount),
					color.HiMagentaString(sourceName),
					color.HiYellowString(abilityName))
			}
		case "heal":
			if verbose && event.Amount != nil {
				fmt.Printf("    • -%s: %s healing\n",
					timeLabel,
					color.HiGreenString("%d", *event.Amount))
			}
		case "death":
			fmt.Printf("    • -%s: ⚰️  DEATH EVENT\n", timeLabel)
		default:
			if verbose {
				fmt.Printf("    • -%s: %s event\n", timeLabel, event.Type)
			}
		}
	}

	if damageCount == 0 {
		fmt.Printf("    💡 No damage events - likely environmental/scripted death\n")
	} else {
		fmt.Printf("    📊 Total damage in window: %s (%d events)\n",
			color.HiRedString("%d", totalDamage), damageCount)
	}

}
