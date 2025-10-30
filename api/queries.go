package api

import (
	"fmt"
)

// GraphQL query constants
// These queries use the Warcraft Logs v2 GraphQL API

const (
	// DamageTableQuery fetches damage data for a specific fight
	// Variables needed: $code (report code) and $fightID (fight number)
	DamageTableQuery = `
		query DamageTable($code: String!, $fightID: Int!) {
			reportData {
				report(code: $code) {
					table(fightIDs: [$fightID], dataType: DamageDone)
				}
			}
		}`

	// HealingTableQuery fetches healing data for a specific fight
	HealingTableQuery = `
		query HealingTable($code: String!, $fightID: Int!) {
			reportData {
				report(code: $code) {
					table(fightIDs: [$fightID], dataType: Healing)
				}
			}
		}`

	// MasterDataQuery fetches all players and their information from a report
	// This is used by the players command and for player name → ID mapping
	MasterDataQuery = `
		query MasterData($code: String!) {
			reportData {
				report(code: $code) {
					masterData {
						actors(type: "player") {
							id
							name
							type
							subType
							server
							icon
						}
					}
				}
			}
		}`

	// AllActorsQuery fetches ALL actors (players, NPCs, pets) from a report
	// This includes boss names and enemy names for death analysis
	AllActorsQuery = `
		query AllActors($code: String!) {
			reportData {
				report(code: $code) {
					masterData {
						actors {
							id
							name
							type
							subType
							server
							icon
							gameID
						}
					}
				}
			}
		}`

	// SingleAbilityLookupQuery fetches a single ability name from game data
	// This is used to resolve ability IDs from events to human-readable names
	SingleAbilityLookupQuery = `
		query SingleAbilityLookup($abilityID: Int!) {
			gameData {
				ability(id: $abilityID) {
					id
					name
					icon
				}
			}
		}`

	// FightInfoQuery fetches fight details including start/end times
	FightInfoQuery = `
		query FightInfo($code: String!) {
			reportData {
				report(code: $code) {
					fights {
						id
						name
						encounterID
						startTime
						endTime
						kill
						difficulty
						fightPercentage
					}
				}
			}
		}`

	// DeathEventsQuery fetches death events from the Events API
	// Note: data field is JSON type, so we can't make subselections on it
	// Supports pagination via startTime parameter
	DeathEventsQuery = `
		query DeathEvents($code: String!, $fightID: Int!, $playerID: Int, $startTime: Float) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						targetID: $playerID,
						dataType: Deaths,
						startTime: $startTime,
						limit: 100
					) {
						data
						nextPageTimestamp
					}
				}
			}
		}`

	// DamageTakenBeforeDeathQuery fetches damage taken events before a death
	DamageTakenBeforeDeathQuery = `
		query DamageTakenBeforeDeath($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						targetID: $playerID,
						dataType: DamageTaken,
						startTime: $startTime,
						endTime: $endTime,
						limit: 1000
					) {
						data
						nextPageTimestamp
					}
				}
			}
		}`

	// HealingReceivedBeforeDeathQuery fetches healing events before death
	HealingReceivedBeforeDeathQuery = `
		query HealingReceivedBeforeDeath($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						targetID: $playerID,
						dataType: Healing,
						startTime: $startTime,
						endTime: $endTime,
						limit: 1000
					) {
						data
						nextPageTimestamp
					}
				}
			}
		}`

	// DefensiveAbilitiesBeforeDeathQuery fetches defensive casts before death
	DefensiveAbilitiesBeforeDeathQuery = `
		query DefensiveAbilitiesBeforeDeath($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						sourceID: $playerID,
						dataType: Casts,
						startTime: $startTime,
						endTime: $endTime,
						limit: 1000
					) {
						data
						nextPageTimestamp
					}
				}
			}
		}`

	// InterruptEventsQuery fetches interrupt events from the Events API
	// Supports pagination via startTime parameter
	InterruptEventsQuery = `
		query InterruptEvents($code: String!, $fightID: Int!, $playerID: Int, $startTime: Float) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						sourceID: $playerID,
						dataType: Interrupts,
						startTime: $startTime,
						limit: 10000
					) {
						data
						nextPageTimestamp
					}
				}
			}
		}`
)

// Table Request Functions

// NewTableRequest creates a generic GraphQL request for any table data type
func NewTableRequest(code string, fightID int, dataType DataType) *GraphQLRequest {
	var query string
	switch dataType {
	case DataTypeDamage:
		query = DamageTableQuery
	case DataTypeHealing:
		query = HealingTableQuery
	default:
		query = DamageTableQuery // fallback
	}

	return &GraphQLRequest{
		Query: query,
		Variables: map[string]any{
			"code":    code,
			"fightID": fightID,
		},
	}
}

// Master Data Request Functions

// NewMasterDataRequest creates a GraphQL request for player information
func NewMasterDataRequest(code string) *GraphQLRequest {
	return &GraphQLRequest{
		Query: MasterDataQuery,
		Variables: map[string]any{
			"code": code,
		},
	}
}

// NewAllActorsRequest creates a GraphQL request for all actors (players, NPCs, pets)
func NewAllActorsRequest(code string) *GraphQLRequest {
	return &GraphQLRequest{
		Query: AllActorsQuery,
		Variables: map[string]any{
			"code": code,
		},
	}
}

// Fight Info Request Functions

// NewFightInfoRequest creates a GraphQL request for fight information
func NewFightInfoRequest(code string) *GraphQLRequest {
	return &GraphQLRequest{
		Query: FightInfoQuery,
		Variables: map[string]any{
			"code": code,
		},
	}
}

// Ability Lookup Request Functions

// NewAbilityLookupRequest creates a GraphQL request for ability name lookup
// This queries the gameData API to resolve ability IDs to names
func NewAbilityLookupRequest(abilityID int) *GraphQLRequest {
	return &GraphQLRequest{
		Query: SingleAbilityLookupQuery,
		Variables: map[string]any{
			"abilityID": abilityID,
		},
	}
}

// Event API Request Functions

// NewDeathEventsRequest creates a GraphQL request for death events
// playerID and startTime are optional (pass nil to omit)
func NewDeathEventsRequest(code string, fightID int, playerID *int, startTime *float64) *GraphQLRequest {
	variables := map[string]any{
		"code":    code,
		"fightID": fightID,
	}

	if playerID != nil {
		variables["playerID"] = *playerID
	}

	if startTime != nil {
		variables["startTime"] = *startTime
	}

	return &GraphQLRequest{
		Query:     DeathEventsQuery,
		Variables: variables,
	}
}

// NewDamageTakenRequest creates a GraphQL request for damage taken before death
func NewDamageTakenRequest(code string, fightID int, playerID int, startTime, endTime float64) *GraphQLRequest {
	return &GraphQLRequest{
		Query: DamageTakenBeforeDeathQuery,
		Variables: map[string]any{
			"code":      code,
			"fightID":   fightID,
			"playerID":  playerID,
			"startTime": startTime,
			"endTime":   endTime,
		},
	}
}

// NewHealingReceivedRequest creates a GraphQL request for healing received before death
func NewHealingReceivedRequest(code string, fightID int, playerID int, startTime, endTime float64) *GraphQLRequest {
	return &GraphQLRequest{
		Query: HealingReceivedBeforeDeathQuery,
		Variables: map[string]any{
			"code":      code,
			"fightID":   fightID,
			"playerID":  playerID,
			"startTime": startTime,
			"endTime":   endTime,
		},
	}
}

// NewDefensiveAbilitiesRequest creates a GraphQL request for defensive abilities used before death
func NewDefensiveAbilitiesRequest(code string, fightID int, playerID int, startTime, endTime float64) *GraphQLRequest {
	return &GraphQLRequest{
		Query: DefensiveAbilitiesBeforeDeathQuery,
		Variables: map[string]any{
			"code":      code,
			"fightID":   fightID,
			"playerID":  playerID,
			"startTime": startTime,
			"endTime":   endTime,
		},
	}
}

// Interrupt and Cast Event Request Functions

// NewInterruptEventsRequest creates a GraphQL request for interrupt events
// playerID and startTime are optional (pass nil to omit)
// startTime is used for pagination - pass nextPageTimestamp from previous response
func NewInterruptEventsRequest(code string, fightID int, playerID *int, startTime *float64) *GraphQLRequest {
	variables := map[string]any{
		"code":    code,
		"fightID": fightID,
	}

	if playerID != nil {
		variables["playerID"] = *playerID
	}

	if startTime != nil {
		variables["startTime"] = *startTime
	}

	return &GraphQLRequest{
		Query:     InterruptEventsQuery,
		Variables: variables,
	}
}

// NewCastEventsRequest creates a GraphQL request for cast events with specific ability
// abilityID is optional (pass nil to get all casts)
// startTime is optional (pass nil for first page, use nextPageTimestamp for pagination)
// hostilityType filters by Enemies or Friendlies
func NewCastEventsRequest(code string, fightID int, abilityID *int, hostilityType EventHostilityType, startTime *float64) *GraphQLRequest {
	// Validate hostility type (security: prevent injection)
	if hostilityType != EventHostilityHostile && hostilityType != EventHostilityFriendly && hostilityType != EventHostilityAll {
		hostilityType = EventHostilityHostile // safe default
	}

	// Build query with hostilityType as literal (WCL API requires enum as literal, not variable)
	query := fmt.Sprintf(`
		query CastEvents($code: String!, $fightID: Int!, $abilityID: Int, $startTime: Float) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						abilityID: $abilityID,
						dataType: Casts,
						hostilityType: %s,
						startTime: $startTime,
						limit: 10000
					) {
						data
						nextPageTimestamp
					}
				}
			}
		}`, hostilityType)

	variables := map[string]any{
		"code":    code,
		"fightID": fightID,
	}

	if abilityID != nil {
		variables["abilityID"] = *abilityID
	}

	if startTime != nil {
		variables["startTime"] = *startTime
	}

	return &GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
}

// NewAllCastEventsRequest creates a GraphQL request for all cast events
// startTime is optional (pass nil for first page, use nextPageTimestamp for pagination)
// hostilityType filters by Enemies or Friendlies
func NewAllCastEventsRequest(code string, fightID int, hostilityType EventHostilityType, startTime *float64) *GraphQLRequest {
	// Validate hostility type (security: prevent injection)
	if hostilityType != EventHostilityHostile && hostilityType != EventHostilityFriendly && hostilityType != EventHostilityAll {
		hostilityType = EventHostilityHostile // safe default
	}

	// Build query with hostilityType as literal (WCL API requires enum as literal, not variable)
	query := fmt.Sprintf(`
		query AllCastEvents($code: String!, $fightID: Int!, $startTime: Float) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						dataType: Casts,
						hostilityType: %s,
						startTime: $startTime,
						limit: 10000
					) {
						data
						nextPageTimestamp
					}
				}
			}
		}`, hostilityType)

	variables := map[string]any{
		"code":    code,
		"fightID": fightID,
	}

	if startTime != nil {
		variables["startTime"] = *startTime
	}

	return &GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
}
