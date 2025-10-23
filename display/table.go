package display

import (
	"fmt"
	"sort"
	"strings"

	"wclogs-cli/models"

	"github.com/fatih/color"
)

// TableOptions configures how the table is displayed
type TableOptions struct {
	TopN      int  // Show only top N players (0 = show all)
	ShowRate  bool // Show rate column (DPS/HPS/etc.)
	ShowClass bool // Show class column
	UseColors bool // Enable color coding by class role
}

// DefaultTableOptions returns sensible defaults
func DefaultTableOptions() TableOptions {
	return TableOptions{
		TopN:      10,   // Show top 10 by default
		ShowRate:  true, // Show rate by default (DPS/HPS/etc.)
		ShowClass: true, // Show class by default
		UseColors: true, // Enable colors by default
	}
}

// DataTypeInfo contains display information for different data types
type DataTypeInfo struct {
	ValueLabel string // "Damage", "Healing", "Deaths", etc.
	RateLabel  string // "DPS", "HPS", "Death Rate", etc.
	TotalLabel string // "Total Damage", "Total Healing", etc.
}

// getDataTypeInfo returns display info for the given data type
func getDataTypeInfo(dataType string) DataTypeInfo {
	switch strings.ToLower(dataType) {
	case "damage":
		return DataTypeInfo{
			ValueLabel: "Damage",
			RateLabel:  "DPS",
			TotalLabel: "Total Damage",
		}
	case "healing":
		return DataTypeInfo{
			ValueLabel: "Healing",
			RateLabel:  "HPS",
			TotalLabel: "Total Healing",
		}
	case "deaths":
		return DataTypeInfo{
			ValueLabel: "Deaths",
			RateLabel:  "Deaths/Min",
			TotalLabel: "Total Deaths",
		}
	case "interrupts":
		return DataTypeInfo{
			ValueLabel: "Interrupts",
			RateLabel:  "Int/Min",
			TotalLabel: "Total Interrupts",
		}
	default:
		return DataTypeInfo{
			ValueLabel: "Value",
			RateLabel:  "Rate",
			TotalLabel: "Total",
		}
	}
}

// getClassColor returns the appropriate color for a class based on role
func getClassColor(class string) *color.Color {
	// DPS classes - Bright Red
	dpsClasses := map[string]bool{
		"Mage": true, "Warlock": true, "Hunter": true, "Rogue": true,
		"DeathKnight": true, "DemonHunter": true, "Warrior": true,
		"Monk": true, // Can be DPS or healer, but often DPS
	}

	// Healer classes - Bright Green
	healerClasses := map[string]bool{
		"Priest": true, "Druid": true, "Paladin": true, "Shaman": true,
		"Evoker": true, // Evokers are primarily healers
	}

	if dpsClasses[class] {
		return color.New(color.FgHiRed, color.Bold)
	} else if healerClasses[class] {
		return color.New(color.FgHiGreen, color.Bold)
	}

	// Default color for unknown classes - Bright Yellow
	return color.New(color.FgHiYellow, color.Bold)
}

// filterMeaningfulPlayers removes players with no relevant data for certain data types
func filterMeaningfulPlayers(players []*models.Player, dataType string) []*models.Player {
	switch strings.ToLower(dataType) {
	case "deaths", "interrupts":
		// Only show players who actually have deaths/interrupts
		var filtered []*models.Player
		for _, player := range players {
			if player.Total > 0 {
				filtered = append(filtered, player)
			}
		}
		return filtered
	default:
		// For damage/healing, show all players (even 0 values can be meaningful)
		return players
	}
}

// DisplayTable displays a formatted table for any data type (replaces DisplayDamageTable)
func DisplayTable(players []*models.Player, dataType string, options TableOptions) {
	if len(players) == 0 {
		fmt.Println("No player data found.")
		return
	}

	// Get display info for this data type
	typeInfo := getDataTypeInfo(dataType)

	// Filter out players with no meaningful data for certain data types
	filteredPlayers := filterMeaningfulPlayers(players, dataType)

	if len(filteredPlayers) == 0 {
		fmt.Printf("ℹ️  No %s data found for this fight.\n", strings.ToLower(typeInfo.ValueLabel))
		fmt.Printf("This could mean:\n")
		switch strings.ToLower(dataType) {
		case "deaths":
			fmt.Printf("  • No players died (great job!)\n")
		case "interrupts":
			fmt.Printf("  • No interrupts were performed\n")
			fmt.Printf("  • This fight may not require interrupts\n")
		}
		return
	}

	// Sort players by total (descending)
	sortedPlayers := make([]*models.Player, len(filteredPlayers))
	copy(sortedPlayers, filteredPlayers)
	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i].Total > sortedPlayers[j].Total
	})

	// Limit to topN if specified
	if options.TopN > 0 && len(sortedPlayers) > options.TopN {
		sortedPlayers = sortedPlayers[:options.TopN]
	}

	// Calculate column widths using modernized max()
	nameWidth := max(calculateMaxWidth(sortedPlayers, func(p *models.Player) string { return p.Name }), 12)

	classWidth := 0
	if options.ShowClass {
		classWidth = max(calculateMaxWidth(sortedPlayers, func(p *models.Player) string { return p.Class }), 8)
	}

	valueWidth := max(
		calculateMaxWidth(sortedPlayers, func(p *models.Player) string { return p.FormatTotal() }),
		len(typeInfo.ValueLabel),
	)

	rateWidth := 0
	if options.ShowRate {
		rateWidth = max(
			calculateMaxWidth(sortedPlayers, func(p *models.Player) string { return p.FormatDPS() }),
			len(typeInfo.RateLabel),
		)
	}

	percentWidth := 8 // "% Total"

	// Print header
	fmt.Println()
	printGenericHeader(nameWidth, classWidth, valueWidth, rateWidth, percentWidth, typeInfo, options)
	printSeparator(nameWidth, classWidth, valueWidth, rateWidth, percentWidth, options)

	// Calculate total for percentage calculation
	totalValue := calculateTotal(sortedPlayers)

	// Print data rows
	for _, player := range sortedPlayers {
		percentage := 0.0
		if totalValue > 0 {
			percentage = (player.Total / totalValue) * 100
		}
		printGenericDataRow(player, percentage, nameWidth, classWidth, valueWidth, rateWidth, percentWidth, options)
	}

	// Print separator and summary
	printSeparator(nameWidth, classWidth, valueWidth, rateWidth, percentWidth, options)
	printGenericSummary(len(filteredPlayers), len(sortedPlayers), totalValue, typeInfo, options)
	fmt.Println()
}

// calculateMaxWidth calculates the maximum width needed for a column
func calculateMaxWidth(players []*models.Player, getter func(*models.Player) string) int {
	maxWidth := 0
	for _, player := range players {
		maxWidth = max(maxWidth, len(getter(player)))
	}
	return maxWidth
}

// printGenericHeader prints the table header based on data type
func printGenericHeader(nameWidth, classWidth, valueWidth, rateWidth, percentWidth int, typeInfo DataTypeInfo, options TableOptions) {
	fmt.Printf("%-*s", nameWidth, "Player Name")

	if options.ShowClass {
		fmt.Printf("  %-*s", classWidth, "Class")
	}

	fmt.Printf("  %*s", valueWidth, typeInfo.ValueLabel)

	if options.ShowRate {
		fmt.Printf("  %*s", rateWidth, typeInfo.RateLabel)
	}

	fmt.Printf("  %*s", percentWidth, "% Total")
	fmt.Println()
}

// printSeparator prints a separator line
func printSeparator(nameWidth, classWidth, valueWidth, rateWidth, percentWidth int, options TableOptions) {
	totalWidth := nameWidth

	if options.ShowClass {
		totalWidth += 2 + classWidth
	}

	totalWidth += 2 + valueWidth

	if options.ShowRate {
		totalWidth += 2 + rateWidth
	}

	totalWidth += 2 + percentWidth

	fmt.Println(strings.Repeat("=", totalWidth))
}

// printGenericDataRow prints a single data row with optional color coding
func printGenericDataRow(player *models.Player, percentage float64, nameWidth, classWidth, valueWidth, rateWidth, percentWidth int, options TableOptions) {
	if options.UseColors {
		classColor := getClassColor(player.Class)

		// Print colored name
		classColor.Printf("%-*s", nameWidth, player.Name)

		if options.ShowClass {
			fmt.Printf("  ")
			classColor.Printf("%-*s", classWidth, player.Class)
		}

		// Print value and rate in normal color
		fmt.Printf("  %*s", valueWidth, player.FormatTotal())

		if options.ShowRate {
			fmt.Printf("  %*s", rateWidth, player.FormatDPS())
		}

		fmt.Printf("  %*.1f%%", percentWidth-1, percentage)
	} else {
		// Non-colored output
		fmt.Printf("%-*s", nameWidth, player.Name)

		if options.ShowClass {
			fmt.Printf("  %-*s", classWidth, player.Class)
		}

		fmt.Printf("  %*s", valueWidth, player.FormatTotal())

		if options.ShowRate {
			fmt.Printf("  %*s", rateWidth, player.FormatDPS())
		}

		fmt.Printf("  %*.1f%%", percentWidth-1, percentage)
	}
	fmt.Println()
}

// calculateTotal calculates the sum of all values
func calculateTotal(players []*models.Player) float64 {
	total := 0.0
	for _, player := range players {
		total += player.Total
	}
	return total
}

// printGenericSummary prints summary information
func printGenericSummary(totalPlayers, shownPlayers int, totalValue float64, typeInfo DataTypeInfo, options TableOptions) {
	// Summary statistics in bold yellow
	summaryColor := color.New(color.FgYellow, color.Bold)

	if totalPlayers != shownPlayers {
		summaryColor.Printf("📊 Showing top %d of %d players", shownPlayers, totalPlayers)
	} else {
		summaryColor.Printf("📊 Showing all %d players", totalPlayers)
	}

	summaryColor.Printf(" | %s: %s", typeInfo.TotalLabel, models.FormatNumber(int64(totalValue)))
	fmt.Println()

	// Add a legend for colors
	if options.UseColors {
		fmt.Println()
		fmt.Print("🎨 Color Legend: ")
		color.New(color.FgHiRed, color.Bold).Print(" DPS ")
		fmt.Print(" | ")
		color.New(color.FgHiGreen, color.Bold).Print(" Healers ")
		fmt.Print(" | ")
		color.New(color.FgHiYellow, color.Bold).Print(" Unknown ")
		fmt.Println()
	}
}
