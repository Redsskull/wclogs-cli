package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/models"
	"wclogs-cli/services"
)

// ExecuteDeathAnalysis provides detailed death analysis using Events API
func ExecuteDeathAnalysis(reportCode string, fightIDStr string, playerName string, verbose bool, enhanced bool) error {
	// Resolve fight ID (handles both numbers and "last" keyword)
	fightID, err := resolveFightID(reportCode, fightIDStr, verbose)
	if err != nil {
		return fmt.Errorf("failed to resolve fight ID '%s': %w", fightIDStr, err)
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

	var startTime *float64 = nil // No pagination in initial call
	request := api.NewDeathEventsRequest(reportCode, fightID, targetPlayerID, startTime)
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
		displayPlayerDeathAnalysis(events, playerLookup, currentFight, lookupService, apiClient, reportCode, fightID, playerName, verbose, enhanced)
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
func displayPlayerDeathAnalysis(events []*models.Event, playerLookup map[int]string, fight *models.Fight, lookupService *services.LookupService, apiClient *api.Client, reportCode string, fightID int, playerName string, verbose bool, enhanced bool) {
	color.HiRed("\n💀 DETAILED DEATH ANALYSIS: %s 💀\n", color.HiYellowString(playerName))

	fightDuration := time.Duration((fight.EndTime - fight.StartTime) * int64(time.Millisecond))
	fmt.Printf("Fight: %s (Duration: %s)\n",
		color.HiYellowString(fight.Name),
		color.HiWhiteString(fightDuration.String()))

	// Find deaths for this specific player
	var playerDeaths []*models.Event
	var targetPlayerID int
	for _, event := range events {
		if event.Type == "death" && event.TargetID != nil {
			if name, exists := playerLookup[*event.TargetID]; exists && strings.EqualFold(name, playerName) {
				playerDeaths = append(playerDeaths, event)
				targetPlayerID = *event.TargetID
			}
		}
	}

	if verbose {
		fmt.Printf("🔍 Debug: Found player ID %d for '%s'\n", targetPlayerID, playerName)
		fmt.Printf("🔍 Debug: Fight start time: %d ms\n", fight.StartTime)
	}

	if len(playerDeaths) == 0 {
		color.HiGreen("🎉 %s survived the entire fight!\n", playerName)
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
		if enhanced {
			fmt.Printf("  🔬 ENHANCED ANALYSIS MODE\n")
			executeEnhancedDeathAnalysis(apiClient, reportCode, fightID, actualPlayerID, startTime, event.Timestamp, lookupService, verbose)
		} else {
			displayDamageTimeline(apiClient, reportCode, fightID, actualPlayerID, startTime, event.Timestamp, lookupService, verbose)
		}

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
		fmt.Printf("• %s died %d times - focus on mechanics and survival\n", playerName, len(playerDeaths))
	}
	if len(playerDeaths) == 1 {
		survivalPct := (playerDeaths[0].Timestamp - fightStartTime) / float64(fight.EndTime-fight.StartTime) * 100
		fmt.Printf("• %s survived %.1f%% of the fight\n", playerName, survivalPct)
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

// displayDamageTimeline shows WCL-style death timeline with damage and healing events
func displayDamageTimeline(apiClient *api.Client, reportCode string, fightID, playerID int, startTime, endTime float64, lookupService *services.LookupService, verbose bool) {
	// Use 3-second window to capture the complete death sequence like WCL CSV
	deathTime := endTime
	windowStart := deathTime - 3000 // 3 seconds before death
	windowEnd := deathTime + 100    // 0.1 seconds after death

	if verbose {
		fmt.Printf("    🔍 Debug: Querying death timeline for player %d\n", playerID)
		fmt.Printf("    🔍 Debug: Death window: %.1f to %.1f (death at %.1f)\n",
			windowStart, windowEnd, deathTime)
	}

	// Query ALL events targeting this player to build complete death timeline
	request := &api.GraphQLRequest{
		Query: `
			query DeathTimelineEvents($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
				reportData {
					report(code: $code) {
						events(
							fightIDs: [$fightID],
							targetID: $playerID,
							startTime: $startTime,
							endTime: $endTime,
							limit: 200
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
		fmt.Printf("    ❌ Failed to fetch death timeline: %v\n", err)
		return
	}

	if response.Data == nil || response.Data.ReportData == nil ||
		response.Data.ReportData.Report == nil ||
		response.Data.ReportData.Report.Events == nil {
		fmt.Printf("    📊 No events found in death timeline\n")
		return
	}

	events, err := models.ParseEventsJSON(response.Data.ReportData.Report.Events.Data)
	if err != nil {
		fmt.Printf("    ❌ Failed to parse death timeline: %v\n", err)
		return
	}

	if verbose {
		fmt.Printf("    🔍 Debug: Found %d events in death timeline\n", len(events))
	}

	if len(events) == 0 {
		fmt.Printf("    📊 No events found - likely instant death\n")
		return
	}

	// Filter and sort events by type and significance
	var deathEvents []DeathTimelineEvent
	totalDamage := 0
	totalHealing := 0

	for _, event := range events {
		timeFromDeath := (event.Timestamp - deathTime) / 1000.0

		switch event.Type {
		case "damage":
			if event.Amount != nil && *event.Amount > 0 {
				totalDamage += *event.Amount
				abilityName := "Unknown"
				if event.AbilityID != nil {
					abilityName = lookupService.GetAbilityName(*event.AbilityID)
				}
				sourceName := "Unknown"
				if event.SourceID != nil {
					sourceName = lookupService.GetActorName(*event.SourceID)
				}

				deathEvents = append(deathEvents, DeathTimelineEvent{
					Time:        timeFromDeath,
					Type:        "damage",
					Amount:      *event.Amount,
					AbilityName: abilityName,
					SourceName:  sourceName,
					EventType:   event.Type,
				})
			}
		case "heal":
			if event.Amount != nil && *event.Amount > 0 {
				totalHealing += *event.Amount
				abilityName := "Unknown"
				if event.AbilityID != nil {
					abilityName = lookupService.GetAbilityName(*event.AbilityID)
				}
				sourceName := "Unknown"
				if event.SourceID != nil {
					sourceName = lookupService.GetActorName(*event.SourceID)
				}

				deathEvents = append(deathEvents, DeathTimelineEvent{
					Time:        timeFromDeath,
					Type:        "heal",
					Amount:      *event.Amount,
					AbilityName: abilityName,
					SourceName:  sourceName,
					EventType:   event.Type,
				})
			}
		case "death":
			abilityName := "Death"
			if event.KillingAbilityGameID != nil {
				abilityName = lookupService.GetAbilityName(*event.KillingAbilityGameID)
			}
			sourceName := "Environment"
			if event.KillerID != nil {
				sourceName = lookupService.GetActorName(*event.KillerID)
			}
			deathEvents = append(deathEvents, DeathTimelineEvent{
				Time:        0.0,
				Type:        "death",
				Amount:      0,
				AbilityName: abilityName,
				SourceName:  sourceName,
				EventType:   "death",
			})
		}
	}

	// Sort events WCL-style: Death first at 0.00s, then backwards chronologically
	sort.Slice(deathEvents, func(i, j int) bool {
		// Death events always first
		if deathEvents[i].Type == "death" && deathEvents[j].Type != "death" {
			return true
		}
		if deathEvents[j].Type == "death" && deathEvents[i].Type != "death" {
			return false
		}
		// Then by time descending (closest to death first, like WCL CSV)
		return deathEvents[i].Time > deathEvents[j].Time
	})

	// Display WCL-style death timeline matching CSV format
	fmt.Printf("    💀 Death Timeline (WCL CSV format):\n")
	fmt.Printf("    ┌─────────┬──────┬─────────────────────────────┬──────────────┐\n")
	fmt.Printf("    │ Time    │ Type │ Ability → Source            │ Amount       │\n")
	fmt.Printf("    ├─────────┼──────┼─────────────────────────────┼──────────────┤\n")

	displayCount := 0
	for _, event := range deathEvents {
		// Limit to most significant events (like WCL web interface)
		if displayCount >= 15 {
			break
		}

		// Skip very small heals/damage unless close to death
		if event.Type == "damage" && event.Amount < 50000 && event.Time < -0.5 {
			continue
		}
		if event.Type == "heal" && event.Amount < 50000 && event.Time < -0.5 {
			continue
		}

		// Format time like WCL CSV
		timeStr := ""
		if event.Type == "death" {
			timeStr = "0.00s"
		} else if event.Time == 0.0 {
			timeStr = "-0.00s" // Special case for killing blow at death time
		} else if event.Time > 0 {
			timeStr = fmt.Sprintf("+%.2fs", event.Time)
		} else {
			timeStr = fmt.Sprintf("%.2fs", event.Time)
		}

		switch event.Type {
		case "damage":
			amountStr := color.HiRedString("%s", models.FormatNumber(int64(event.Amount)))
			typeStr := color.HiRedString("DMG")
			abilityStr := fmt.Sprintf("%s ← %s", event.AbilityName, event.SourceName)
			if len(abilityStr) > 29 {
				abilityStr = abilityStr[:26] + "..."
			}
			fmt.Printf("    │ %-7s │ %-4s │ %-29s │ %12s │\n",
				timeStr, typeStr, abilityStr, amountStr)
		case "heal":
			amountStr := color.HiGreenString("+%s", models.FormatNumber(int64(event.Amount)))
			typeStr := color.HiGreenString("HEAL")
			abilityStr := fmt.Sprintf("%s ← %s", event.AbilityName, event.SourceName)
			if len(abilityStr) > 29 {
				abilityStr = abilityStr[:26] + "..."
			}
			fmt.Printf("    │ %-7s │ %-4s │ %-29s │ %12s │\n",
				timeStr, typeStr, abilityStr, amountStr)
		case "death":
			typeStr := color.HiRedString("DEATH")
			deathStr := fmt.Sprintf("%s ← %s", event.AbilityName, event.SourceName)
			if len(deathStr) > 29 {
				deathStr = deathStr[:26] + "..."
			}
			fmt.Printf("    │ %-7s │ %-4s │ %-29s │ %12s │\n",
				timeStr, typeStr, deathStr, "")
		}
		displayCount++
	}

	fmt.Printf("    └─────────┴──────┴─────────────────────────────┴──────────────┘\n")

	// Summary statistics matching WCL format
	damageStr := color.HiRedString("%s", models.FormatNumber(int64(totalDamage)))
	healingStr := color.HiGreenString("%s", models.FormatNumber(int64(totalHealing)))
	fmt.Printf("    📊 Timeline Summary: %s damage taken, %s healing received\n", damageStr, healingStr)

	// Calculate death window duration (from first major damage to death)
	deathWindowDuration := 0.0
	if len(deathEvents) > 1 {
		for i := len(deathEvents) - 1; i >= 0; i-- {
			if deathEvents[i].Type == "damage" && deathEvents[i].Amount > 500000 {
				deathWindowDuration = -deathEvents[i].Time
				break
			}
		}
	}

	if deathWindowDuration > 0 {
		fmt.Printf("    ⏱️  Death window: %.1fs (from first major damage to death)\n", deathWindowDuration)
	}
}

// DeathTimelineEvent represents an event in the death timeline
type DeathTimelineEvent struct {
	Time        float64
	Type        string
	Amount      int
	AbilityName string
	SourceName  string
	EventType   string
}
