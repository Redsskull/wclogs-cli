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

	// DeathsTableQuery fetches death event data for a specific fight
	DeathsTableQuery = `
		query DeathsTable($code: String!, $fightID: Int!) {
			reportData {
				report(code: $code) {
					table(fightIDs: [$fightID], dataType: Deaths)
				}
			}
		}`

	// InterruptsTableQuery fetches interrupt data for a specific fight
	InterruptsTableQuery = `
		query InterruptsTable($code: String!, $fightID: Int!) {
			reportData {
				report(code: $code) {
					table(fightIDs: [$fightID], dataType: Interrupts)
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
)

// NewTableRequest creates a generic GraphQL request for any table data type
func NewTableRequest(code string, fightID int, dataType DataType) *GraphQLRequest {
	var query string
	switch dataType {
	case DataTypeDamage:
		query = DamageTableQuery
	case DataTypeHealing:
		query = HealingTableQuery
	case DataTypeDeaths:
		query = DeathsTableQuery
	case DataTypeInterrupts:
		query = InterruptsTableQuery
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

// NewDamageTableRequest creates a GraphQL request for damage data (backwards compatibility)
func NewDamageTableRequest(code string, fightID int) *GraphQLRequest {
	return NewTableRequest(code, fightID, DataTypeDamage)
}
