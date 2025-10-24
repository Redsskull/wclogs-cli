package api

// GraphQL query constants
// I define these as constants so they're reusable and easy to maintain

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

	// AbilityLookupQuery fetches ability names from game data
	AbilityLookupQuery = `
		query AbilityLookup($abilityIDs: [Int!]!) {
			gameData {
				abilities: [
					# We need to query each ability individually
					# This will be constructed dynamically per ability ID
				]
			}
		}`

	// SingleAbilityLookupQuery fetches a single ability name
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
	DeathEventsQuery = `
		query DeathEvents($code: String!, $fightID: Int!, $playerID: Int) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						targetID: $playerID,
						dataType: Deaths,
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
	InterruptEventsQuery = `
		query InterruptEvents($code: String!, $fightID: Int!, $playerID: Int) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						sourceID: $playerID,
						dataType: Interrupts,
						limit: 100
					) {
						data
						nextPageTimestamp
					}
				}
			}
		}`

	// TestEventsQuery is a simple query to test the Events API structure
	TestEventsQuery = `
		query TestEvents($code: String!, $fightID: Int!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						limit: 10
					) {
						data
					}
				}
			}
		}`
)

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

// NewMasterDataRequest creates a GraphQL request for player information
func NewMasterDataRequest(code string) *GraphQLRequest {
	return &GraphQLRequest{
		Query: MasterDataQuery,
		Variables: map[string]any{
			"code": code,
		},
	}
}

// NewFightInfoRequest creates a GraphQL request for fight information
func NewFightInfoRequest(code string) *GraphQLRequest {
	return &GraphQLRequest{
		Query: FightInfoQuery,
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

// NewAbilityLookupRequest creates a GraphQL request for ability name lookup
func NewAbilityLookupRequest(abilityID int) *GraphQLRequest {
	return &GraphQLRequest{
		Query: SingleAbilityLookupQuery,
		Variables: map[string]any{
			"abilityID": abilityID,
		},
	}
}

// Event API request functions

// NewDeathEventsRequest creates a GraphQL request for death events
func NewDeathEventsRequest(code string, fightID int, playerID *int) *GraphQLRequest {
	variables := map[string]any{
		"code":    code,
		"fightID": fightID,
	}

	if playerID != nil {
		variables["playerID"] = *playerID
	}

	return &GraphQLRequest{
		Query:     DeathEventsQuery,
		Variables: variables,
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

// NewTestEventsRequest creates a simple test query for the Events API
func NewTestEventsRequest(code string, fightID int) *GraphQLRequest {
	return &GraphQLRequest{
		Query: TestEventsQuery,
		Variables: map[string]any{
			"code":    code,
			"fightID": fightID,
		},
	}
}
