package models

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// SortPlayersByTotal sorts players by their total value (descending order)
// This is useful for damage/healing tables where you want top performers first
func SortPlayersByTotal(players []*Player) {
	sort.Slice(players, func(i, j int) bool {
		return players[i].Total > players[j].Total
	})
}

// FilterPlayersByClass filters players to only include the specified class
func FilterPlayersByClass(players []*Player, class string) []*Player {
	var filtered []*Player
	for _, player := range players {
		if strings.EqualFold(player.Class, class) {
			filtered = append(filtered, player)
		}
	}
	return filtered
}

// CalculatePercentage calculates what percentage of the total this player represents
func (p *Player) CalculatePercentage(totalSum float64) float64 {
	if totalSum == 0 {
		return 0
	}
	return (p.Total / totalSum) * 100
}

// GetTableSum calculates the total sum of all players in the table
func (t *TableData) GetTableSum() float64 {
	var sum float64
	for _, entry := range t.Entries {
		sum += entry.Total
	}
	return sum
}

// GetPlayerCount returns the number of players in the table
func (t *TableData) GetPlayerCount() int {
	return len(t.Entries)
}

// FindPlayerByName finds a player entry by name (case-insensitive)
func (t *TableData) FindPlayerByName(name string) *PlayerEntry {
	for i := range t.Entries {
		if strings.EqualFold(t.Entries[i].Name, name) {
			return &t.Entries[i]
		}
	}
	return nil
}

// GetTopPlayers returns the top N players sorted by total damage/healing
func GetTopPlayers(players []*Player, n int) []*Player {
	if n <= 0 || n >= len(players) {
		// Return all players if n is invalid or >= total count
		sorted := make([]*Player, len(players))
		copy(sorted, players)
		SortPlayersByTotal(sorted)
		return sorted
	}

	// Sort and return top N
	sorted := make([]*Player, len(players))
	copy(sorted, players)
	SortPlayersByTotal(sorted)
	return sorted[:n]
}

// GetClassBreakdown returns a map of class name to total damage/healing
func GetClassBreakdown(players []*Player) map[string]float64 {
	breakdown := make(map[string]float64)
	for _, player := range players {
		breakdown[player.Class] += player.Total
	}
	return breakdown
}

// FormatClassBreakdown returns a formatted string of class breakdown
func FormatClassBreakdown(breakdown map[string]float64) string {
	if len(breakdown) == 0 {
		return "No data available"
	}

	var result strings.Builder
	for class, total := range breakdown {
		result.WriteString(fmt.Sprintf("%s: %s\n", class, FormatNumber(int64(total))))
	}
	return strings.TrimSpace(result.String())
}

// ParseEventsJSON parses the raw JSON data from Events API into Event structs
func ParseEventsJSON(data json.RawMessage) ([]*Event, error) {
	var events []*Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("failed to parse events JSON: %w", err)
	}
	return events, nil
}

// ParseInterruptEventsJSON parses interrupt events specifically from raw JSON data
func ParseInterruptEventsJSON(data json.RawMessage) ([]*Event, error) {
	var events []*Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("failed to parse interrupt events JSON: %w", err)
	}

	// Filter to only interrupt events
	var interruptEvents []*Event
	for _, event := range events {
		if event.Type == "interrupt" {
			interruptEvents = append(interruptEvents, event)
		}
	}

	return interruptEvents, nil
}

// ParseCastEventsJSON parses cast events specifically from raw JSON data
func ParseCastEventsJSON(data json.RawMessage) ([]*Event, error) {
	var events []*Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("failed to parse cast events JSON: %w", err)
	}

	// Filter to only cast events
	var castEvents []*Event
	for _, event := range events {
		if event.Type == "cast" || event.Type == "begincast" {
			castEvents = append(castEvents, event)
		}
	}

	return castEvents, nil
}

// GetPlayerLookup creates a player ID to name mapping
func GetPlayerLookup(actors []*Actor) map[int]string {
	lookup := make(map[int]string)
	for _, actor := range actors {
		lookup[actor.ID] = actor.Name
	}
	return lookup
}
