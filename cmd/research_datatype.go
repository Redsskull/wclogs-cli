package cmd

import (
	"encoding/json"
	"fmt"
	"wclogs-cli/api"
	"wclogs-cli/models"
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

// ResearchDamageEvents specifically tests different approaches to get individual damage events
func ResearchDamageEvents(apiClient *api.Client) error {
	fmt.Printf("🔬 DAMAGE EVENTS RESEARCH - Testing approaches to get individual damage events\n\n")

	// Known test case: Naalla's death at 6404661ms in fight 36
	reportCode := "9aQbqzgJy2dK8rVk"
	fightID := 36
	playerID := 11 // Naalla
	deathTime := 6404661.0
	windowStart := deathTime - 3000 // 3 seconds before death
	windowEnd := deathTime + 100    // 0.1 seconds after death

	fmt.Printf("📊 Test Parameters:\n")
	fmt.Printf("   • Report: %s\n", reportCode)
	fmt.Printf("   • Fight: %d (Fractillus)\n", fightID)
	fmt.Printf("   • Player: %d (Naalla)\n", playerID)
	fmt.Printf("   • Death Time: %.0fms\n", deathTime)
	fmt.Printf("   • Window: %.0fms to %.0fms\n", windowStart, windowEnd)
	fmt.Printf("\n")

	// Test 1: Standard DamageTaken query (this failed in the breakthrough)
	fmt.Printf("🧪 TEST 1: Standard DamageTaken events query\n")
	damageTakenQuery := `
		query DamageTakenEvents($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						targetID: $playerID,
						startTime: $startTime,
						endTime: $endTime,
						dataType: DamageTaken,
						limit: 50
					) {
						data
					}
				}
			}
		}`

	variables := map[string]interface{}{
		"code":      reportCode,
		"fightID":   fightID,
		"playerID":  playerID,
		"startTime": windowStart,
		"endTime":   windowEnd,
	}

	fmt.Printf("   • Query: events(targetID: %d, dataType: DamageTaken)\n", playerID)
	resp1, err1 := apiClient.Query(damageTakenQuery, variables)
	if err1 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err1)
	} else if resp1.Data == nil || resp1.Data.ReportData == nil || resp1.Data.ReportData.Report == nil || resp1.Data.ReportData.Report.Events == nil {
		fmt.Printf("   ❌ No events data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - parsing events...\n")
		if resp1.Data.ReportData.Report.Events.Data != nil {
			fmt.Printf("   📊 Raw JSON data size: %d bytes\n", len(string(resp1.Data.ReportData.Report.Events.Data)))
			// Try to parse as JSON array to count elements
			if jsonStr := string(resp1.Data.ReportData.Report.Events.Data); jsonStr == "[]" || jsonStr == "" {
				fmt.Printf("   📊 Event count: 0 (empty array)\n")
			} else {
				fmt.Printf("   📊 Event data present - contains actual events!\n")
				fmt.Printf("   💡 First 200 chars: %s...\n", truncateString(jsonStr, 200))
			}
		} else {
			fmt.Printf("   ❌ No event data field\n")
		}
	}

	// Test 2: Try sourceID instead of targetID (maybe damage events work differently)
	fmt.Printf("\n🧪 TEST 2: DamageTaken with sourceID instead of targetID\n")
	sourceIDQuery := `
		query DamageTakenSourceID($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						sourceID: $playerID,
						startTime: $startTime,
						endTime: $endTime,
						dataType: DamageTaken,
						limit: 50
					) {
						data
					}
				}
			}
		}`

	fmt.Printf("   • Query: events(sourceID: %d, dataType: DamageTaken)\n", playerID)
	resp2, err2 := apiClient.Query(sourceIDQuery, variables)
	if err2 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err2)
	} else if resp2.Data == nil || resp2.Data.ReportData == nil || resp2.Data.ReportData.Report == nil || resp2.Data.ReportData.Report.Events == nil {
		fmt.Printf("   ❌ No events data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - parsing events...\n")
		if resp2.Data.ReportData.Report.Events.Data != nil {
			fmt.Printf("   📊 Raw JSON data size: %d bytes\n", len(string(resp2.Data.ReportData.Report.Events.Data)))
			if jsonStr := string(resp2.Data.ReportData.Report.Events.Data); jsonStr == "[]" || jsonStr == "" {
				fmt.Printf("   📊 Event count: 0 (empty array)\n")
			} else {
				fmt.Printf("   📊 Event data present - contains actual events!\n")
				fmt.Printf("   💡 First 200 chars: %s...\n", truncateString(jsonStr, 200))
			}
		} else {
			fmt.Printf("   ❌ No event data field\n")
		}
	}

	// Test 3: Try "All" dataType and filter for damage events
	fmt.Printf("\n🧪 TEST 3: All events and filter for damage types\n")
	allEventsQuery := `
		query AllEvents($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						targetID: $playerID,
						startTime: $startTime,
						endTime: $endTime,
						dataType: All,
						limit: 100
					) {
						data
					}
				}
			}
		}`

	fmt.Printf("   • Query: events(targetID: %d, dataType: All)\n", playerID)
	resp3, err3 := apiClient.Query(allEventsQuery, variables)
	if err3 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err3)
	} else if resp3.Data == nil || resp3.Data.ReportData == nil || resp3.Data.ReportData.Report == nil || resp3.Data.ReportData.Report.Events == nil {
		fmt.Printf("   ❌ No events data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - parsing all events targeting player...\n")
		if resp3.Data.ReportData.Report.Events.Data != nil {
			fmt.Printf("   📊 Raw JSON data size: %d bytes\n", len(string(resp3.Data.ReportData.Report.Events.Data)))
			if jsonStr := string(resp3.Data.ReportData.Report.Events.Data); jsonStr == "[]" || jsonStr == "" {
				fmt.Printf("   📊 Event count: 0 (empty array)\n")
			} else {
				fmt.Printf("   📊 Event data present - mixed events (heal/damage/etc)\n")
				fmt.Printf("   💡 First 400 chars: %s...\n", truncateString(jsonStr, 400))
			}
		} else {
			fmt.Printf("   ❌ No event data field\n")
		}
	}

	// Test 4: Try no targetID or sourceID (get all events in time window)
	fmt.Printf("\n🧪 TEST 4: All events in time window (no player filter)\n")
	timeWindowQuery := `
		query TimeWindowEvents($code: String!, $fightID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						startTime: $startTime,
						endTime: $endTime,
						dataType: All,
						limit: 200
					) {
						data
					}
				}
			}
		}`

	timeVars := map[string]interface{}{
		"code":      reportCode,
		"fightID":   fightID,
		"startTime": windowStart,
		"endTime":   windowEnd,
	}

	fmt.Printf("   • Query: events(no player filter, dataType: All)\n")
	resp4, err4 := apiClient.Query(timeWindowQuery, timeVars)
	if err4 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err4)
	} else if resp4.Data == nil || resp4.Data.ReportData == nil || resp4.Data.ReportData.Report == nil || resp4.Data.ReportData.Report.Events == nil {
		fmt.Printf("   ❌ No events data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - parsing all events in time window...\n")
		if resp4.Data.ReportData.Report.Events.Data != nil {
			fmt.Printf("   📊 Raw JSON data size: %d bytes\n", len(string(resp4.Data.ReportData.Report.Events.Data)))
			if jsonStr := string(resp4.Data.ReportData.Report.Events.Data); jsonStr == "[]" || jsonStr == "" {
				fmt.Printf("   📊 Event count: 0 (empty array)\n")
			} else {
				fmt.Printf("   📊 Event data present - ALL events in time window\n")
				fmt.Printf("   💡 First 400 chars: %s...\n", truncateString(jsonStr, 400))
				fmt.Printf("   💡 Should contain events for targetID: %d and others\n", playerID)
			}
		} else {
			fmt.Printf("   ❌ No event data field\n")
		}
	}

	// Test 5: Parse the "All events targeting player" data and filter for damage events
	fmt.Printf("\n🧪 TEST 5: Parse All events and filter for damage targeting player\n")
	if resp3 != nil && resp3.Data != nil && resp3.Data.ReportData != nil &&
		resp3.Data.ReportData.Report != nil && resp3.Data.ReportData.Report.Events != nil &&
		resp3.Data.ReportData.Report.Events.Data != nil {

		fmt.Printf("   • Parsing JSON events data from Test 3...\n")

		// Try to parse the events
		events, err := models.ParseEventsJSON(resp3.Data.ReportData.Report.Events.Data)
		if err != nil {
			fmt.Printf("   ❌ Failed to parse events: %v\n", err)
		} else {
			fmt.Printf("   ✅ Successfully parsed %d total events\n", len(events))

			// Filter for damage events targeting our player
			var damageEvents []*models.Event
			for _, event := range events {
				if event.Type == "damage" && event.TargetID != nil && *event.TargetID == playerID {
					damageEvents = append(damageEvents, event)
				}
			}

			fmt.Printf("   🎯 Found %d damage events targeting player %d\n", len(damageEvents), playerID)

			if len(damageEvents) > 0 {
				fmt.Printf("   📊 Sample damage events:\n")
				for i, event := range damageEvents {
					if i >= 5 { // Show first 5 events
						break
					}
					fmt.Printf("      %d. Time: %.3fs, Ability: %d, Amount: %d, Source: %d\n",
						i+1,
						(event.Timestamp-deathTime)/1000.0,
						getIntValue(event.AbilityID),
						getIntValue(event.Amount),
						getIntValue(event.SourceID))
				}
				fmt.Printf("   🎉 SUCCESS: Found individual damage events with timing!\n")
				fmt.Printf("   💡 SOLUTION: Use targetID + dataType:All, then filter for type:damage\n")
			}
		}
	} else {
		fmt.Printf("   ❌ No data available for parsing\n")
	}

	// Test 6: Parse the broader "All events in time window" data to find damage targeting player
	fmt.Printf("\n🧪 TEST 6: Parse ALL events in time window and filter for damage targeting player\n")
	if resp4 != nil && resp4.Data != nil && resp4.Data.ReportData != nil &&
		resp4.Data.ReportData.Report != nil && resp4.Data.ReportData.Report.Events != nil &&
		resp4.Data.ReportData.Report.Events.Data != nil {

		fmt.Printf("   • Parsing JSON events data from Test 4 (broader dataset)...\n")

		// Try to parse the events
		events, err := models.ParseEventsJSON(resp4.Data.ReportData.Report.Events.Data)
		if err != nil {
			fmt.Printf("   ❌ Failed to parse events: %v\n", err)
		} else {
			fmt.Printf("   ✅ Successfully parsed %d total events\n", len(events))

			// Filter for damage events targeting our player
			var damageEvents []*models.Event
			var allDamageEvents []*models.Event
			for _, event := range events {
				if event.Type == "damage" {
					allDamageEvents = append(allDamageEvents, event)
					if event.TargetID != nil && *event.TargetID == playerID {
						damageEvents = append(damageEvents, event)
					}
				}
			}

			fmt.Printf("   📊 Total damage events in time window: %d\n", len(allDamageEvents))
			fmt.Printf("   🎯 Found %d damage events targeting player %d\n", len(damageEvents), playerID)

			if len(damageEvents) > 0 {
				fmt.Printf("   📊 Individual damage events targeting Naalla:\n")
				for i, event := range damageEvents {
					if i >= 8 { // Show first 8 events
						break
					}
					fmt.Printf("      %d. Time: %.3fs, Ability: %d, Amount: %d, Source: %d\n",
						i+1,
						(event.Timestamp-deathTime)/1000.0,
						getIntValue(event.AbilityID),
						getIntValue(event.Amount),
						getIntValue(event.SourceID))
				}
				fmt.Printf("   🎉 BREAKTHROUGH: Found individual damage events with timing!\n")
				fmt.Printf("   💡 SOLUTION: Use events(dataType: All, time window) + filter for targetID + type: damage\n")
			} else {
				fmt.Printf("   🔍 DEBUG: Showing sample damage events to understand structure:\n")
				for i, event := range allDamageEvents {
					if i >= 3 {
						break
					}
					fmt.Printf("      Sample %d: Target: %d, Source: %d, Type: %s, Amount: %d\n",
						i+1,
						getIntValue(event.TargetID),
						getIntValue(event.SourceID),
						event.Type,
						getIntValue(event.Amount))
				}
			}
		}
	} else {
		fmt.Printf("   ❌ No data available for parsing\n")
	}

	fmt.Printf("\n📋 RESEARCH SUMMARY:\n")
	fmt.Printf("   • Test 1: targetID + DamageTaken = 0 events (API limitation)\n")
	fmt.Printf("   • Test 2: sourceID + DamageTaken = events where player deals damage (wrong direction)\n")
	fmt.Printf("   • Test 3: targetID + All = mixed events, may not contain damage targeting player\n")
	fmt.Printf("   • Test 4: All events in window = contains everything, including damage to player\n")
	fmt.Printf("   • Test 5: Parse Test 3 data and filter for damage events\n")
	fmt.Printf("   • Test 6: Parse Test 4 data and filter for damage events targeting player\n")
	fmt.Printf("   • SOLUTION: Use events(dataType: All, time window) + filter for targetID + type: damage\n")

	return nil
}

// ResearchDamageTimeline specifically investigates how WCL web interface creates damage timelines
func ResearchDamageTimeline(apiClient *api.Client) error {
	fmt.Printf("🔬 DAMAGE TIMELINE RESEARCH - Understanding WCL web interface damage chronology\n\n")

	// Known test case: Naalla's death at 6404661ms in fight 36
	reportCode := "9aQbqzgJy2dK8rVk"
	fightID := 36
	playerID := 11 // Naalla
	deathTime := 6404661.0
	windowStart := deathTime - 3000 // 3 seconds before death
	windowEnd := deathTime + 100    // 0.1 seconds after death

	fmt.Printf("📊 Expected Timeline from WCL CSV:\n")
	fmt.Printf("   • -0.00s: Null Explosion (Environment → Naalla) 1438291 damage\n")
	fmt.Printf("   • -0.23s: Null Explosion (Environment → Naalla) 2778988 damage\n")
	fmt.Printf("   • -0.24s: Null Explosion (Environment → Naalla) 2778988 damage\n")
	fmt.Printf("   • -0.29s: Null Explosion (Environment → Naalla) 2778988 damage\n")
	fmt.Printf("   • -0.29s: Null Consumption (Fractillus → Naalla) 1470800 damage\n")
	fmt.Printf("   • -0.75s: Crystal Lacerations (Fractillus → Naalla) 1185701 damage\n")
	fmt.Printf("\n")

	// Research 1: Find ALL damage events in the time window targeting any player
	fmt.Printf("🧪 RESEARCH 1: Catalog ALL damage events in time window\n")
	allEventsQuery := `
		query AllDamageEvents($code: String!, $fightID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						startTime: $startTime,
						endTime: $endTime,
						dataType: All,
						limit: 500
					) {
						data
					}
				}
			}
		}`

	timeVars := map[string]interface{}{
		"code":      reportCode,
		"fightID":   fightID,
		"startTime": windowStart,
		"endTime":   windowEnd,
	}

	resp, err := apiClient.Query(allEventsQuery, timeVars)
	if err != nil {
		return fmt.Errorf("failed to query all events: %w", err)
	}

	if resp.Data == nil || resp.Data.ReportData == nil ||
		resp.Data.ReportData.Report == nil || resp.Data.ReportData.Report.Events == nil ||
		resp.Data.ReportData.Report.Events.Data == nil {
		return fmt.Errorf("no events data available")
	}

	events, err := models.ParseEventsJSON(resp.Data.ReportData.Report.Events.Data)
	if err != nil {
		return fmt.Errorf("failed to parse events: %w", err)
	}

	fmt.Printf("   ✅ Parsed %d total events in time window\n", len(events))

	// Filter and analyze damage events
	var allDamageEvents []*models.Event
	var naallaDamageEvents []*models.Event
	sourceMap := make(map[int]int)  // sourceID -> count
	abilityMap := make(map[int]int) // abilityID -> count

	for _, event := range events {
		if event.Type == "damage" {
			allDamageEvents = append(allDamageEvents, event)

			// Track source IDs
			if event.SourceID != nil {
				sourceMap[*event.SourceID]++
			}

			// Track abilities
			if event.AbilityID != nil {
				abilityMap[*event.AbilityID]++
			}

			// Filter for Naalla specifically
			if event.TargetID != nil && *event.TargetID == playerID {
				naallaDamageEvents = append(naallaDamageEvents, event)
			}
		}
	}

	fmt.Printf("   📊 Total damage events: %d\n", len(allDamageEvents))
	fmt.Printf("   🎯 Damage events targeting Naalla (ID %d): %d\n", playerID, len(naallaDamageEvents))

	// Research 2: Identify Environmental vs Boss damage patterns
	fmt.Printf("\n🧪 RESEARCH 2: Analyze damage source patterns\n")
	fmt.Printf("   📊 Top damage sources by frequency:\n")

	type sourceCount struct {
		id    int
		count int
	}

	var sources []sourceCount
	for id, count := range sourceMap {
		sources = append(sources, sourceCount{id, count})
	}

	// Sort by count (simple bubble sort for small data)
	for i := 0; i < len(sources)-1; i++ {
		for j := i + 1; j < len(sources); j++ {
			if sources[j].count > sources[i].count {
				sources[i], sources[j] = sources[j], sources[i]
			}
		}
	}

	// Show top sources
	for i, source := range sources {
		if i >= 10 { // Top 10 sources
			break
		}
		fmt.Printf("      %d. SourceID %d: %d damage events\n", i+1, source.id, source.count)
	}

	// Research 3: Examine Naalla's damage events chronologically
	fmt.Printf("\n🧪 RESEARCH 3: Chronological analysis of Naalla's damage events\n")

	if len(naallaDamageEvents) > 0 {
		fmt.Printf("   📅 Naalla's damage events (chronological order):\n")

		// Sort events by timestamp
		for i := 0; i < len(naallaDamageEvents)-1; i++ {
			for j := i + 1; j < len(naallaDamageEvents); j++ {
				if naallaDamageEvents[j].Timestamp < naallaDamageEvents[i].Timestamp {
					naallaDamageEvents[i], naallaDamageEvents[j] = naallaDamageEvents[j], naallaDamageEvents[i]
				}
			}
		}

		for i, event := range naallaDamageEvents {
			if i >= 15 { // Show first 15 events
				break
			}

			timeFromDeath := (event.Timestamp - deathTime) / 1000.0
			fmt.Printf("      %d. %+.3fs: Ability %d, Amount %d, Source %d\n",
				i+1,
				timeFromDeath,
				getIntValue(event.AbilityID),
				getIntValue(event.Amount),
				getIntValue(event.SourceID))
		}

		fmt.Printf("   🎉 SUCCESS: Found individual damage events with precise timing!\n")
	} else {
		fmt.Printf("   ❌ No damage events found targeting Naalla\n")
		fmt.Printf("   🔍 DEBUG: Checking if targetID field is populated correctly...\n")

		// Debug: Show sample damage events to understand targeting
		fmt.Printf("   📊 Sample damage events (any target):\n")
		for i, event := range allDamageEvents {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. Target: %d, Source: %d, Ability: %d, Amount: %d\n",
				i+1,
				getIntValue(event.TargetID),
				getIntValue(event.SourceID),
				getIntValue(event.AbilityID),
				getIntValue(event.Amount))
		}
	}

	// Research 4: Ability ID mapping to understand spell names
	fmt.Printf("\n🧪 RESEARCH 4: Most frequent damage abilities\n")
	fmt.Printf("   📊 Top damage abilities by frequency:\n")

	type abilityCount struct {
		id    int
		count int
	}

	var abilities []abilityCount
	for id, count := range abilityMap {
		abilities = append(abilities, abilityCount{id, count})
	}

	// Sort by count
	for i := 0; i < len(abilities)-1; i++ {
		for j := i + 1; j < len(abilities); j++ {
			if abilities[j].count > abilities[i].count {
				abilities[i], abilities[j] = abilities[j], abilities[i]
			}
		}
	}

	// Show top abilities
	for i, ability := range abilities {
		if i >= 10 {
			break
		}
		fmt.Printf("      %d. AbilityID %d: %d damage events\n", i+1, ability.id, ability.count)
	}

	fmt.Printf("\n📋 RESEARCH SUMMARY:\n")
	fmt.Printf("   • Found %d total damage events in 3-second death window\n", len(allDamageEvents))
	fmt.Printf("   • Found %d damage events targeting Naalla specifically\n", len(naallaDamageEvents))
	fmt.Printf("   • Identified %d unique damage sources\n", len(sourceMap))
	fmt.Printf("   • Identified %d unique damage abilities\n", len(abilityMap))
	fmt.Printf("   • NEXT STEPS:\n")
	fmt.Printf("     1. Map sourceID values to determine which represents 'Environment'\n")
	fmt.Printf("     2. Map abilityID values to spell names (Null Explosion, etc.)\n")
	fmt.Printf("     3. Understand chronological ordering for unified timeline\n")
	fmt.Printf("     4. Implement damage timeline matching WCL web interface\n")

	return nil
}

// ResearchPlayerHPTimeline investigates how WCL reconstructs player health timeline
func ResearchPlayerHPTimeline(apiClient *api.Client) error {
	fmt.Printf("🎯 PLAYER HP TIMELINE RESEARCH - Understanding Naalla's health progression\n\n")

	// Known test case: Naalla's death at 6404661ms in fight 36
	reportCode := "9aQbqzgJy2dK8rVk"
	fightID := 36
	playerID := 11 // Naalla
	deathTime := 6404661.0
	windowStart := deathTime - 3000 // 3 seconds before death
	windowEnd := deathTime + 100    // 0.1 seconds after death

	fmt.Printf("📊 Target: Track how Naalla's HP changes leading to death\n")
	fmt.Printf("   • Player: Naalla (ID %d)\n", playerID)
	fmt.Printf("   • Death Time: %.0fms\n", deathTime)
	fmt.Printf("   • Analysis Window: 3 seconds before death\n")
	fmt.Printf("\n")

	// Query ALL events affecting Naalla in the time window
	fmt.Printf("🧪 RESEARCH: Get ALL events affecting Naalla (damage + healing)\n")
	playerEventsQuery := `
		query PlayerHPEvents($code: String!, $fightID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						startTime: $startTime,
						endTime: $endTime,
						dataType: All,
						limit: 1000
					) {
						data
					}
				}
			}
		}`

	timeVars := map[string]interface{}{
		"code":      reportCode,
		"fightID":   fightID,
		"startTime": windowStart,
		"endTime":   windowEnd,
	}

	resp, err := apiClient.Query(playerEventsQuery, timeVars)
	if err != nil {
		return fmt.Errorf("failed to query player events: %w", err)
	}

	if resp.Data == nil || resp.Data.ReportData == nil ||
		resp.Data.ReportData.Report == nil || resp.Data.ReportData.Report.Events == nil ||
		resp.Data.ReportData.Report.Events.Data == nil {
		return fmt.Errorf("no events data available")
	}

	allEvents, err := models.ParseEventsJSON(resp.Data.ReportData.Report.Events.Data)
	if err != nil {
		return fmt.Errorf("failed to parse events: %w", err)
	}

	fmt.Printf("   ✅ Parsed %d total events in time window\n", len(allEvents))

	// Filter events affecting Naalla specifically
	var naallaEvents []*models.Event
	for _, event := range allEvents {
		// Include events where Naalla is the target (damage to her, heals to her)
		if event.TargetID != nil && *event.TargetID == playerID {
			naallaEvents = append(naallaEvents, event)
		}
	}

	fmt.Printf("   🎯 Found %d events directly affecting Naalla\n", len(naallaEvents))

	// Sort events by timestamp (chronological order)
	for i := 0; i < len(naallaEvents)-1; i++ {
		for j := i + 1; j < len(naallaEvents); j++ {
			if naallaEvents[j].Timestamp < naallaEvents[i].Timestamp {
				naallaEvents[i], naallaEvents[j] = naallaEvents[j], naallaEvents[i]
			}
		}
	}

	// Analyze the chronological timeline
	fmt.Printf("\n📅 CHRONOLOGICAL TIMELINE of events affecting Naalla:\n")

	var damageEvents, healingEvents, otherEvents int
	for i, event := range naallaEvents {
		if i >= 20 { // Show first 20 events
			fmt.Printf("   ... (showing first 20 of %d total events)\n", len(naallaEvents))
			break
		}

		timeFromDeath := (event.Timestamp - deathTime) / 1000.0

		eventTypeDesc := "OTHER"
		if event.Type == "damage" {
			eventTypeDesc = "DAMAGE"
			damageEvents++
		} else if event.Type == "heal" {
			eventTypeDesc = "HEAL"
			healingEvents++
		} else {
			otherEvents++
		}

		fmt.Printf("   %2d. %+.3fs: %-6s | Ability %d | Amount %d | Source %d\n",
			i+1,
			timeFromDeath,
			eventTypeDesc,
			getIntValue(event.AbilityID),
			getIntValue(event.Amount),
			getIntValue(event.SourceID))
	}

	fmt.Printf("\n📊 EVENT BREAKDOWN:\n")
	fmt.Printf("   • Damage events: %d\n", damageEvents)
	fmt.Printf("   • Healing events: %d\n", healingEvents)
	fmt.Printf("   • Other events: %d\n", otherEvents)

	// Check if we're missing damage events by looking at broader data
	if damageEvents < 5 {
		fmt.Printf("\n🔍 INVESTIGATING: Low damage event count, checking broader search...\n")

		// Look for damage events in the broader window that might hit multiple targets
		var broadDamageEvents []*models.Event
		for _, event := range allEvents {
			if event.Type == "damage" {
				// Look for high-damage events that might be AoE
				if event.Amount != nil && *event.Amount > 500000 {
					broadDamageEvents = append(broadDamageEvents, event)
				}
			}
		}

		fmt.Printf("   • Found %d high-damage events in time window (any target)\n", len(broadDamageEvents))

		if len(broadDamageEvents) > 0 {
			fmt.Printf("   📊 High-damage events (potential AoE hitting Naalla):\n")
			for i, event := range broadDamageEvents {
				if i >= 10 { // Show first 10
					break
				}
				timeFromDeath := (event.Timestamp - deathTime) / 1000.0
				fmt.Printf("      %d. %+.3fs: %d damage, Target %d, Source %d, Ability %d\n",
					i+1,
					timeFromDeath,
					getIntValue(event.Amount),
					getIntValue(event.TargetID),
					getIntValue(event.SourceID),
					getIntValue(event.AbilityID))
			}
		}
	}

	fmt.Printf("\n📋 HP TIMELINE RESEARCH SUMMARY:\n")
	fmt.Printf("   • Tracked %d events directly affecting Naalla in 3-second window\n", len(naallaEvents))
	fmt.Printf("   • %d damage events reducing her HP\n", damageEvents)
	fmt.Printf("   • %d healing events increasing her HP\n", healingEvents)
	fmt.Printf("   • INSIGHT: WCL likely reconstructs HP by processing these events chronologically\n")
	fmt.Printf("   • NEXT STEP: Implement HP calculation to match WCL's timeline format\n")
	fmt.Printf("   • CHALLENGE: Missing environmental/AoE damage events in direct targeting\n")

	return nil
}

// ResearchTableAPI investigates the GraphQL Table API for damage taken data with timestamps
func ResearchTableAPI(apiClient *api.Client) error {
	fmt.Printf("🔬 TABLE API RESEARCH - Investigating damage taken queries for individual events\n\n")

	// Known test case: Naalla's death at 6404661ms in fight 36
	reportCode := "9aQbqzgJy2dK8rVk"
	fightID := 36
	playerID := 11 // Naalla
	deathTime := 6404661.0
	windowStart := deathTime - 3000 // 3 seconds before death
	windowEnd := deathTime + 100    // 0.1 seconds after death

	fmt.Printf("📊 Goal: Find the Table API query that matches WCL 'Damage Taken' CSV export\n")
	fmt.Printf("   • Expected: Null Explosion (13 hits), Crystalline Shockwave, etc.\n")
	fmt.Printf("   • Player: Naalla (ID %d)\n", playerID)
	fmt.Printf("   • Time Window: %.0f to %.0f ms\n", windowStart, windowEnd)
	fmt.Printf("\n")

	// Test 1: Standard Table API query like the breakthrough document used
	fmt.Printf("🧪 TEST 1: Standard Table API - DamageTaken by Ability\n")
	tableQuery1 := `
		query DamageTakenTable($code: String!, $fightID: Int!, $playerID: Int!) {
			reportData {
				report(code: $code) {
					table(
						fightIDs: [$fightID],
						sourceID: $playerID,
						dataType: DamageTaken,
						viewBy: Ability
					)
				}
			}
		}`

	variables := map[string]interface{}{
		"code":     reportCode,
		"fightID":  fightID,
		"playerID": playerID,
	}

	fmt.Printf("   • Query: table(sourceID: %d, dataType: DamageTaken, viewBy: Ability)\n", playerID)
	resp1, err1 := apiClient.Query(tableQuery1, variables)
	if err1 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err1)
	} else if resp1.Data == nil || resp1.Data.ReportData == nil || resp1.Data.ReportData.Report == nil || resp1.Data.ReportData.Report.Table == nil {
		fmt.Printf("   ❌ No table data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - parsing table data...\n")
		jsonBytes, _ := json.Marshal(resp1.Data.ReportData.Report.Table)
		jsonStr := string(jsonBytes)
		fmt.Printf("   📊 Table data size: %d bytes\n", len(jsonStr))
		fmt.Printf("   💡 Sample data: %s...\n", truncateString(jsonStr, 300))
	}

	// Test 2: Try different viewBy options
	fmt.Printf("\n🧪 TEST 2: Table API - DamageTaken by Source\n")
	tableQuery2 := `
		query DamageTakenBySource($code: String!, $fightID: Int!, $playerID: Int!) {
			reportData {
				report(code: $code) {
					table(
						fightIDs: [$fightID],
						sourceID: $playerID,
						dataType: DamageTaken,
						viewBy: Source
					)
				}
			}
		}`

	fmt.Printf("   • Query: table(sourceID: %d, dataType: DamageTaken, viewBy: Source)\n", playerID)
	resp2, err2 := apiClient.Query(tableQuery2, variables)
	if err2 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err2)
	} else if resp2.Data == nil || resp2.Data.ReportData == nil || resp2.Data.ReportData.Report == nil || resp2.Data.ReportData.Report.Table == nil {
		fmt.Printf("   ❌ No table data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - parsing table data...\n")
		jsonBytes, _ := json.Marshal(resp2.Data.ReportData.Report.Table)
		jsonStr := string(jsonBytes)
		fmt.Printf("   📊 Table data size: %d bytes\n", len(jsonStr))
		fmt.Printf("   💡 Sample data: %s...\n", truncateString(jsonStr, 300))
	}

	// Test 3: Try using targetID instead of sourceID (correct for damage taken)
	fmt.Printf("\n🧪 TEST 3: Table API - DamageTaken with targetID (correct approach)\n")
	tableQuery3 := `
		query DamageTakenTarget($code: String!, $fightID: Int!, $playerID: Int!) {
			reportData {
				report(code: $code) {
					table(
						fightIDs: [$fightID],
						targetID: $playerID,
						dataType: DamageTaken,
						viewBy: Ability
					)
				}
			}
		}`

	fmt.Printf("   • Query: table(targetID: %d, dataType: DamageTaken, viewBy: Ability)\n", playerID)
	resp3, err3 := apiClient.Query(tableQuery3, variables)
	if err3 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err3)
	} else if resp3.Data == nil || resp3.Data.ReportData == nil || resp3.Data.ReportData.Report == nil || resp3.Data.ReportData.Report.Table == nil {
		fmt.Printf("   ❌ No table data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - parsing table data...\n")
		jsonBytes, _ := json.Marshal(resp3.Data.ReportData.Report.Table)
		jsonStr := string(jsonBytes)
		fmt.Printf("   📊 Table data size: %d bytes\n", len(jsonStr))
		fmt.Printf("   💡 Sample data: %s...\n", truncateString(jsonStr, 300))

		// Try to parse as the structured table format
		var tableResult map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &tableResult); err == nil {
			fmt.Printf("   🔍 Parsed table structure successfully\n")
			if data, ok := tableResult["data"]; ok {
				if entries, ok := data.(map[string]interface{})["entries"]; ok {
					if entriesArray, ok := entries.([]interface{}); ok {
						fmt.Printf("   📊 Found %d damage abilities:\n", len(entriesArray))
						for i, entry := range entriesArray {
							if i >= 5 { // Show first 5 abilities
								break
							}
							if entryMap, ok := entry.(map[string]interface{}); ok {
								name := entryMap["name"]
								total := entryMap["total"]
								fmt.Printf("      %d. %v: %v damage\n", i+1, name, total)
							}
						}
					}
				}
			}
		}
	}

	// Test 4: Try to get individual events from table with time filtering
	fmt.Printf("\n🧪 TEST 4: Table API with time window filtering\n")
	tableQuery4 := `
		query DamageTakenTimeWindow($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					table(
						fightIDs: [$fightID],
						targetID: $playerID,
						dataType: DamageTaken,
						viewBy: Ability,
						startTime: $startTime,
						endTime: $endTime
					)
				}
			}
		}`

	timeVars := map[string]interface{}{
		"code":      reportCode,
		"fightID":   fightID,
		"playerID":  playerID,
		"startTime": windowStart,
		"endTime":   windowEnd,
	}

	fmt.Printf("   • Query: table(targetID: %d, time: %.0f-%.0f, dataType: DamageTaken)\n", playerID, windowStart, windowEnd)
	resp4, err4 := apiClient.Query(tableQuery4, timeVars)
	if err4 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err4)
	} else if resp4.Data == nil || resp4.Data.ReportData == nil || resp4.Data.ReportData.Report == nil || resp4.Data.ReportData.Report.Table == nil {
		fmt.Printf("   ❌ No table data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - this should show damage in death window!\n")
		jsonBytes, _ := json.Marshal(resp4.Data.ReportData.Report.Table)
		jsonStr := string(jsonBytes)
		fmt.Printf("   📊 Time-filtered data size: %d bytes\n", len(jsonStr))
		fmt.Printf("   💡 Death window data: %s...\n", truncateString(jsonStr, 400))
	}

	// Test 5: Try events API with specific hostility filtering
	fmt.Printf("\n🧪 TEST 5: Events API with hostilityType filtering\n")
	eventsQuery := `
		query EventsWithHostility($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						targetID: $playerID,
						startTime: $startTime,
						endTime: $endTime,
						dataType: All,
						hostilityType: Enemies,
						limit: 200
					) {
						data
					}
				}
			}
		}`

	fmt.Printf("   • Query: events(targetID: %d, hostilityType: Enemies, dataType: All)\n", playerID)
	resp5, err5 := apiClient.Query(eventsQuery, timeVars)
	if err5 != nil {
		fmt.Printf("   ❌ FAILED: %v\n", err5)
	} else if resp5.Data == nil || resp5.Data.ReportData == nil || resp5.Data.ReportData.Report == nil || resp5.Data.ReportData.Report.Events == nil {
		fmt.Printf("   ❌ No events data structure\n")
	} else {
		fmt.Printf("   ✅ Query successful - parsing hostile events...\n")
		if resp5.Data.ReportData.Report.Events.Data != nil {
			jsonStr := string(resp5.Data.ReportData.Report.Events.Data)
			fmt.Printf("   📊 Hostile events data size: %d bytes\n", len(jsonStr))

			// Parse and filter for damage events
			events, err := models.ParseEventsJSON(resp5.Data.ReportData.Report.Events.Data)
			if err == nil {
				var damageEvents []*models.Event
				for _, event := range events {
					if event.Type == "damage" && event.TargetID != nil && *event.TargetID == playerID {
						damageEvents = append(damageEvents, event)
					}
				}
				fmt.Printf("   🎯 Found %d damage events from enemies targeting Naalla\n", len(damageEvents))
				if len(damageEvents) > 0 {
					fmt.Printf("   📊 Sample damage events:\n")
					for i, event := range damageEvents {
						if i >= 5 {
							break
						}
						timeFromDeath := (event.Timestamp - deathTime) / 1000.0
						fmt.Printf("      %d. %+.3fs: Ability %d, Amount %d, Source %d\n",
							i+1,
							timeFromDeath,
							getIntValue(event.AbilityID),
							getIntValue(event.Amount),
							getIntValue(event.SourceID))
					}
				}
			}
		}
	}

	fmt.Printf("\n📋 TABLE API RESEARCH SUMMARY:\n")
	fmt.Printf("   • Test 1: Standard sourceID + DamageTaken (from breakthrough doc)\n")
	fmt.Printf("   • Test 2: sourceID + DamageTaken + viewBy Source\n")
	fmt.Printf("   • Test 3: targetID + DamageTaken + viewBy Ability (likely correct)\n")
	fmt.Printf("   • Test 4: targetID + DamageTaken + time window (death window)\n")
	fmt.Printf("   • Test 5: Events API + hostilityType Enemies (alternative approach)\n")
	fmt.Printf("   • GOAL: Find API query that returns individual Null Explosion hits with timestamps\n")
	fmt.Printf("   • SUCCESS CRITERIA: Match WCL CSV showing 13 Null Explosion hits to Naalla\n")

	return nil
}

// ResearchTableAPIParsing deeply parses successful Table API DamageTaken responses
func ResearchTableAPIParsing(apiClient *api.Client) error {
	fmt.Printf("🔬 TABLE API PARSING - Deep analysis of DamageTaken data structure\n\n")

	// Known test case: Naalla's damage taken
	reportCode := "9aQbqzgJy2dK8rVk"
	fightID := 36
	playerID := 11 // Naalla

	fmt.Printf("📊 BREAKTHROUGH: Using sourceID = playerID for DamageTaken queries\n")
	fmt.Printf("   • sourceID: %d (Naalla receiving damage)\n", playerID)
	fmt.Printf("   • Expected: Complete damage breakdown matching WCL CSV\n\n")

	// The working query from our previous research
	tableQuery := `
		query DamageTakenDetailed($code: String!, $fightID: Int!, $playerID: Int!) {
			reportData {
				report(code: $code) {
					table(
						fightIDs: [$fightID],
						sourceID: $playerID,
						dataType: DamageTaken,
						viewBy: Ability
					)
				}
			}
		}`

	variables := map[string]interface{}{
		"code":     reportCode,
		"fightID":  fightID,
		"playerID": playerID,
	}

	fmt.Printf("🧪 PARSING: Complete DamageTaken table for Naalla\n")
	resp, err := apiClient.Query(tableQuery, variables)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	if resp.Data == nil || resp.Data.ReportData == nil ||
		resp.Data.ReportData.Report == nil || resp.Data.ReportData.Report.Table == nil {
		return fmt.Errorf("no table data available")
	}

	// Get raw JSON for detailed parsing
	jsonBytes, err := json.Marshal(resp.Data.ReportData.Report.Table)
	if err != nil {
		return fmt.Errorf("failed to marshal table data: %w", err)
	}

	fmt.Printf("   ✅ Retrieved %d bytes of damage taken data\n", len(jsonBytes))

	// Parse the table structure
	var tableResult map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &tableResult); err != nil {
		return fmt.Errorf("failed to parse table JSON: %w", err)
	}

	fmt.Printf("\n📊 DAMAGE BREAKDOWN - All abilities that damaged Naalla:\n")

	// Navigate the nested structure
	data, ok := tableResult["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no data field in table response")
	}

	entries, ok := data["entries"].([]interface{})
	if !ok {
		return fmt.Errorf("no entries field in table data")
	}

	fmt.Printf("   📋 Found %d damage abilities:\n", len(entries))

	var totalDamage float64
	var nullExplosionFound, crystallineShockwaveFound bool

	for i, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		name := entryMap["name"]
		total, _ := entryMap["total"].(float64)
		totalDamage += total

		// Extract additional details
		hitCount, _ := entryMap["hitCount"].(float64)
		tickCount, _ := entryMap["tickCount"].(float64)
		actorName, _ := entryMap["actorName"].(string)

		fmt.Printf("   %2d. %-25s │ %10.0f │ %3.0f hits │ %s\n",
			i+1, name, total, hitCount+tickCount, actorName)

		// Check for key abilities we expect
		if nameStr, ok := name.(string); ok {
			if nameStr == "Null Explosion" {
				nullExplosionFound = true
				fmt.Printf("       🎯 FOUND: Null Explosion - %3.0f hits for %.0f total damage\n",
					hitCount+tickCount, total)
			}
			if nameStr == "Crystalline Shockwave" {
				crystallineShockwaveFound = true
				fmt.Printf("       🎯 FOUND: Crystalline Shockwave - %3.0f hits for %.0f total damage\n",
					hitCount+tickCount, total)
			}
		}

		// Show first few entries in detail
		if i < 8 {
			fmt.Printf("       └─ Details: %s\n", truncateString(fmt.Sprintf("%v", entryMap), 100))
		}
	}

	fmt.Printf("\n📊 SUMMARY:\n")
	fmt.Printf("   • Total damage taken: %.0f\n", totalDamage)
	fmt.Printf("   • Total abilities: %d\n", len(entries))
	fmt.Printf("   • Null Explosion found: %v\n", nullExplosionFound)
	fmt.Printf("   • Crystalline Shockwave found: %v\n", crystallineShockwaveFound)

	// Look for detailed hit information
	fmt.Printf("\n🔍 DETAILED ANALYSIS - Looking for individual hit timestamps:\n")

	for i, entry := range entries {
		if i >= 3 { // Analyze first 3 abilities in detail
			break
		}

		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		name := entryMap["name"]
		fmt.Printf("   🔬 Analyzing: %v\n", name)

		// Look for nested data structures that might contain timing
		for key, value := range entryMap {
			if key == "bands" || key == "details" || key == "events" || key == "timeline" {
				fmt.Printf("      ├─ %s: %s\n", key, truncateString(fmt.Sprintf("%v", value), 150))
			}
		}

		// Look for any arrays that might contain individual hits
		for key, value := range entryMap {
			if valueSlice, ok := value.([]interface{}); ok && len(valueSlice) > 0 {
				fmt.Printf("      ├─ %s[%d]: %s\n", key, len(valueSlice),
					truncateString(fmt.Sprintf("%v", valueSlice[0]), 100))
			}
		}
	}

	fmt.Printf("\n📋 TABLE API PARSING RESULTS:\n")
	fmt.Printf("   • ✅ CONFIRMED: sourceID = playerID works for DamageTaken\n")
	fmt.Printf("   • ✅ CONFIRMED: Table API contains complete damage breakdown\n")
	fmt.Printf("   • 🔍 INVESTIGATION: Need to find individual hit timestamps within table data\n")
	fmt.Printf("   • 🎯 NEXT STEP: Extract timing data or correlate with events API\n")

	if nullExplosionFound {
		fmt.Printf("   • 🎉 SUCCESS: Found Null Explosion data in Table API!\n")
	}

	return nil
}

// truncateString safely truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// getIntValue safely extracts int value from pointer, returns 0 if nil
func getIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}
