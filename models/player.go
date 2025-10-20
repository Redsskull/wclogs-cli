package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Player represents a simplified view of player data for display purposes
// This is derived from PlayerEntry but with a cleaner interface
type Player struct {
	Name      string  `json:"name"`
	Class     string  `json:"class"`
	Total     float64 `json:"total"`
	Icon      string  `json:"icon"`
	ItemLevel int     `json:"itemLevel"`
	DPS       float64 `json:"dps"`
}

// NewPlayerFromEntry creates a Player from a PlayerEntry
func NewPlayerFromEntry(entry *PlayerEntry) *Player {
	return &Player{
		Name:      entry.Name,
		Class:     entry.Type,
		Total:     entry.Total,
		Icon:      entry.Icon,
		ItemLevel: entry.ItemLevel,
		DPS:       entry.DPS(),
	}
}

// FormatTotal returns the total as a formatted string with commas
func (p *Player) FormatTotal() string {
	return FormatNumber(int64(p.Total))
}

// FormatDPS returns formatted DPS string
func (p *Player) FormatDPS() string {
	return FormatNumber(int64(p.DPS))
}

// ParseTableData parses raw JSON table data into a TableData struct
func ParseTableData(rawJSON json.RawMessage) (*TableData, error) {
	var wrapper TableResponseWrapper
	if err := json.Unmarshal(rawJSON, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse table data: %w", err)
	}

	return &wrapper.Data, nil
}

// GetPlayersFromTable extracts Player objects from TableData
func GetPlayersFromTable(tableData *TableData) []*Player {
	players := make([]*Player, 0, len(tableData.Entries))

	for _, entry := range tableData.Entries {
		player := NewPlayerFromEntry(&entry)
		players = append(players, player)
	}

	return players
}

// FormatNumber formats a number with commas for thousands separators
func FormatNumber(n int64) string {
	str := strconv.FormatInt(n, 10)
	if len(str) <= 3 {
		return str
	}

	// Add commas every 3 digits from the right
	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}
	return result.String()
}

// FormatDuration formats seconds into a human readable duration
func FormatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60

	if minutes < 60 {
		return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60

	return fmt.Sprintf("%dh %dm %ds", hours, remainingMinutes, remainingSeconds)
}
