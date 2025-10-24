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

// Legacy functions for backwards compatibility
func NewDamageTableRequest(code string, fightID int) *GraphQLRequest {
	return NewTableRequest(code, fightID, DataTypeDamage)
}
