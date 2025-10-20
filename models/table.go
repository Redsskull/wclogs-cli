package models

import "encoding/json"

// TableResponseWrapper represents the outer wrapper of the table response
// The actual structure is: {"data": {"entries": [...]}, "totalTime": ..., etc}
type TableResponseWrapper struct {
	Data            TableData `json:"data"`
	TotalTime       int64     `json:"totalTime"`
	DamageDowntime  int64     `json:"damageDowntime,omitempty"`
	HealingDowntime int64     `json:"healingDowntime,omitempty"`
	LogVersion      int       `json:"logVersion,omitempty"`
	GameVersion     int       `json:"gameVersion,omitempty"`
}

// TableData represents the inner data structure containing the entries
type TableData struct {
	Entries []PlayerEntry `json:"entries"`
}

// PlayerEntry represents a single player's data from the table
// This includes basic stats plus complex nested data
type PlayerEntry struct {
	Name              string  `json:"name"`
	ID                int     `json:"id"`
	GUID              int64   `json:"guid"`
	Type              string  `json:"type"` // Player class
	Icon              string  `json:"icon"` // Class icon
	ItemLevel         int     `json:"itemLevel"`
	Total             float64 `json:"total"` // Total damage/healing
	TotalReduced      float64 `json:"totalReduced,omitempty"`
	ActiveTime        int64   `json:"activeTime"`
	ActiveTimeReduced int64   `json:"activeTimeReduced,omitempty"`

	// Complex nested data - keeping as RawMessage for now to avoid deep parsing
	Abilities       json.RawMessage `json:"abilities,omitempty"`
	DamageAbilities json.RawMessage `json:"damageAbilities,omitempty"`
	Targets         json.RawMessage `json:"targets,omitempty"`
	Talents         json.RawMessage `json:"talents,omitempty"`
	Gear            json.RawMessage `json:"gear,omitempty"`
	Pets            json.RawMessage `json:"pets,omitempty"`
}

// FormatTotal returns the total as a formatted string with commas
func (p *PlayerEntry) FormatTotal() string {
	return FormatNumber(int64(p.Total))
}

// FormatActiveTime returns active time in a human readable format (seconds)
func (p *PlayerEntry) FormatActiveTime() string {
	seconds := p.ActiveTime / 1000 // Convert from milliseconds to seconds
	return FormatDuration(seconds)
}

// DPS calculates damage per second
func (p *PlayerEntry) DPS() float64 {
	if p.ActiveTime == 0 {
		return 0
	}
	// ActiveTime is in milliseconds, so divide by 1000 to get seconds
	return p.Total / (float64(p.ActiveTime) / 1000.0)
}

// FormatDPS returns formatted DPS string
func (p *PlayerEntry) FormatDPS() string {
	return FormatNumber(int64(p.DPS()))
}
