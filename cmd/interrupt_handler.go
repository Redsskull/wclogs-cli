package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/display"
	"wclogs-cli/models"
	"wclogs-cli/services"
)

// ExecuteInterruptAnalysis provides comprehensive interrupt analysis using the correlation system
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

	// Get fight information first
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

	// Load all actors for name lookups
	if verbose {
		color.HiBlue("👥 Loading actors and game data...")
	}

	err = lookupService.LoadActorsFromReport(reportCode)
	if err != nil {
		return fmt.Errorf("failed to load actors: %w", err)
	}

	// Create fight info for the correlator
	fightInfo := &models.FightInfo{
		FightID:   fightID,
		Name:      currentFight.Name,
		Duration:  time.Duration((currentFight.EndTime - currentFight.StartTime) * int64(time.Millisecond)),
		StartTime: currentFight.StartTime,
		EndTime:   currentFight.EndTime,
	}

	// Create interrupt correlator
	correlator := services.NewInterruptCorrelator(apiClient, lookupService, reportCode, fightID, fightInfo, verbose)

	// Perform the comprehensive analysis
	if verbose {
		color.HiBlue("🔗 Performing interrupt correlation analysis...")
	}

	analysis, err := correlator.AnalyzeInterrupts()
	if err != nil {
		return fmt.Errorf("interrupt analysis failed: %w", err)
	}

	// Display results based on player filter
	if playerName != "" {
		// Player-specific detailed analysis
		displayPlayerSpecificAnalysis(analysis, playerName)
	} else {
		// Full analysis display
		display.DisplayInterruptAnalysis(analysis, verbose)
	}

	// Show cache statistics if verbose
	if verbose {
		abilityCount, actorCount := lookupService.GetCacheStats()
		color.HiBlue("📊 Cache Stats: %d abilities, %d actors cached", abilityCount, actorCount)
	}

	return nil
}

// displayPlayerSpecificAnalysis shows detailed analysis for a specific player
func displayPlayerSpecificAnalysis(analysis *models.InterruptAnalysisResult, playerName string) {
	player, exists := analysis.PlayerSummary[playerName]
	if !exists {
		color.HiYellow("⚠️  Player '%s' not found in interrupt data", playerName)

		// Show available players
		if len(analysis.PlayerSummary) > 0 {
			color.HiCyan("\nAvailable players with interrupts:")
			topPlayers := analysis.GetTopPlayers(0)
			for _, p := range topPlayers {
				fmt.Printf("  • %s (%d interrupts)\n", p.PlayerName, p.TotalStops)
			}
		} else {
			color.HiYellow("No players with interrupts found in this fight")
		}
		return
	}

	// Display player-specific header
	color.HiCyan("\n👤 INTERRUPT ANALYSIS: %s", playerName)
	fmt.Printf("📊 Total Stops: %d | Success Rate: %.1f%%\n", player.TotalStops, player.SuccessRate)
	fmt.Println(strings.Repeat("=", 60))

	if player.TotalStops == 0 {
		color.HiYellow("This player performed no interrupts during this fight")
		return
	}

	// Show detailed breakdown by spell
	color.HiGreen("\n🔮 SPELLS INTERRUPTED:")
	for spellName, count := range player.SpellsStoppedDetailed {
		spell := analysis.SpellAnalysis[spellName]
		if spell != nil {
			percentage := float64(count) / float64(spell.TotalCasts) * 100
			fmt.Printf("  • %s: %d/%d stops (%.1f%% of casts)\n",
				spellName, count, spell.TotalCasts, percentage)
		} else {
			fmt.Printf("  • %s: %d stops\n", spellName, count)
		}
	}

	// Show timing analysis for this player
	displayPlayerTimingAnalysis(analysis, playerName)
}

// displayPlayerTimingAnalysis shows when a specific player performed their interrupts
func displayPlayerTimingAnalysis(analysis *models.InterruptAnalysisResult, playerName string) {
	color.HiBlue("\n⏱️  INTERRUPT TIMING:")

	spells := analysis.GetSpellsByActivity()
	hasInterrupts := false

	for _, spell := range spells {
		for _, detail := range spell.StoppedBy {
			if detail.PlayerName == playerName {
				hasInterrupts = true
				color.HiWhite("\n%s (%d interrupts):", spell.SpellName)

				for i, timestamp := range detail.Timestamps {
					fightTime := time.Duration((timestamp - float64(analysis.FightInfo.StartTime)) * float64(time.Millisecond))
					target := "Unknown"
					if i < len(detail.Targets) {
						target = detail.Targets[i]
					}

					fmt.Printf("  %s → %s\n", formatFightTime(fightTime), target)
				}
			}
		}
	}

	if !hasInterrupts {
		color.HiYellow("No detailed timing data available")
	}
}

// formatFightTime formats a duration as MM:SS for fight time display
func formatFightTime(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// validateInterruptCapability checks if a player has interrupt abilities
func validateInterruptCapability(playerName string, playerClass string) bool {
	// Define classes with interrupt abilities
	interruptClasses := map[string]bool{
		"Death Knight": true,
		"Demon Hunter": true,
		"Hunter":       true,
		"Mage":         true,
		"Monk":         true,
		"Paladin":      true,
		"Rogue":        true,
		"Shaman":       true,
		"Warlock":      true,
		"Warrior":      true,
		// Note: Some specs have interrupts while others don't, but this is a general check
	}

	return interruptClasses[playerClass]
}

// displayInterruptCapabilityWarning shows a warning if a class typically can't interrupt
func displayInterruptCapabilityWarning(playerName string, playerClass string) {
	if !validateInterruptCapability(playerName, playerClass) {
		color.HiYellow("⚠️  Note: %s (%s) may not have interrupt abilities in their current spec",
			playerName, playerClass)
	}
}
