package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"wclogs-cli/models"

	"github.com/fatih/color"
)

// OutputData represents the structured data that commands return
type OutputData struct {
	Players    []*models.Player `json:"players"`
	ReportCode string           `json:"report_code"`
	FightID    int              `json:"fight_id"`
	Title      string           `json:"title"`
	Total      int64            `json:"total_damage,omitempty"`
}

// HandleOutput processes the output based on flags - either display to terminal or save to file
func HandleOutput(data *OutputData, outputPath string, topN int, noColor bool, verbose bool) error {
	// If no output file specified, display to terminal
	if outputPath == "" {
		return displayToTerminal(data, topN, noColor)
	}

	// Determine format from file extension
	format := detectFormat(outputPath)
	if format == "" {
		return fmt.Errorf("unsupported file format. Use .csv or .json extension")
	}

	// Create saved_reports directory if it doesn't exist
	reportsDir := "saved_reports"
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	// Prepend the saved_reports directory to the output path
	fullPath := filepath.Join(reportsDir, outputPath)

	if verbose {
		color.HiBlue("💾 Saving to file: %s (format: %s)", fullPath, format)
	}

	// Save to file
	if err := saveToFile(data, fullPath, format, topN); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// Show success message
	color.HiGreen("✅ Data saved to: %s", fullPath)
	color.HiCyan("📊 %d players saved", len(data.Players))

	return nil
}

// detectFormat determines output format from file extension
func detectFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".csv":
		return "csv"
	case ".json":
		return "json"
	default:
		return ""
	}
}

// displayToTerminal shows the data in table format (existing behavior)
func displayToTerminal(data *OutputData, topN int, noColor bool) error {
	// Import and use existing display logic
	// This maintains the current beautiful table display
	fmt.Println("\n" + strings.Repeat("=", 60))
	color.HiCyan("🗡️  WARCRAFT LOGS DAMAGE TABLE")
	fmt.Printf("📊 Report: %s | ⚔️  Fight: %d | 👥 Players: %d\n",
		data.ReportCode, data.FightID, len(data.Players))
	fmt.Println(strings.Repeat("=", 60))

	// Use existing display package - we'll need to refactor this
	// For now, recreate basic table logic here
	players := data.Players
	if topN > 0 && topN < len(players) {
		players = players[:topN]
	}

	// Header
	fmt.Printf("%-13s %-15s %13s %9s %8s\n",
		"Player Name", "Class", "Damage", "DPS", "% Total")
	fmt.Println(strings.Repeat("=", 61))

	// Rows
	for _, player := range players {
		percentage := (player.Total / float64(data.Total)) * 100
		fmt.Printf("%-13s %-15s %13s %9s %7.1f%%\n",
			truncate(player.Name, 13),
			truncate(player.Class, 15),
			formatNumber(int64(player.Total)),
			formatNumber(int64(player.DPS)),
			percentage)
	}

	fmt.Println(strings.Repeat("=", 61))
	if topN > 0 && topN < len(data.Players) {
		fmt.Printf("📊 Showing top %d of %d players | Total Damage: %s\n",
			topN, len(data.Players), formatNumber(data.Total))
	} else {
		fmt.Printf("📊 Total Damage: %s\n", formatNumber(data.Total))
	}

	fmt.Println("\n🎨 Color Legend:  DPS  |  Healers  |  Tanks  |  Unknown")
	color.HiGreen("\n🎉 Success! Damage table displayed!")

	return nil
}

// saveToFile writes data to file in specified format
func saveToFile(data *OutputData, filename string, format string, topN int) error {
	// Apply topN filter if specified
	players := data.Players
	if topN > 0 && topN < len(players) {
		players = players[:topN]
		data = &OutputData{
			Players:    players,
			ReportCode: data.ReportCode,
			FightID:    data.FightID,
			Title:      data.Title,
			Total:      data.Total,
		}
	}

	switch format {
	case "csv":
		return saveCSV(data, filename)
	case "json":
		return saveJSON(data, filename)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// saveCSV writes data as CSV
func saveCSV(data *OutputData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	if err := writer.Write([]string{
		"Player Name", "Class", "Damage", "DPS", "Percent", "Report Code", "Fight ID",
	}); err != nil {
		return err
	}

	// Data rows
	for _, player := range data.Players {
		percentage := (player.Total / float64(data.Total)) * 100
		record := []string{
			player.Name,
			player.Class,
			fmt.Sprintf("%.0f", player.Total),
			fmt.Sprintf("%.0f", player.DPS),
			fmt.Sprintf("%.1f", percentage),
			data.ReportCode,
			fmt.Sprintf("%d", data.FightID),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// saveJSON writes data as JSON
func saveJSON(data *OutputData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	return encoder.Encode(data)
}

// Helper functions
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-1] + "…"
}

func formatNumber(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	// Add commas for thousands
	result := ""
	for i, digit := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}
