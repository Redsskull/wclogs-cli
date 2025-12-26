package models

import (
	"fmt"
	"sort"
	"time"
)

// InterruptAnalysisResult represents the complete interrupt analysis for a fight
type InterruptAnalysisResult struct {
	PlayerSummary map[string]*PlayerInterruptSummary `json:"player_summary"`
	SpellAnalysis map[string]*SpellInterruptAnalysis `json:"spell_analysis"`
	FightInfo     *FightInfo                         `json:"fight_info"`
	TotalCasts    int                                `json:"total_casts"`
	TotalStops    int                                `json:"total_stops"`
	OverallRate   float64                            `json:"overall_rate"`
}

// PlayerInterruptSummary represents a player's interrupt performance
type PlayerInterruptSummary struct {
	PlayerName            string         `json:"player_name"`
	PlayerID              int            `json:"player_id"`
	TotalStops            int            `json:"total_stops"`
	MissedCasts           int            `json:"missed_casts"`
	SuccessRate           float64        `json:"success_rate"`
	SpellsStoppedDetailed map[string]int `json:"spells_stopped_detailed"`
}

// SpellInterruptAnalysis represents interrupt analysis for a specific spell
type SpellInterruptAnalysis struct {
	SpellName     string              `json:"spell_name"`
	SpellID       int                 `json:"spell_id"`
	TotalCasts    int                 `json:"total_casts"`
	StoppedCasts  int                 `json:"stopped_casts"`
	MissedCasts   int                 `json:"missed_casts"`
	StopRate      float64             `json:"stop_rate"`
	StoppedBy     []*InterruptDetail  `json:"stopped_by"`
	MissedDetails []*MissedCastDetail `json:"missed_details"`
}

// InterruptDetail represents a successful interrupt
type InterruptDetail struct {
	PlayerName     string    `json:"player_name"`
	PlayerID       int       `json:"player_id"`
	Count          int       `json:"count"`
	PercentOfCasts float64   `json:"percent_of_casts"`
	Timestamps     []float64 `json:"timestamps"`
	Targets        []string  `json:"targets"` // Who was being interrupted
}

// MissedCastDetail represents a cast that completed (wasn't interrupted)
type MissedCastDetail struct {
	CasterName           string  `json:"caster_name"`
	CasterID             int     `json:"caster_id"`
	TargetName           string  `json:"target_name"`
	TargetID             int     `json:"target_id"`
	Timestamp            float64 `json:"timestamp"`
	TimeInFight          string  `json:"time_in_fight"`
	CouldHaveBeenStopped bool    `json:"could_have_been_stopped"`
}

// BeginCastEvent represents the start of a spell cast
type BeginCastEvent struct {
	CasterID   int     `json:"caster_id"`
	CasterName string  `json:"caster_name"`
	TargetID   *int    `json:"target_id"`
	TargetName string  `json:"target_name"`
	SpellID    int     `json:"spell_id"`
	SpellName  string  `json:"spell_name"`
	Timestamp  float64 `json:"timestamp"`
	Duration   float64 `json:"duration"` // Cast time in milliseconds
}

// InterruptEvent represents a successful interrupt
type InterruptEventDetail struct {
	InterrupterID        int     `json:"interrupter_id"`
	InterrupterName      string  `json:"interrupter_name"`
	TargetID             int     `json:"target_id"`
	TargetName           string  `json:"target_name"`
	SpellID              int     `json:"spell_id"`
	SpellName            string  `json:"spell_name"`
	Timestamp            float64 `json:"timestamp"`
	InterruptSpellID     int     `json:"interrupt_spell_id"`
	InterruptSpellName   string  `json:"interrupt_spell_name"`
	InterruptedSpellID   int     `json:"interrupted_spell_id"`   // The spell that was interrupted
	InterruptedSpellName string  `json:"interrupted_spell_name"` // Name of interrupted spell
}

// CastCompleteEvent represents a spell that finished casting
type CastCompleteEvent struct {
	CasterID   int     `json:"caster_id"`
	CasterName string  `json:"caster_name"`
	TargetID   *int    `json:"target_id"`
	TargetName string  `json:"target_name"`
	SpellID    int     `json:"spell_id"`
	SpellName  string  `json:"spell_name"`
	Timestamp  float64 `json:"timestamp"`
}

// FightInfo contains basic fight information for context
type FightInfo struct {
	FightID   int           `json:"fight_id"`
	Name      string        `json:"name"`
	Duration  time.Duration `json:"duration"`
	StartTime int64         `json:"start_time"`
	EndTime   int64         `json:"end_time"`
}

// NewInterruptAnalysisResult creates a new interrupt analysis result
func NewInterruptAnalysisResult(fightInfo *FightInfo) *InterruptAnalysisResult {
	return &InterruptAnalysisResult{
		PlayerSummary: make(map[string]*PlayerInterruptSummary),
		SpellAnalysis: make(map[string]*SpellInterruptAnalysis),
		FightInfo:     fightInfo,
	}
}

// AddPlayerInterrupt adds a successful interrupt to the analysis
func (iar *InterruptAnalysisResult) AddPlayerInterrupt(playerName string, playerID int, spellName string, spellID int, timestamp float64, targetName string) {
	// Initialize player summary if needed
	if iar.PlayerSummary[playerName] == nil {
		iar.PlayerSummary[playerName] = &PlayerInterruptSummary{
			PlayerName:            playerName,
			PlayerID:              playerID,
			SpellsStoppedDetailed: make(map[string]int),
		}
	}

	// Update player totals
	player := iar.PlayerSummary[playerName]
	player.TotalStops++
	player.SpellsStoppedDetailed[spellName]++

	// Initialize spell analysis if needed
	if iar.SpellAnalysis[spellName] == nil {
		iar.SpellAnalysis[spellName] = &SpellInterruptAnalysis{
			SpellName:     spellName,
			SpellID:       spellID,
			StoppedBy:     make([]*InterruptDetail, 0),
			MissedDetails: make([]*MissedCastDetail, 0),
		}
	}

	// Update spell analysis
	spell := iar.SpellAnalysis[spellName]
	spell.StoppedCasts++

	// Find or create interrupt detail for this player
	var detail *InterruptDetail
	for _, d := range spell.StoppedBy {
		if d.PlayerName == playerName {
			detail = d
			break
		}
	}

	if detail == nil {
		detail = &InterruptDetail{
			PlayerName: playerName,
			PlayerID:   playerID,
			Timestamps: make([]float64, 0),
			Targets:    make([]string, 0),
		}
		spell.StoppedBy = append(spell.StoppedBy, detail)
	}

	detail.Count++
	detail.Timestamps = append(detail.Timestamps, timestamp)
	detail.Targets = append(detail.Targets, targetName)

	iar.TotalStops++
}

// AddMissedCast adds a cast that completed without being interrupted
func (iar *InterruptAnalysisResult) AddMissedCast(casterName string, casterID int, targetName string, targetID int, spellName string, spellID int, timestamp float64, fightStartTime int64) {
	// Initialize spell analysis if needed
	if iar.SpellAnalysis[spellName] == nil {
		iar.SpellAnalysis[spellName] = &SpellInterruptAnalysis{
			SpellName:     spellName,
			SpellID:       spellID,
			StoppedBy:     make([]*InterruptDetail, 0),
			MissedDetails: make([]*MissedCastDetail, 0),
		}
	}

	spell := iar.SpellAnalysis[spellName]
	spell.MissedCasts++

	// Calculate time in fight
	timeInFight := time.Duration((timestamp - float64(fightStartTime)) * float64(time.Millisecond))

	missedDetail := &MissedCastDetail{
		CasterName:           casterName,
		CasterID:             casterID,
		TargetName:           targetName,
		TargetID:             targetID,
		Timestamp:            timestamp,
		TimeInFight:          formatDuration(timeInFight),
		CouldHaveBeenStopped: true, // Assume it could have been stopped for now
	}

	spell.MissedDetails = append(spell.MissedDetails, missedDetail)
	iar.TotalCasts++
}

// CalculateRates calculates success rates for all players and spells
func (iar *InterruptAnalysisResult) CalculateRates() {
	// Calculate player success rates
	for _, player := range iar.PlayerSummary {
		total := player.TotalStops + player.MissedCasts
		if total > 0 {
			player.SuccessRate = float64(player.TotalStops) / float64(total) * 100
		}
	}

	// Calculate spell stop rates and per-player percentages
	for _, spell := range iar.SpellAnalysis {
		spell.TotalCasts = spell.StoppedCasts + spell.MissedCasts
		if spell.TotalCasts > 0 {
			spell.StopRate = float64(spell.StoppedCasts) / float64(spell.TotalCasts) * 100

			// Calculate percentage for each player who stopped this spell
			for _, detail := range spell.StoppedBy {
				detail.PercentOfCasts = float64(detail.Count) / float64(spell.TotalCasts) * 100
			}
		}

		// Sort stopped by details by count (descending)
		sort.Slice(spell.StoppedBy, func(i, j int) bool {
			return spell.StoppedBy[i].Count > spell.StoppedBy[j].Count
		})

		// Sort missed details by timestamp
		sort.Slice(spell.MissedDetails, func(i, j int) bool {
			return spell.MissedDetails[i].Timestamp < spell.MissedDetails[j].Timestamp
		})
	}

	// Calculate overall rate
	iar.TotalCasts = iar.TotalStops + iar.countTotalMissed()
	if iar.TotalCasts > 0 {
		iar.OverallRate = float64(iar.TotalStops) / float64(iar.TotalCasts) * 100
	}
}

// countTotalMissed counts total missed casts across all spells
func (iar *InterruptAnalysisResult) countTotalMissed() int {
	total := 0
	for _, spell := range iar.SpellAnalysis {
		total += spell.MissedCasts
	}
	return total
}

// GetTopPlayers returns players sorted by total stops (descending)
func (iar *InterruptAnalysisResult) GetTopPlayers(limit int) []*PlayerInterruptSummary {
	players := make([]*PlayerInterruptSummary, 0, len(iar.PlayerSummary))
	for _, player := range iar.PlayerSummary {
		players = append(players, player)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].TotalStops > players[j].TotalStops
	})

	if limit > 0 && len(players) > limit {
		return players[:limit]
	}
	return players
}

// GetSpellsByActivity returns spells sorted by total casts (descending)
func (iar *InterruptAnalysisResult) GetSpellsByActivity() []*SpellInterruptAnalysis {
	spells := make([]*SpellInterruptAnalysis, 0, len(iar.SpellAnalysis))
	for _, spell := range iar.SpellAnalysis {
		spells = append(spells, spell)
	}

	sort.Slice(spells, func(i, j int) bool {
		return spells[i].TotalCasts > spells[j].TotalCasts
	})

	return spells
}

// formatDuration formats a duration as MM:SS
func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// IsInterruptableSpell determines if a spell should be considered for interrupt analysis
// This filters out non-interruptable spells like instant abilities
func IsInterruptableSpell(spellName string, duration float64) bool {
	// Filter out obvious non-interruptable spells
	nonInterruptable := map[string]bool{
		"Melee":       true,
		"Auto Attack": true,
		"Arcane Orb":  true,  // Example instant spell
		"Shadow Bolt": false, // This is typically interruptable
		// Add more as needed based on actual spell data
	}

	if excluded, exists := nonInterruptable[spellName]; exists {
		return !excluded
	}

	// Generally, spells with cast time > 0 are interruptable
	// Instant spells (duration = 0) are not interruptable
	return duration > 0
}

// ValidateAnalysis performs basic validation on the analysis results
func (iar *InterruptAnalysisResult) ValidateAnalysis() []string {
	var warnings []string

	if iar.TotalCasts == 0 {
		warnings = append(warnings, "No interruptable casts found in this fight")
	}

	if iar.TotalStops == 0 {
		warnings = append(warnings, "No successful interrupts found in this fight")
	}

	if len(iar.PlayerSummary) == 0 {
		warnings = append(warnings, "No players with interrupt capabilities detected")
	}

	if len(iar.SpellAnalysis) == 0 {
		warnings = append(warnings, "No interruptable spells detected")
	}

	return warnings
}
