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
	ShowDPS   bool // Show DPS column
	ShowClass bool // Show class column
	UseColors bool // Enable color coding by class role
}

// DefaultTableOptions returns sensible defaults
func DefaultTableOptions() TableOptions {
	return TableOptions{
		TopN:      10,   // Show top 10 by default
		ShowDPS:   true, // Show DPS by default
		ShowClass: true, // Show class by default
		UseColors: true, // Enable colors by default
	}
}

// getClassColor returns the appropriate color for a class based on role
func getClassColor(class string) *color.Color {
	// DPS classes - Bright Red with Background
	dpsClasses := map[string]bool{
		"Mage": true, "Warlock": true, "Hunter": true, "Rogue": true,
		"DeathKnight": true, "DemonHunter": true, "Evoker": true,
		"Monk": true, "Shaman": true, // Assuming DPS specs
	}

	// Healer classes - Bright Green with Background
	healerClasses := map[string]bool{
		"Priest": true, "Druid": true, "Paladin": true,
	}

	// Tank classes - Bright Blue with Background (though many can DPS too)
	tankClasses := map[string]bool{
		"Warrior": true,
	}

	if dpsClasses[class] {
		return color.New(color.FgHiRed, color.Bold, color.BgHiBlack)
	} else if healerClasses[class] {
		return color.New(color.FgHiGreen, color.Bold, color.BgHiBlack)
	} else if tankClasses[class] {
		return color.New(color.FgHiBlue, color.Bold, color.BgHiBlack)
	}

	// Default color for unknown classes - Bright Yellow
	return color.New(color.FgHiYellow, color.Bold)
}

// DisplayDamageTable displays a formatted damage table
func DisplayDamageTable(players []*models.Player, options TableOptions) {
	if len(players) == 0 {
		fmt.Println("No player data found.")
		return
	}

	// Sort players by total damage (descending)
	sortedPlayers := make([]*models.Player, len(players))
	copy(sortedPlayers, players)
	sort.Slice(sortedPlayers, func(i, j int) bool {
		return sortedPlayers[i].Total > sortedPlayers[j].Total
	})

	// Limit to topN if specified
	if options.TopN > 0 && len(sortedPlayers) > options.TopN {
		sortedPlayers = sortedPlayers[:options.TopN]
	}

	// Calculate column widths
	nameWidth := calculateMaxWidth(sortedPlayers, func(p *models.Player) string { return p.Name })
	if nameWidth < 12 { // Minimum width for "Player Name"
		nameWidth = 12
	}

	classWidth := 0
	if options.ShowClass {
		classWidth = calculateMaxWidth(sortedPlayers, func(p *models.Player) string { return p.Class })
		if classWidth < 8 { // Minimum width for "Class"
			classWidth = 8
		}
	}

	damageWidth := calculateMaxWidth(sortedPlayers, func(p *models.Player) string { return p.FormatTotal() })
	if damageWidth < 10 { // Minimum width for "Damage"
		damageWidth = 10
	}

	dpsWidth := 0
	if options.ShowDPS {
		dpsWidth = calculateMaxWidth(sortedPlayers, func(p *models.Player) string { return p.FormatDPS() })
		if dpsWidth < 8 { // Minimum width for "DPS"
			dpsWidth = 8
		}
	}

	percentWidth := 8 // "% Total"

	// Print header
	fmt.Println()
	printHeader(nameWidth, classWidth, damageWidth, dpsWidth, percentWidth, options)
	printSeparator(nameWidth, classWidth, damageWidth, dpsWidth, percentWidth, options)

	// Calculate total damage for percentage calculation
	totalDamage := calculateTotalDamage(sortedPlayers)

	// Print data rows
	for _, player := range sortedPlayers {
		percentage := (player.Total / totalDamage) * 100
		printDataRow(player, percentage, nameWidth, classWidth, damageWidth, dpsWidth, percentWidth, options)
	}

	// Print separator and summary
	printSeparator(nameWidth, classWidth, damageWidth, dpsWidth, percentWidth, options)
	printSummary(len(players), len(sortedPlayers), totalDamage, options)
	fmt.Println()
}

// calculateMaxWidth calculates the maximum width needed for a column
func calculateMaxWidth(players []*models.Player, getter func(*models.Player) string) int {
	maxWidth := 0
	for _, player := range players {
		width := len(getter(player))
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

// printHeader prints the table header
func printHeader(nameWidth, classWidth, damageWidth, dpsWidth, percentWidth int, options TableOptions) {
	fmt.Printf("%-*s", nameWidth, "Player Name")

	if options.ShowClass {
		fmt.Printf("  %-*s", classWidth, "Class")
	}

	fmt.Printf("  %*s", damageWidth, "Damage")

	if options.ShowDPS {
		fmt.Printf("  %*s", dpsWidth, "DPS")
	}

	fmt.Printf("  %*s", percentWidth, "% Total")
	fmt.Println()
}

// printSeparator prints a separator line
func printSeparator(nameWidth, classWidth, damageWidth, dpsWidth, percentWidth int, options TableOptions) {
	totalWidth := nameWidth

	if options.ShowClass {
		totalWidth += 2 + classWidth
	}

	totalWidth += 2 + damageWidth

	if options.ShowDPS {
		totalWidth += 2 + dpsWidth
	}

	totalWidth += 2 + percentWidth

	fmt.Println(strings.Repeat("=", totalWidth))
}

// printDataRow prints a single data row with optional color coding
func printDataRow(player *models.Player, percentage float64, nameWidth, classWidth, damageWidth, dpsWidth, percentWidth int, options TableOptions) {
	if options.UseColors {
		classColor := getClassColor(player.Class)

		// Print colored name
		classColor.Printf("%-*s", nameWidth, player.Name)

		if options.ShowClass {
			fmt.Printf("  ")
			classColor.Printf("%-*s", classWidth, player.Class)
		}

		// Print damage and DPS in normal color
		fmt.Printf("  %*s", damageWidth, player.FormatTotal())

		if options.ShowDPS {
			fmt.Printf("  %*s", dpsWidth, player.FormatDPS())
		}

		fmt.Printf("  %*.1f%%", percentWidth-1, percentage)
	} else {
		// Original non-colored output
		fmt.Printf("%-*s", nameWidth, player.Name)

		if options.ShowClass {
			fmt.Printf("  %-*s", classWidth, player.Class)
		}

		fmt.Printf("  %*s", damageWidth, player.FormatTotal())

		if options.ShowDPS {
			fmt.Printf("  %*s", dpsWidth, player.FormatDPS())
		}

		fmt.Printf("  %*.1f%%", percentWidth-1, percentage)
	}
	fmt.Println()
}

// calculateTotalDamage calculates the sum of all damage
func calculateTotalDamage(players []*models.Player) float64 {
	total := 0.0
	for _, player := range players {
		total += player.Total
	}
	return total
}

// printSummary prints summary information with a total damage row
func printSummary(totalPlayers, shownPlayers int, totalDamage float64, options TableOptions) {
	// Print a total row separator
	fmt.Println()

	// Summary statistics in bold yellow
	summaryColor := color.New(color.FgYellow, color.Bold)

	if totalPlayers != shownPlayers {
		summaryColor.Printf("📊 Showing top %d of %d players", shownPlayers, totalPlayers)
	} else {
		summaryColor.Printf("📊 Showing all %d players", totalPlayers)
	}

	summaryColor.Printf(" | Total Damage: %s", models.FormatNumber(int64(totalDamage)))
	fmt.Println()

	// Add a legend for colors
	if options.UseColors {
		fmt.Println()
		fmt.Print("🎨 Color Legend: ")
		color.New(color.FgHiRed, color.Bold, color.BgHiBlack).Print(" DPS ")
		fmt.Print(" | ")
		color.New(color.FgHiGreen, color.Bold, color.BgHiBlack).Print(" Healers ")
		fmt.Print(" | ")
		color.New(color.FgHiBlue, color.Bold, color.BgHiBlack).Print(" Tanks ")
		fmt.Print(" | ")
		color.New(color.FgHiYellow, color.Bold).Print(" Unknown ")
		fmt.Println()
	}
}
