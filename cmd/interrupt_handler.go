package cmd

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/models"
	"wclogs-cli/services"
)

// ExecuteInterruptAnalysis provides detailed interrupt analysis using Events API
func ExecuteInterruptAnalysis(reportCode string, fightIDStr string, playerName string, verbose bool) error {
	fightID, err := strconv.Atoi(fightIDStr)
	if err != nil {
		return fmt.Errorf("fight-id must be a number, got: %s", fightIDStr)
	}

	if verbose {
		color.HiBlue("🎛️ Starting comprehensive interrupt analysis for report %s, fight %d", reportCode, fightID)
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

	if verbose {
		color.HiBlue("🤖 Fetching interrupt events...")
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

	// Fetch interrupt events
	var startTime *float64 = nil // No pagination in initial call
	interruptRequest := api.NewInterruptEventsRequest(reportCode, fightID, targetPlayerID, startTime)
	interruptResponse, err := apiClient.Query(interruptRequest.Query, interruptRequest.Variables)
	if err != nil {
		return fmt.Errorf("failed to fetch interrupt events: %w", err)
	}

	if interruptResponse.Data == nil || interruptResponse.Data.ReportData == nil ||
		interruptResponse.Data.ReportData.Report == nil ||
		interruptResponse.Data.ReportData.Report.Events == nil {
		color.HiYellow("⚠️  No interrupt events found! 🤷")
		return nil
	}

	// Parse interrupt events JSON
	interruptEvents, err := models.ParseInterruptEventsJSON(interruptResponse.Data.ReportData.Report.Events.Data)
	if err != nil {
		return fmt.Errorf("failed to parse interrupt events: %w", err)
	}

	// If no interrupt events found, show message and return
	if len(interruptEvents) == 0 {
		if playerName != "" {
			color.HiYellow("🤔 Player '%s' did not perform any interrupts in this fight", playerName)
		} else {
			color.HiYellow("🤔 No interrupts occurred in this fight")
		}
		return nil
	}

	if verbose {
		color.HiBlue("✅ Found %d interrupt events", len(interruptEvents))
	}

	// Preload ability names for all interrupt events to reduce API calls
	var abilityIDs []int
	for _, event := range interruptEvents {
		if event.AbilityID != nil {
			abilityIDs = append(abilityIDs, *event.AbilityID)
		}
		if event.SourceID != nil {
			abilityIDs = append(abilityIDs, *event.SourceID)
		}
	}
	if len(abilityIDs) > 0 {
		if verbose {
			color.HiBlue("🔍 Loading ability names...")
		}
		lookupService.PreloadAbilities(abilityIDs)
	}

	// Display interrupt analysis
	if playerName != "" {
		// Single player detailed analysis
		displayPlayerInterruptAnalysis(interruptEvents, playerLookup, currentFight, lookupService, apiClient, reportCode, fightID, playerName, verbose)
	} else {
		// Fight summary for all interrupts
		displayInterruptSummary(interruptEvents, playerLookup, currentFight, lookupService, apiClient, reportCode, fightID, verbose)
	}

	return nil
}

// displayInterruptSummary shows a concise overview of all interrupts in the fight
func displayInterruptSummary(events []*models.Event, playerLookup map[int]string, fight *models.Fight, lookupService *services.LookupService, apiClient *api.Client, reportCode string, fightID int, verbose bool) {
	color.HiBlue("\n🎛️  INTERRUPT ANALYSIS SUMMARY 🎛️\n")

	fightDuration := time.Duration((fight.EndTime - fight.StartTime) * int64(time.Millisecond))
	fmt.Printf("Fight: %s (Duration: %s)\n",
		color.HiYellowString(fight.Name),
		color.HiWhiteString(fightDuration.String()))

	result := color.HiGreenString("SUCCESS ✅")
	if !fight.Kill {
		result = color.HiRedString("WIPE ❌") + fmt.Sprintf(" (%.1f%%)", fight.FightPercentage)
	}
	fmt.Printf("Result: %s\n", result)
	fmt.Printf("Total Interrupts: %s\n\n", color.HiBlueString("%d", len(events)))

	if len(events) == 0 {
		color.HiGreen("🎉 No interrupts in this fight!\n")
		return
	}

	// Group interrupts by interrupt source (players)
	interruptsByPlayer := make(map[string]int)

	for _, event := range events {
		playerName := "Unknown"
		if event.SourceID != nil {
			if name, exists := playerLookup[*event.SourceID]; exists {
				playerName = name
			} else {
				playerName = fmt.Sprintf("Player-%d", *event.SourceID)
			}
		}

		interruptsByPlayer[playerName]++
	}

	// Display top interrupters
	fmt.Printf("🏆 TOP INTERRUPTERS:\n")
	type playerInterrupt struct {
		name  string
		count int
	}
	var sortedPlayers []playerInterrupt
	for name, count := range interruptsByPlayer {
		sortedPlayers = append(sortedPlayers, playerInterrupt{name, count})
	}
	// Simple sort by count (descending)
	for i := 0; i < len(sortedPlayers)-1; i++ {
		for j := i + 1; j < len(sortedPlayers); j++ {
			if sortedPlayers[j].count > sortedPlayers[i].count {
				sortedPlayers[i], sortedPlayers[j] = sortedPlayers[j], sortedPlayers[i]
			}
		}
	}

	for _, player := range sortedPlayers {
		fmt.Printf("  • %s: %s interrupts\n",
			color.HiYellowString(player.name),
			color.HiBlueString("%d", player.count))
	}

	// Now let's add the correlation analysis to the summary as well
	fmt.Printf("\n🔄 CORRELATING INTERRUPTS WITH TARGET CASTS...\n")

	// Correlate interrupts with casts to determine what was actually interrupted
	analysis, err := CorrelateInterruptsAndCasts(apiClient, reportCode, fightID, events, verbose, fight.StartTime)
	if err != nil {
		fmt.Printf("❌ Error correlating interrupts with target casts: %v\n", err)
		fmt.Printf("📊 Summary will show interrupt abilities used instead of what was interrupted\n")

		// Show fallback information with interrupt abilities used
		interruptAbilitiesUsed := make(map[string]int)
		for _, event := range events {
			abilityName := "Unknown Ability"
			if event.AbilityID != nil {
				abilityName = lookupService.GetAbilityName(*event.AbilityID)
			}
			interruptAbilitiesUsed[abilityName]++
		}

		if len(interruptAbilitiesUsed) > 0 {
			fmt.Printf("\n🎭 INTERRUPT ABILITIES USED (without correlation):\n")
			type abilityUsed struct {
				name  string
				count int
			}
			var sortedAbilityUsed []abilityUsed
			for name, count := range interruptAbilitiesUsed {
				sortedAbilityUsed = append(sortedAbilityUsed, abilityUsed{name, count})
			}
			// Simple sort by count (descending)
			for i := 0; i < len(sortedAbilityUsed)-1; i++ {
				for j := i + 1; j < len(sortedAbilityUsed); j++ {
					if sortedAbilityUsed[j].count > sortedAbilityUsed[i].count {
						sortedAbilityUsed[i], sortedAbilityUsed[j] = sortedAbilityUsed[j], sortedAbilityUsed[i]
					}
				}
			}

			for _, ability := range sortedAbilityUsed {
				fmt.Printf("  • %s: %s times\n",
					color.HiYellowString(ability.name),
					color.HiBlueString("%d", ability.count))
			}
		}
	} else if len(analysis) > 0 {
		fmt.Printf("\n🏆 WHAT WAS ACTUALLY INTERRUPTED:\n")

		// Calculate totals for percentages
		totalInterrupted := 0
		totalCompleted := 0
		for _, details := range analysis {
			totalInterrupted += details.Stopped
			totalCompleted += details.Missed
		}
		totalCasts := totalInterrupted + totalCompleted

		for abilityName, details := range analysis {
			stoppedPct := float64(details.Stopped) / float64(details.TotalCasts) * 100
			missedPct := float64(details.Missed) / float64(details.TotalCasts) * 100

			fmt.Printf("\n=== %s ===\n", color.HiYellowString(abilityName))
			fmt.Printf("Total Casts: %d\n", details.TotalCasts)
			fmt.Printf("Stopped: %.1f%% (%d)\n", stoppedPct, details.Stopped)
			fmt.Printf("Completed: %.1f%% (%d)\n", missedPct, details.Missed)

			// Show who interrupted the casts
			if len(details.InterruptedBy) > 0 {
				fmt.Println("\nInterrupted By:")
				for interrupter, count := range details.InterruptedBy {
					interrupterPct := float64(count) / float64(details.TotalCasts) * 100
					fmt.Printf("  %s: %d (%.1f%%)\n", interrupter, count, interrupterPct)
				}
			}
		}

		if totalCasts > 0 {
			overallPct := float64(totalInterrupted) / float64(totalCasts) * 100
			fmt.Printf("\n📊 OVERALL SUMMARY:\n")
			fmt.Printf("Total Interrupted: %d\n", totalInterrupted)
			fmt.Printf("Total Completed: %d\n", totalCompleted)
			fmt.Printf("Overall Interrupt Effectiveness: %.1f%%\n", overallPct)
		}
	} else {
		fmt.Printf("📊 No target cast correlations found - targets may not have cast interruptible abilities\n")

		// Show fallback information with interrupt abilities used
		interruptAbilitiesUsed := make(map[string]int)
		for _, event := range events {
			abilityName := "Unknown Ability"
			if event.AbilityID != nil {
				abilityName = lookupService.GetAbilityName(*event.AbilityID)
			}
			interruptAbilitiesUsed[abilityName]++
		}

		if len(interruptAbilitiesUsed) > 0 {
			fmt.Printf("\n🎭 INTERRUPT ABILITIES USED:\n")
			type abilityUsed struct {
				name  string
				count int
			}
			var sortedAbilityUsed []abilityUsed
			for name, count := range interruptAbilitiesUsed {
				sortedAbilityUsed = append(sortedAbilityUsed, abilityUsed{name, count})
			}
			// Simple sort by count (descending)
			for i := 0; i < len(sortedAbilityUsed)-1; i++ {
				for j := i + 1; j < len(sortedAbilityUsed); j++ {
					if sortedAbilityUsed[j].count > sortedAbilityUsed[i].count {
						sortedAbilityUsed[i], sortedAbilityUsed[j] = sortedAbilityUsed[j], sortedAbilityUsed[i]
					}
				}
			}

			for _, ability := range sortedAbilityUsed {
				fmt.Printf("  • %s: %s times\n",
					color.HiYellowString(ability.name),
					color.HiBlueString("%d", ability.count))
			}
		}
	}

	fmt.Println()
}

// displayPlayerInterruptAnalysis shows detailed analysis for a specific player
func displayPlayerInterruptAnalysis(events []*models.Event, playerLookup map[int]string, fight *models.Fight, lookupService *services.LookupService, apiClient *api.Client, reportCode string, fightID int, targetPlayerName string, verbose bool) {
	color.HiBlue("\n🎛️  DETAILED INTERRUPT ANALYSIS: %s 🎛️\n", color.HiYellowString(targetPlayerName))

	fightDuration := time.Duration((fight.EndTime - fight.StartTime) * int64(time.Millisecond))
	fmt.Printf("Fight: %s (Duration: %s)\n",
		color.HiYellowString(fight.Name),
		color.HiWhiteString(fightDuration.String()))

	// Find interrupts for this specific player
	var playerInterrupts []*models.Event
	var targetPlayerID int
	for _, event := range events {
		if event.SourceID != nil {
			if name, exists := playerLookup[*event.SourceID]; exists && strings.EqualFold(name, targetPlayerName) {
				playerInterrupts = append(playerInterrupts, event)
				targetPlayerID = *event.SourceID
			}
		}
	}

	if verbose {
		fmt.Printf("🔍 Debug: Found player ID %d for '%s'\n", targetPlayerID, targetPlayerName)
		fmt.Printf("🔍 Debug: Fight start time: %d ms\n", fight.StartTime)
	}

	if len(playerInterrupts) == 0 {
		color.HiGreen("🎉 %s did not perform any interrupts!\n", targetPlayerName)
		return
	}

	fmt.Printf("Interrupts: %s\n\n", color.HiBlueString("%d", len(playerInterrupts)))

	fightStartTime := float64(fight.StartTime)

	// Show timeline of interrupts - here we show what the player CAST to interrupt (e.g., Wind Shear)
	fmt.Printf("⏰ INTERRUPT TIMELINE (Player's Interrupt Ability):\n")
	for _, event := range playerInterrupts {
		timeIntoFight := time.Duration((event.Timestamp - fightStartTime) * float64(time.Millisecond))
		interruptAbilityName := "Unknown Ability"
		if event.AbilityID != nil {
			interruptAbilityName = lookupService.GetAbilityName(*event.AbilityID)
		}
		targetName := "Unknown Target"
		if event.Target != nil {
			targetName = event.Target.Name
		} else if event.TargetID != nil {
			if name, exists := playerLookup[*event.TargetID]; exists {
				targetName = name
			}
		}

		fmt.Printf("  • %s: %s cast on %s\n",
			color.HiWhiteString(timeIntoFight.String()),
			color.HiYellowString(interruptAbilityName),
			color.HiMagentaString(targetName))
	}

	// Show interrupt target breakdown
	interruptTargets := make(map[string]int)

	for _, event := range playerInterrupts {
		targetName := "Unknown Target"
		if event.Target != nil {
			targetName = event.Target.Name
		} else if event.TargetID != nil {
			if name, exists := playerLookup[*event.TargetID]; exists {
				targetName = name
			}
		}
		interruptTargets[targetName]++
	}

	if len(interruptTargets) > 0 {
		fmt.Printf("\n🎯 INTERRUPT TARGETS (What player interrupted):\n")
		for target, count := range interruptTargets {
			fmt.Printf("  • %s: %s\n",
				color.HiMagentaString(target),
				color.HiBlueString("%d interrupts", count))
		}
	}

	// Player-specific insights
	color.HiBlue("📊 INSIGHTS:")
	totalInterrupts := len(playerInterrupts)
	fmt.Printf("• %s performed %d interrupts\n", targetPlayerName, totalInterrupts)

	// Now we need to implement the core logic that finds what was actually interrupted
	// This is where we call the advanced interrupt analysis
	fmt.Printf("\n🔄 CORRELATING WITH TARGET CASTS (Finding what was actually interrupted)...\n")

	// Correlate interrupts with casts to determine what was interrupted vs allowed to complete
	analysis, err := CorrelateInterruptsAndCasts(apiClient, reportCode, fightID, playerInterrupts, verbose, fight.StartTime)
	if err != nil {
		fmt.Printf("❌ Error correlating interrupts with target casts: %v\n", err)
		// Show a fallback message
		fmt.Printf("\n📊 INTERRUPT SUMMARY (without cast correlation):\n")
		fmt.Printf("• Total interrupts performed: %d\n", totalInterrupts)
		fmt.Printf("• Note: Detailed cast correlation failed, unable to determine what was interrupted vs completed\n")
		return
	}

	// Display the correlation results with detailed breakdown
	if len(analysis) > 0 {
		fmt.Printf("\n🏆 CORRELATION ANALYSIS (What abilities were interrupted vs completed):\n")
		for abilityName, details := range analysis {
			stoppedPct := float64(details.Stopped) / float64(details.TotalCasts) * 100
			missedPct := float64(details.Missed) / float64(details.TotalCasts) * 100

			fmt.Printf("\n=== %s ===\n", color.HiYellowString(abilityName))
			fmt.Printf("Total Casts: %d\n", details.TotalCasts)
			fmt.Printf("Stopped: %.1f%% (%d)\n", stoppedPct, details.Stopped)
			fmt.Printf("Completed: %.1f%% (%d)\n", missedPct, details.Missed)

			// Show who interrupted the casts (should be our target player)
			if len(details.InterruptedBy) > 0 {
				fmt.Println("\nInterrupted By:")
				for interrupter, count := range details.InterruptedBy {
					pct := float64(count) / float64(details.TotalCasts) * 100
					fmt.Printf("  %s: %d (%.1f%%)\n", interrupter, count, pct)
				}

				// Show detailed stopped cast log
				fmt.Println("\nStopped Casts:")
				for _, stoppedCast := range details.StoppedCasts {
					// Convert timestamp from milliseconds to seconds and format properly
					timestampSecs := int(stoppedCast.Timestamp / 1000)
					minutes := timestampSecs / 60
					seconds := timestampSecs % 60
					fmt.Printf("  [%d:%02d] %s cast interrupted by %s\n",
						minutes, seconds, stoppedCast.CasterName, stoppedCast.InterruptedBy)
				}
			}

			// Show completed casts if any
			if len(details.MissedCasts) > 0 {
				fmt.Println("\nCompleted Casts (Not Interrupted):")
				for _, missedCast := range details.MissedCasts {
					timestampSecs := int(missedCast.Timestamp / 1000)
					minutes := timestampSecs / 60
					seconds := timestampSecs % 60
					fmt.Printf("  [%d:%02d] %s - Cast Completed\n",
						minutes, seconds, missedCast.CasterName)
				}
			}
		}

		// Summary stats
		totalInterrupted := 0
		totalCompleted := 0
		for _, details := range analysis {
			totalInterrupted += details.Stopped
			totalCompleted += details.Missed
		}
		totalCasts := totalInterrupted + totalCompleted
		if totalCasts > 0 {
			overallPct := float64(totalInterrupted) / float64(totalCasts) * 100
			fmt.Printf("\n📊 OVERALL CORRELATION SUMMARY:\n")
			fmt.Printf("Total Interrupted: %d\n", totalInterrupted)
			fmt.Printf("Total Completed: %d\n", totalCompleted)
			fmt.Printf("Overall Interrupt Effectiveness: %.1f%%\n", overallPct)
		}
	} else {
		fmt.Printf("📊 No correlations found - target may not have cast any abilities during interrupt windows\n")
	}
}

// CastAnalysis holds the analysis results for a specific ability
type CastAnalysis struct {
	AbilityName   string
	TotalCasts    int
	Stopped       int
	Missed        int
	InterruptedBy map[string]int
	StoppedCasts  []StoppedCast
	MissedCasts   []MissedCast
}

// StoppedCast represents a cast that was successfully interrupted
type StoppedCast struct {
	CasterName    string
	InterruptedBy string
	Timestamp     float64
}

// MissedCast represents a cast that was not interrupted
type MissedCast struct {
	CasterName string
	Timestamp  float64
}

// CorrelateInterruptsAndCasts analyzes the relationship between interrupts and casts
func CorrelateInterruptsAndCasts(apiClient *api.Client, reportCode string, fightID int, interruptEvents []*models.Event, verbose bool, fightStartTime int64) (map[string]*CastAnalysis, error) {
	if verbose {
		fmt.Printf("🔍 Fetching hostile cast events to correlate with interrupts...\n")
	}

	// Fetch hostile cast events to see what was cast by enemies
	var startTime *float64 = nil // No pagination in initial call
	castRequest := api.NewAllCastEventsRequest(reportCode, fightID, api.EventHostilityHostile, startTime)
	castResponse, err := apiClient.Query(castRequest.Query, castRequest.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cast events: %w", err)
	}

	if castResponse.Data == nil || castResponse.Data.ReportData == nil ||
		castResponse.Data.ReportData.Report == nil ||
		castResponse.Data.ReportData.Report.Events == nil {
		return map[string]*CastAnalysis{}, nil
	}

	// Parse cast events
	castEvents, err := models.ParseCastEventsJSON(castResponse.Data.ReportData.Report.Events.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cast events: %w", err)
	}

	if verbose {
		fmt.Printf("✅ Found %d cast events to analyze\n", len(castEvents))
	}

	// Load lookup service to get actor names
	lookupService := services.NewLookupService(apiClient)
	err = lookupService.LoadActorsFromReport(reportCode)
	if err != nil {
		return nil, fmt.Errorf("failed to load actors for correlation: %w", err)
	}

	// Create a map of cast events by the target IDs from interrupt events
	// We only care about casts from NPCs that were interrupted
	interruptedNPCs := make(map[int]bool)
	for _, interrupt := range interruptEvents {
		if interrupt.TargetID != nil {
			interruptedNPCs[*interrupt.TargetID] = true
		}
	}

	// Filter cast events to only include casts from NPCs that were interrupted
	// This makes the correlation more focused and accurate
	var relevantCastEvents []*models.Event
	for _, castEvent := range castEvents {
		if castEvent.SourceID != nil && interruptedNPCs[*castEvent.SourceID] {
			relevantCastEvents = append(relevantCastEvents, castEvent)
		}
	}

	if verbose {
		fmt.Printf("✅ Found %d relevant cast events from interrupted NPCs\n", len(relevantCastEvents))
	}

	// Group cast events by ability ID and fetch ability names
	analysis := make(map[string]*CastAnalysis)

	// First, collect all ability IDs from relevant cast events to preload names
	abilityIDs := make(map[int]bool)
	actorIDs := make(map[int]bool)

	for _, castEvent := range relevantCastEvents {
		if castEvent.AbilityID != nil {
			abilityIDs[*castEvent.AbilityID] = true
		}
		if castEvent.SourceID != nil {
			actorIDs[*castEvent.SourceID] = true
		}
	}

	// Preload ability names to avoid multiple API calls
	abilityList := make([]int, 0, len(abilityIDs))
	for id := range abilityIDs {
		abilityList = append(abilityList, id)
	}
	if len(abilityList) > 0 {
		lookupService.PreloadAbilities(abilityList)
	}

	// Create a lookup map of interrupt events by target ID and time for efficient correlation
	interruptMap := make(map[string]*models.Event)
	for _, interrupt := range interruptEvents {
		if interrupt.TargetID != nil {
			// Create key with target ID and timestamp for exact matching
			key := fmt.Sprintf("%d|%.0f", *interrupt.TargetID, interrupt.Timestamp)
			interruptMap[key] = interrupt
		}
	}

	// Also add approximate time matching for interrupts (within 500ms window)
	interruptTimeMap := make(map[int][]*models.Event)
	for _, interrupt := range interruptEvents {
		if interrupt.TargetID != nil {
			interruptTimeMap[*interrupt.TargetID] = append(interruptTimeMap[*interrupt.TargetID], interrupt)
		}
	}

	// Process each relevant cast event
	for _, castEvent := range relevantCastEvents {
		if castEvent.AbilityID == nil || castEvent.SourceID == nil {
			continue
		}

		// Get ability name
		abilityName := lookupService.GetAbilityName(*castEvent.AbilityID)
		if abilityName == "" {
			abilityName = fmt.Sprintf("Unknown Ability (%d)", *castEvent.AbilityID)
		}

		// Initialize analysis for this ability if not exists
		if _, exists := analysis[abilityName]; !exists {
			analysis[abilityName] = &CastAnalysis{
				AbilityName:   abilityName,
				InterruptedBy: make(map[string]int),
				StoppedCasts:  make([]StoppedCast, 0),
				MissedCasts:   make([]MissedCast, 0),
			}
		}

		// Check if this cast was interrupted
		// An interrupted cast means there was an interrupt event from a player on this NPC at a close timestamp
		wasInterrupted := false
		interruptedBy := "Unknown"

		// Look for interrupts on this specific NPC around this timestamp
		if interruptsForNPC, exists := interruptTimeMap[*castEvent.SourceID]; exists {
			for _, interrupt := range interruptsForNPC {
				// Check if the interrupt time is close to cast time (within 300ms)
				timeDiff := math.Abs(castEvent.Timestamp - interrupt.Timestamp)
				if timeDiff <= 300.0 { // 300ms window should capture most interrupts
					wasInterrupted = true
					if interrupt.SourceID != nil {
						interruptedBy = lookupService.GetActorName(*interrupt.SourceID)
						if interruptedBy == "" {
							interruptedBy = fmt.Sprintf("Player-%d", *interrupt.SourceID)
						}
					}
					break // Found the interrupt for this cast
				}
			}
		}

		// Update analysis based on whether cast was interrupted
		analysis[abilityName].TotalCasts++

		casterName := lookupService.GetActorName(*castEvent.SourceID)
		if casterName == "" {
			casterName = fmt.Sprintf("NPC-%d", *castEvent.SourceID)
		}

		if wasInterrupted {
			analysis[abilityName].Stopped++
			// Calculate timestamp relative to fight start time (WCL timestamps are in milliseconds)
			fightRelativeTimestamp := castEvent.Timestamp - float64(fightStartTime)
			if fightRelativeTimestamp < 0 {
				fightRelativeTimestamp = castEvent.Timestamp // fallback if calculation is wrong
			}

			analysis[abilityName].StoppedCasts = append(analysis[abilityName].StoppedCasts, StoppedCast{
				CasterName:    casterName,
				InterruptedBy: interruptedBy,
				Timestamp:     fightRelativeTimestamp,
			})
			analysis[abilityName].InterruptedBy[interruptedBy]++
		} else {
			analysis[abilityName].Missed++
			// Calculate timestamp relative to fight start time (WCL timestamps are in milliseconds)
			fightRelativeTimestamp := castEvent.Timestamp - float64(fightStartTime)
			if fightRelativeTimestamp < 0 {
				fightRelativeTimestamp = castEvent.Timestamp // fallback if calculation is wrong
			}

			analysis[abilityName].MissedCasts = append(analysis[abilityName].MissedCasts, MissedCast{
				CasterName: casterName,
				Timestamp:  fightRelativeTimestamp,
			})
		}
	}

	return analysis, nil
}
