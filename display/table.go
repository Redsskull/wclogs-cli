package display

import (
	"fmt"
	"sort"
	"strings"

	"wclogs-cli/models"
)

// TableOptions configures how the table is displayed
type TableOptions struct {
	TopN      int  // Show only top N players (0 = show all)
	ShowDPS   bool // Show DPS column
	ShowClass bool // Show class column
}

// DefaultTableOptions returns sensible defaults
func DefaultTableOptions() TableOptions {
	return TableOptions{
		TopN:      10,   // Show top 10 by default
		ShowDPS:   true, // Show DPS by default
		ShowClass: true, // Show class by default
	}
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

// printDataRow prints a single data row
func printDataRow(player *models.Player, percentage float64, nameWidth, classWidth, damageWidth, dpsWidth, percentWidth int, options TableOptions) {
	fmt.Printf("%-*s", nameWidth, player.Name)

	if options.ShowClass {
		fmt.Printf("  %-*s", classWidth, player.Class)
	}

	fmt.Printf("  %*s", damageWidth, player.FormatTotal())

	if options.ShowDPS {
		fmt.Printf("  %*s", dpsWidth, player.FormatDPS())
	}

	fmt.Printf("  %*.1f%%", percentWidth-1, percentage)
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

// printSummary prints summary information
func printSummary(totalPlayers, shownPlayers int, totalDamage float64, options TableOptions) {
	if totalPlayers != shownPlayers {
		fmt.Printf("Showing top %d of %d players", shownPlayers, totalPlayers)
	} else {
		fmt.Printf("Showing all %d players", totalPlayers)
	}

	fmt.Printf(" | Total Damage: %s\n", models.FormatNumber(int64(totalDamage)))
}
