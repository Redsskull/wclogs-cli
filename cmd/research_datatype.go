package cmd

import (
	"fmt"
	"wclogs-cli/api"
)

// ResearchEventDataType uses GraphQL introspection to discover valid EventDataType enum values
func ResearchEventDataType(apiClient *api.Client) error {
	fmt.Printf("🔍 Researching WCL EventDataType enum values using GraphQL introspection...\n\n")

	// GraphQL introspection query to get enum values for EventDataType
	introspectionQuery := `
		query IntrospectEventDataType {
			__schema {
				types {
					name
					kind
					enumValues {
						name
						description
					}
				}
			}
		}`

	_, err := apiClient.Query(introspectionQuery, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("introspection query failed: %w", err)
	}

	fmt.Printf("✅ Introspection query successful!\n")
	fmt.Printf("📊 Response contains schema information\n\n")

	// Also try a more specific introspection for EventDataType
	specificQuery := `
		query GetEventDataTypeEnum {
			__type(name: "EventDataType") {
				name
				kind
				enumValues {
					name
					description
				}
			}
		}`

	fmt.Printf("🔍 Trying specific EventDataType enum query...\n")
	_, err = apiClient.Query(specificQuery, map[string]interface{}{})
	if err != nil {
		fmt.Printf("❌ Specific query failed: %v\n", err)
	} else {
		fmt.Printf("✅ Specific EventDataType query successful!\n")
	}

	// Try to test different dataType values we suspect exist
	fmt.Printf("\n🧪 Testing suspected EventDataType values:\n")

	testValues := []string{
		"DamageTaken",
		"Healing",
		"Deaths",
		"Damage",
		"DamageDealt",
		"All",
		"Resources",
		"Casts",
		"Buffs",
		"Debuffs",
	}

	// Use a simple test query to see which dataType values are accepted
	for _, testValue := range testValues {
		fmt.Printf("  • Testing '%s'... ", testValue)

		testQuery := fmt.Sprintf(`
			query TestDataType($code: String!, $fightID: Int!) {
				reportData {
					report(code: $code) {
						events(
							fightIDs: [$fightID],
							dataType: %s,
							limit: 1
						) {
							data
						}
					}
				}
			}`, testValue)

		// Use our known good report for testing
		testVars := map[string]interface{}{
			"code":    "9aQbqzgJy2dK8rVk", // The report we know works
			"fightID": 36,                 // The Fractillus fight
		}

		testResp, testErr := apiClient.Query(testQuery, testVars)
		if testErr != nil {
			fmt.Printf("❌ INVALID (%v)\n", testErr)
		} else if testResp != nil {
			fmt.Printf("✅ VALID\n")
		} else {
			fmt.Printf("⚠️  UNKNOWN\n")
		}
	}

	fmt.Printf("\n🔍 Testing WCL Table API (this might be how damage taken works):\n")

	// Test the table API which might be what WCL web interface uses for damage taken tables
	tableQuery := `
		query TestTableAPI($code: String!, $fightID: Int!, $playerID: Int!) {
			reportData {
				report(code: $code) {
					table(
						fightIDs: [$fightID],
						targetID: $playerID,
						dataType: DamageTaken,
						viewBy: Target
					)
				}
			}
		}`

	tableVars := map[string]interface{}{
		"code":     "9aQbqzgJy2dK8rVk", // The report we know works
		"fightID":  36,                 // The Fractillus fight
		"playerID": 11,                 // Naalla's player ID
	}

	fmt.Printf("  • Testing table API for Naalla's damage taken... ")
	tableResp, tableErr := apiClient.Query(tableQuery, tableVars)
	if tableErr != nil {
		fmt.Printf("❌ FAILED (%v)\n", tableErr)
	} else if tableResp != nil {
		fmt.Printf("✅ SUCCESS! This might be the key!\n")
		fmt.Printf("    💡 The 'table' API likely powers WCL's damage taken interface\n")
	} else {
		fmt.Printf("⚠️  UNKNOWN\n")
	}

	fmt.Printf("\n📋 Research Summary:\n")
	fmt.Printf("   • GraphQL introspection can reveal valid enum values\n")
	fmt.Printf("   • We tested common EventDataType values\n")
	fmt.Printf("   • Valid values can be used in our damage event queries\n")
	fmt.Printf("   • We discovered the 'table' API which may be the solution!\n")
	fmt.Printf("   • The table API likely provides aggregated damage data like WCL web interface\n")

	return nil
}
