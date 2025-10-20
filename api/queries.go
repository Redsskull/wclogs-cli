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
				table(fightIDs: [$fightID], dataType: HealingDone)

			}
		}`

	// I'll add more query types later (deaths, interrupts, etc.)
)

// NewDamageTableRequest creates a GraphQL request for damage data
func NewDamageTableRequest(code string, fightID int) *GraphQLRequest {
	return &GraphQLRequest{
		Query: DamageTableQuery,
		Variables: map[string]any{
			"code":    code,
			"fightID": fightID,
		},
	}
}
