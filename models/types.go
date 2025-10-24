package models

import (
	"time"
)

// TableType represents the different types of data we can query
type TableType string

const (
	// Available table data types from the API
	TableTypeDamageDone  TableType = "DamageDone"
	TableTypeHealingDone TableType = "HealingDone"
	TableTypeDeaths      TableType = "Deaths"
	TableTypeInterrupts  TableType = "Interrupts"
	TableTypeDamageTaken TableType = "DamageTaken"
	TableTypeDispels     TableType = "Dispels"
	TableTypeBuffs       TableType = "Buffs"
	TableTypeDebuffs     TableType = "Debuffs"
	TableTypeCasts       TableType = "Casts"
	TableTypeSummons     TableType = "Summons"
)

// NewGraphQLResponse creates a new GraphQLResponse with initialized fields
func NewGraphQLResponse() *GraphQLResponse {
	return &GraphQLResponse{
		Data:   &ResponseData{},
		Errors: make([]GraphQLError, 0),
	}
}

// NewPlayer creates a new Player with the given values
func NewPlayer(name, class string, total float64, icon string) *Player {
	return &Player{
		Name:  name,
		Class: class,
		Total: total,
		Icon:  icon,
	}
}

// NewTableData creates a new TableData with initialized entries slice
func NewTableData() *TableData {
	return &TableData{
		Entries: make([]PlayerEntry, 0),
	}
}

// IsValid checks if the GraphQL response contains valid data
func (r *GraphQLResponse) IsValid() bool {
	return r.Data != nil && len(r.Errors) == 0
}

// HasErrors checks if the GraphQL response contains any errors
func (r *GraphQLResponse) HasErrors() bool {
	return len(r.Errors) > 0
}

// FirstError returns the first error message, or empty string if no errors
func (r *GraphQLResponse) FirstError() string {
	if len(r.Errors) > 0 {
		return r.Errors[0].Message
	}
	return ""
}

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

// PlayerInfo represents a player with their basic information (NEW for Day 6)
type PlayerInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Class  string `json:"class"`
	Server string `json:"server"`
	Icon   string `json:"icon"`
}

// PlayerLookup provides player name → ID mapping functionality (NEW for Day 6)
type PlayerLookup struct {
	playersByName map[string]*PlayerInfo // Name → PlayerInfo mapping
	playersByID   map[int]*PlayerInfo    // ID → PlayerInfo mapping
}

// NewPlayerFromEntry creates a Player from a PlayerEntry (ORIGINAL - KEEP)
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

// Event represents a single combat log event
type Event struct {
	Timestamp float64 `json:"timestamp"`
	Type      string  `json:"type"`
	SourceID  *int    `json:"sourceID"`
	TargetID  *int    `json:"targetID"`
	AbilityID *int    `json:"abilityGameID"`
	Amount    *int    `json:"amount"`
	HitType   *int    `json:"hitType"`
	Overkill  *int    `json:"overkill"`
	Tick      *bool   `json:"tick"`

	// These are available in some event types
	Ability *EventAbility `json:"ability"`
	Source  *EventActor   `json:"source"`
	Target  *EventActor   `json:"target"`
}

// EventAbility represents ability information in events
type EventAbility struct {
	Name   string `json:"name"`
	GameID int    `json:"gameID"`
	Type   int    `json:"type"`
	Icon   string `json:"icon"`
}

// EventActor represents source/target information in events
type EventActor struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	Type string `json:"type"`
	Icon string `json:"icon"`
}

// DeathEvent represents a death event with parsed details
type DeathEvent struct {
	PlayerID             int
	PlayerName           string
	Timestamp            float64
	TimeFromStart        time.Duration
	KillingAbility       *EventAbility
	KillingSource        *EventActor
	Overkill             int
	DamageLeadingToDeath []*DamageEvent
}

// DamageEvent represents damage taken before death
type DamageEvent struct {
	Timestamp     float64
	TimeFromStart time.Duration
	Ability       *EventAbility
	Source        *EventActor
	Amount        int
	HitType       int
	Tick          bool
}

// InterruptEvent represents an interrupt event
type InterruptEvent struct {
	PlayerID           int
	PlayerName         string
	Timestamp          float64
	TimeFromStart      time.Duration
	Ability            *EventAbility
	Target             *EventActor
	InterruptedAbility *EventAbility
}

// InterruptAnalysis represents interrupt statistics for a player
type InterruptAnalysis struct {
	PlayerID             int
	PlayerName           string
	SuccessfulInterrupts int
	MissedOpportunities  int
	InterruptSuccess     float64
	InterruptDetails     []*InterruptEvent
}

// DeathAnalysis represents death analysis for a player
type DeathAnalysis struct {
	PlayerID               int
	PlayerName             string
	Deaths                 []*DeathEvent
	TotalDeaths            int
	AverageSurvivalTime    time.Duration
	AverageDamagePerSecond float64
	TopCausesOfDeath       []string
}
