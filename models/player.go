package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// FormatTotal returns the total as a formatted string with commas (ORIGINAL - KEEP)
func (p *Player) FormatTotal() string {
	return FormatNumber(int64(p.Total))
}

// FormatDPS returns formatted DPS string (ORIGINAL - KEEP)
func (p *Player) FormatDPS() string {
	return FormatNumber(int64(p.DPS))
}

// ParseTableData parses raw JSON table data into a TableData struct (ORIGINAL - KEEP)
func ParseTableData(rawJSON json.RawMessage) (*TableData, error) {
	var wrapper TableResponseWrapper
	if err := json.Unmarshal(rawJSON, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse table data: %w", err)
	}

	return &wrapper.Data, nil
}

// GetPlayersFromTable extracts Player objects from TableData (ORIGINAL - KEEP)
func GetPlayersFromTable(tableData *TableData) []*Player {
	players := make([]*Player, 0, len(tableData.Entries))

	for _, entry := range tableData.Entries {
		player := NewPlayerFromEntry(&entry)
		players = append(players, player)
	}

	return players
}

// NewPlayerLookup creates a new PlayerLookup from masterData actors (NEW for Day 6)
func NewPlayerLookup(actors []Actor) *PlayerLookup {
	lookup := &PlayerLookup{
		playersByName: make(map[string]*PlayerInfo),
		playersByID:   make(map[int]*PlayerInfo),
	}

	for _, actor := range actors {
		player := &PlayerInfo{
			ID:     actor.ID,
			Name:   actor.Name,
			Class:  actor.SubType, // SubType contains the class name
			Server: actor.Server,
			Icon:   actor.Icon,
		}

		// Store by name (case-insensitive)
		lookup.playersByName[strings.ToLower(actor.Name)] = player
		// Store by ID
		lookup.playersByID[actor.ID] = player
	}

	return lookup
}

// FindPlayerByName finds a player by name (case-insensitive) (NEW for Day 6)
func (pl *PlayerLookup) FindPlayerByName(name string) (*PlayerInfo, bool) {
	player, exists := pl.playersByName[strings.ToLower(name)]
	return player, exists
}

// FindPlayerByID finds a player by ID (NEW for Day 6)
func (pl *PlayerLookup) FindPlayerByID(id int) (*PlayerInfo, bool) {
	player, exists := pl.playersByID[id]
	return player, exists
}

// GetAllPlayers returns all players sorted by name (NEW for Day 6)
func (pl *PlayerLookup) GetAllPlayers() []*PlayerInfo {
	players := make([]*PlayerInfo, 0, len(pl.playersByName))
	for _, player := range pl.playersByName {
		players = append(players, player)
	}

	// Sort by name
	for i := 0; i < len(players)-1; i++ {
		for j := i + 1; j < len(players); j++ {
			if strings.ToLower(players[i].Name) > strings.ToLower(players[j].Name) {
				players[i], players[j] = players[j], players[i]
			}
		}
	}

	return players
}

// ValidatePlayerName checks if a player name exists and returns suggestions if not (NEW for Day 6)
func (pl *PlayerLookup) ValidatePlayerName(name string) error {
	if _, exists := pl.FindPlayerByName(name); !exists {
		// Find similar names for suggestions
		suggestions := pl.findSimilarNames(name, 3)
		if len(suggestions) > 0 {
			return fmt.Errorf("player '%s' not found. Did you mean: %s?", name, strings.Join(suggestions, ", "))
		}
		return fmt.Errorf("player '%s' not found in this report", name)
	}
	return nil
}

// findSimilarNames finds players with similar names (simple substring matching) (NEW for Day 6)
func (pl *PlayerLookup) findSimilarNames(target string, maxSuggestions int) []string {
	target = strings.ToLower(target)
	var suggestions []string

	for name := range pl.playersByName {
		// Simple similarity check: substring or similar length
		if strings.Contains(name, target) || strings.Contains(target, name) {
			suggestions = append(suggestions, pl.playersByName[name].Name)
			if len(suggestions) >= maxSuggestions {
				break
			}
		}
	}

	return suggestions
}

// FormatNumber formats a number with commas for thousands separators (ORIGINAL - KEEP)
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

// FormatDuration formats seconds into a human readable duration (ORIGINAL - KEEP)
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
