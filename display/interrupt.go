package display

import (
	"fmt"
	"strings"

	"wclogs-cli/models"

	"github.com/fatih/color"
)

// DisplayInterruptAnalysis displays the complete interrupt analysis in the WCL format
func DisplayInterruptAnalysis(analysis *models.InterruptAnalysisResult, showPlayerDetails bool) {
	if analysis == nil {
		color.HiRed("❌ No interrupt analysis data available")
		return
	}

	// Validate analysis and show warnings if needed
	warnings := analysis.ValidateAnalysis()
	if len(warnings) > 0 {
		color.HiYellow("⚠️  Analysis Warnings:")
		for _, warning := range warnings {
			fmt.Printf("   • %s\n", warning)
		}
		fmt.Println()
	}

	// Header with fight information
	displayFightHeader(analysis.FightInfo, analysis.TotalCasts, analysis.TotalStops, analysis.OverallRate)

	// Display total interrupts & stops summary
	displayPlayerSummaryTable(analysis)

	// Display per-spell analysis
	spells := analysis.GetSpellsByActivity()
	for _, spell := range spells {
		if spell.TotalCasts > 0 { // Only show spells that were actually cast
			displaySpellAnalysis(spell)
		}
	}

	// Show player details if requested
	if showPlayerDetails {
		displayPlayerDetails(analysis)
	}
}

// displayFightHeader shows fight information and overall statistics
func displayFightHeader(fightInfo *models.FightInfo, totalCasts, totalStops int, overallRate float64) {
	color.HiCyan("\n🎛️  INTERRUPT ANALYSIS")
	fmt.Printf("⚔️  Fight: %s | 🕐 Duration: %s | 📊 Overall: %d/%d (%.1f%%)\n",
		fightInfo.Name, fightInfo.Duration.String(), totalStops, totalCasts, overallRate)
	fmt.Println(strings.Repeat("=", 80))
}

// displayPlayerSummaryTable shows the "Total Interrupts & Stops" section
func displayPlayerSummaryTable(analysis *models.InterruptAnalysisResult) {
	topPlayers := analysis.GetTopPlayers(0) // Get all players

	if len(topPlayers) == 0 {
		color.HiYellow("ℹ️  No players with interrupts found")
		return
	}

	// Header
	color.New(color.FgHiGreen, color.Bold).Println("Total Interrupts & Stops")
	fmt.Println()

	// Create table frame similar to WCL
	fmt.Printf("┌─────────────────────┬─────────┐\n")
	fmt.Printf("│ %-19s │ %7s │\n", "Player", "Count")
	fmt.Printf("├─────────────────────┼─────────┤\n")

	// Player rows
	for _, player := range topPlayers {
		playerColor := getPlayerNameColor(player.PlayerName)
		fmt.Printf("│ ")
		playerColor.Printf("%-19s", truncateString(player.PlayerName, 19))
		fmt.Printf(" │ %7d │\n", player.TotalStops)
	}

	fmt.Printf("└─────────────────────┴─────────┘\n")
	fmt.Println()
}

// displaySpellAnalysis shows analysis for a specific spell (like "Void Bolt", "Dark Strike")
func displaySpellAnalysis(spell *models.SpellInterruptAnalysis) {
	// Spell header with icon placeholder
	spellColor := color.New(color.FgHiBlue, color.Bold)
	fmt.Printf("🔮 ")
	spellColor.Printf("%s\n", spell.SpellName)
	fmt.Println()

	// Create two-column layout: Stopped | Missed
	displayStoppedAndMissedColumns(spell)
	fmt.Println()
}

// displayStoppedAndMissedColumns creates the side-by-side display
func displayStoppedAndMissedColumns(spell *models.SpellInterruptAnalysis) {
	stoppedTitle := fmt.Sprintf("Stopped: %.1f%% (%d)", spell.StopRate, spell.StoppedCasts)
	missedTitle := fmt.Sprintf("Missed: %.1f%% (%d)", 100-spell.StopRate, spell.MissedCasts)

	// Left column: Stopped
	fmt.Printf("┌─────────────────────────────────────┬─────────────────────────────────────┐\n")
	fmt.Printf("│ ")
	color.New(color.FgHiGreen, color.Bold).Printf("%-35s", stoppedTitle)
	fmt.Printf(" │ ")
	color.New(color.FgHiRed, color.Bold).Printf("%-35s", missedTitle)
	fmt.Printf(" │\n")
	fmt.Printf("├─────────────────────────────────────┼─────────────────────────────────────┤\n")

	// Column headers
	fmt.Printf("│ %-12s %5s %9s   + │ %-15s %-15s %6s │\n",
		"Name", "Count", "% of Casts", "Caster", "Target", "Time")
	fmt.Printf("├─────────────────────────────────────┼─────────────────────────────────────┤\n")

	// Get max rows needed
	maxStoppedRows := len(spell.StoppedBy)
	maxMissedRows := len(spell.MissedDetails)
	maxRows := maxStoppedRows
	if maxMissedRows > maxRows {
		maxRows = maxMissedRows
	}

	// Display rows side by side
	for i := 0; i < maxRows; i++ {
		// Left column (Stopped)
		if i < len(spell.StoppedBy) {
			stopped := spell.StoppedBy[i]
			playerColor := getPlayerNameColor(stopped.PlayerName)
			fmt.Printf("│ ")
			playerColor.Printf("%-12s", truncateString(stopped.PlayerName, 12))
			fmt.Printf(" %5d %8.2f%%   + │",
				stopped.Count, stopped.PercentOfCasts)
		} else {
			fmt.Printf("│                                     │")
		}

		// Right column (Missed)
		if i < len(spell.MissedDetails) {
			missed := spell.MissedDetails[i]
			casterColor := getCasterNameColor(missed.CasterName)
			targetColor := getTargetNameColor(missed.TargetName)
			fmt.Printf(" ")
			casterColor.Printf("%-15s", truncateString(missed.CasterName, 15))
			fmt.Printf(" ")
			targetColor.Printf("%-15s", truncateString(missed.TargetName, 15))
			fmt.Printf(" %6s │\n", missed.TimeInFight)
		} else {
			fmt.Printf("                                     │\n")
		}
	}

	fmt.Printf("└─────────────────────────────────────┴─────────────────────────────────────┘\n")
}

// displayPlayerDetails shows detailed breakdown for individual players
func displayPlayerDetails(analysis *models.InterruptAnalysisResult) {
	fmt.Println()
	color.HiCyan("📋 PLAYER DETAILS")
	fmt.Println(strings.Repeat("-", 60))

	topPlayers := analysis.GetTopPlayers(5) // Show top 5 players
	for i, player := range topPlayers {
		if i > 0 {
			fmt.Println()
		}

		playerColor := getPlayerNameColor(player.PlayerName)
		playerColor.Printf("👤 %s", player.PlayerName)
		fmt.Printf(" - %d total stops", player.TotalStops)

		if len(player.SpellsStoppedDetailed) > 0 {
			fmt.Printf(" (")
			first := true
			for spellName, count := range player.SpellsStoppedDetailed {
				if !first {
					fmt.Printf(", ")
				}
				fmt.Printf("%s: %d", spellName, count)
				first = false
			}
			fmt.Printf(")")
		}
		fmt.Println()
	}
}

// Color helper functions to match WCL styling

func getPlayerNameColor(playerName string) *color.Color {
	// Use different colors for different players to distinguish them
	hash := 0
	for _, char := range playerName {
		hash += int(char)
	}

	colors := []*color.Color{
		color.New(color.FgHiCyan),
		color.New(color.FgHiMagenta),
		color.New(color.FgHiYellow),
		color.New(color.FgHiGreen),
		color.New(color.FgHiBlue),
		color.New(color.FgHiWhite),
	}

	return colors[hash%len(colors)]
}

func getCasterNameColor(casterName string) *color.Color {
	// Red-ish colors for enemy casters
	return color.New(color.FgRed)
}

func getTargetNameColor(targetName string) *color.Color {
	// Player targets get player colors, others get default
	if strings.Contains(targetName, "Shadowguard") || strings.Contains(targetName, "boss") {
		return color.New(color.FgYellow)
	}
	return getPlayerNameColor(targetName)
}

// truncateString truncates a string to the specified length with ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-1] + "…"
}

// DisplayInterruptSummary shows a brief summary for quick overview
func DisplayInterruptSummary(analysis *models.InterruptAnalysisResult) {
	if analysis == nil {
		return
	}

	fmt.Println()
	color.HiCyan("📊 INTERRUPT SUMMARY")
	fmt.Printf("Overall Performance: %d/%d casts interrupted (%.1f%%)\n",
		analysis.TotalStops, analysis.TotalCasts, analysis.OverallRate)

	topPlayers := analysis.GetTopPlayers(3)
	if len(topPlayers) > 0 {
		fmt.Printf("Top Interrupters: ")
		for i, player := range topPlayers {
			if i > 0 {
				fmt.Printf(", ")
			}
			playerColor := getPlayerNameColor(player.PlayerName)
			playerColor.Printf("%s (%d)", player.PlayerName, player.TotalStops)
		}
		fmt.Println()
	}

	activeSpells := 0
	for _, spell := range analysis.SpellAnalysis {
		if spell.TotalCasts > 0 {
			activeSpells++
		}
	}
	fmt.Printf("Spells Analyzed: %d different interruptible abilities\n", activeSpells)
	fmt.Println()
}

// DisplayNoInterruptsFound shows a message when no interrupt data is available
func DisplayNoInterruptsFound(fightName string) {
	fmt.Println()
	color.HiYellow("ℹ️  No Interrupt Data Found")
	fmt.Printf("Fight: %s\n\n", fightName)
	fmt.Println("This could mean:")
	fmt.Println("  • No players performed interrupts during this fight")
	fmt.Println("  • No enemy abilities required interrupting")
	fmt.Println("  • This fight type doesn't involve interruptible mechanics")
	fmt.Println()

	color.HiBlue("💡 Tip: Try a different fight or check if this encounter has interrupt mechanics")
	fmt.Println()
}
