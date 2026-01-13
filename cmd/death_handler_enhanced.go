package cmd

import (
	"encoding/json"
	"fmt"
	"sort"

	"wclogs-cli/api"
	"wclogs-cli/models"
	"wclogs-cli/services"

	"github.com/fatih/color"
)

// DamageSource represents damage taken from the table API
type DamageSource struct {
	AbilityName string  `json:"name"`
	Amount      int64   `json:"amount"`
	Percentage  float64 `json:"percentage"`
}

// EnhancedTimelineEvent represents an event for WCL-style display with HP tracking
type EnhancedTimelineEvent struct {
	TimeFromDeath  float64
	Type           string
	AbilityName    string
	SourceName     string
	Amount         int
	HitPoints      int
	MaxHitPoints   int
	IsOverkill     bool
	OverkillAmount int
}

// executeEnhancedDeathAnalysis provides comprehensive death analysis using table + events APIs
func executeEnhancedDeathAnalysis(apiClient *api.Client, reportCode string, fightID, playerID int, startTime, deathTime float64, lookupService *services.LookupService, verbose bool) {
	if verbose {
		fmt.Printf("    🔬 Enhanced: Starting comprehensive death analysis for player %d\n", playerID)
	}

	// Step 1: Get damage taken data using table API (like WCL web interface)
	damageSources, err := getDamageTakenData(apiClient, reportCode, fightID, playerID, verbose)
	if err != nil {
		fmt.Printf("    ❌ Failed to get damage taken data: %v\n", err)
		return
	}

	// Step 2: Get healing events using events API (this works well)
	windowStart := deathTime - 3000 // 3 seconds before death
	windowEnd := deathTime + 100    // 0.1 seconds after death
	healingEvents := queryPlayerEvents(apiClient, reportCode, fightID, playerID, windowStart, windowEnd, "Healing", verbose)

	// Step 3: Display comprehensive death analysis
	displayComprehensiveDeathAnalysis(damageSources, healingEvents, deathTime, lookupService, verbose)
}

// getDamageTakenData uses the table API to get aggregated damage data like WCL web interface
func getDamageTakenData(apiClient *api.Client, reportCode string, fightID, playerID int, verbose bool) ([]DamageSource, error) {
	if verbose {
		fmt.Printf("    🔬 Enhanced: Querying table API for damage taken data\n")
	}

	// Query for damage taken by the specific player using the approach that works
	query := `
		query GetDamageTakenTable($code: String!, $fightID: Int!, $playerID: Int!) {
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

	response, err := apiClient.Query(query, variables)
	if err != nil {
		return nil, fmt.Errorf("damage taken table query failed: %w", err)
	}

	if response.Data == nil || response.Data.ReportData == nil ||
		response.Data.ReportData.Report == nil || response.Data.ReportData.Report.Table == nil {
		return nil, fmt.Errorf("no damage taken data found")
	}

	// Parse the JSON response from table API
	jsonBytes, err := json.Marshal(response.Data.ReportData.Report.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal table data: %w", err)
	}

	if verbose {
		fmt.Printf("    🔬 Enhanced: Damage taken table response: %s\n", string(jsonBytes))
	}

	var tableResult map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &tableResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal table data: %w", err)
	}

	// Extract damage sources from table data
	damageSources, err := parseDamageTakenFromTable(tableResult, verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to parse damage table: %w", err)
	}

	if verbose {
		fmt.Printf("    🔬 Enhanced: Found %d damage sources targeting player\n", len(damageSources))
	}

	return damageSources, nil
}

// parseDamageTakenFromTable extracts damage taken information from WCL table API response
func parseDamageTakenFromTable(tableData map[string]interface{}, verbose bool) ([]DamageSource, error) {
	var damageSources []DamageSource

	// Navigate to the entries array
	data, ok := tableData["data"]
	if !ok {
		return damageSources, fmt.Errorf("no data field in table response")
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return damageSources, fmt.Errorf("data field is not a map")
	}

	entries, ok := dataMap["entries"]
	if !ok {
		return damageSources, fmt.Errorf("no entries field in data")
	}

	entriesArray, ok := entries.([]interface{})
	if !ok {
		return damageSources, fmt.Errorf("entries field is not an array")
	}

	// Parse each entry to extract damage sources
	for _, item := range entriesArray {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract ability name
		abilityName, ok := entry["name"].(string)
		if !ok || abilityName == "" {
			continue
		}

		// Extract damage amount
		var damageAmount int64
		if total, ok := entry["total"].(float64); ok {
			damageAmount = int64(total)
		}

		// Skip if no damage
		if damageAmount <= 0 {
			continue
		}

		// Look for sources that damaged this player with this ability
		if sources, ok := entry["sources"].([]interface{}); ok {
			for _, sourceItem := range sources {
				if sourceMap, ok := sourceItem.(map[string]interface{}); ok {
					if sourceName, ok := sourceMap["name"].(string); ok {
						if sourceTotal, ok := sourceMap["total"].(float64); ok && sourceTotal > 0 {
							damageSources = append(damageSources, DamageSource{
								AbilityName: fmt.Sprintf("%s ← %s", abilityName, sourceName),
								Amount:      int64(sourceTotal),
								Percentage:  0, // Will calculate later
							})
						}
					}
				}
			}
		} else {
			// If no sources breakdown, use the total for the ability
			damageSources = append(damageSources, DamageSource{
				AbilityName: abilityName,
				Amount:      damageAmount,
				Percentage:  0, // Will calculate later
			})
		}
	}

	// Calculate total damage and percentages
	var totalDamage int64
	for _, source := range damageSources {
		totalDamage += source.Amount
	}

	for i := range damageSources {
		if totalDamage > 0 {
			damageSources[i].Percentage = float64(damageSources[i].Amount) / float64(totalDamage) * 100
		}
	}

	// Sort by damage amount (highest first)
	sort.Slice(damageSources, func(i, j int) bool {
		return damageSources[i].Amount > damageSources[j].Amount
	})

	if verbose && len(damageSources) > 0 {
		fmt.Printf("    🔬 Enhanced: Top damage source: %s (%s)\n",
			damageSources[0].AbilityName,
			models.FormatNumber(damageSources[0].Amount))
	}

	return damageSources, nil
}

// queryPlayerEvents queries for specific event types targeting the player (same as before but simplified)
func queryPlayerEvents(apiClient *api.Client, reportCode string, fightID, playerID int, startTime, endTime float64, dataType string, verbose bool) []*models.Event {
	query := fmt.Sprintf(`
		query PlayerEvents($code: String!, $fightID: Int!, $playerID: Int!, $startTime: Float!, $endTime: Float!) {
			reportData {
				report(code: $code) {
					events(
						fightIDs: [$fightID],
						targetID: $playerID,
						startTime: $startTime,
						endTime: $endTime,
						dataType: %s,
						limit: 100,
						includeResources: true
					) {
						data
					}
				}
			}
		}`, dataType)

	variables := map[string]interface{}{
		"code":      reportCode,
		"fightID":   fightID,
		"playerID":  playerID,
		"startTime": startTime,
		"endTime":   endTime,
	}

	response, err := apiClient.Query(query, variables)
	if err != nil {
		if verbose {
			fmt.Printf("    ❌ Failed to query %s events: %v\n", dataType, err)
		}
		return nil
	}

	if response.Data == nil || response.Data.ReportData == nil ||
		response.Data.ReportData.Report == nil ||
		response.Data.ReportData.Report.Events == nil {
		return nil
	}

	events, err := models.ParseEventsJSON(response.Data.ReportData.Report.Events.Data)
	if err != nil {
		if verbose {
			fmt.Printf("    ❌ Failed to parse %s events: %v\n", dataType, err)
		}
		return nil
	}

	return events
}

// displayComprehensiveDeathAnalysis shows both damage sources and healing timeline
func displayComprehensiveDeathAnalysis(damageSources []DamageSource, healingEvents []*models.Event, deathTime float64, lookupService *services.LookupService, verbose bool) {
	fmt.Printf("    💀 Enhanced Death Analysis (WCL Style):\n\n")

	// Section 1: Damage Taken Summary (like WCL damage taken table)
	if len(damageSources) > 0 {
		fmt.Printf("    ⚔️  DAMAGE SOURCES (from Table API):\n")
		fmt.Printf("    ┌─────────────────────────────────────────────┬──────────────────┬─────────┐\n")
		fmt.Printf("    │ Ability Name                                │ Damage Amount    │ Percent │\n")
		fmt.Printf("    ├─────────────────────────────────────────────┼──────────────────┼─────────┤\n")

		var totalDamage int64
		for _, source := range damageSources {
			totalDamage += source.Amount
		}

		displayCount := 0
		for _, source := range damageSources {
			if displayCount >= 10 { // Show top 10 damage sources
				break
			}

			abilityName := source.AbilityName
			if len(abilityName) > 43 {
				abilityName = abilityName[:40] + "..."
			}

			damageStr := color.HiRedString(models.FormatNumber(source.Amount))

			percentage := float64(source.Amount) / float64(totalDamage) * 100
			percentStr := fmt.Sprintf("%.1f%%", percentage)

			fmt.Printf("    │ %-43s │ %-16s │ %-7s │\n", abilityName, damageStr, percentStr)
			displayCount++
		}

		fmt.Printf("    └─────────────────────────────────────────────┴──────────────────┴─────────┘\n")
		fmt.Printf("    📊 Total Damage Taken: %s\n\n", color.HiRedString(models.FormatNumber(totalDamage)))
	} else {
		fmt.Printf("    ⚠️  No damage sources found from table API\n\n")
	}

	// Section 2: Healing Timeline (like death CSV timeline but focused on healing)
	if len(healingEvents) > 0 {
		// Convert healing events to timeline format
		var timelineEvents []*EnhancedTimelineEvent

		for _, event := range healingEvents {
			if event.Amount != nil && *event.Amount > 0 {
				timeFromDeath := (event.Timestamp - deathTime) / 1000.0

				timelineEvent := &EnhancedTimelineEvent{
					TimeFromDeath: timeFromDeath,
					Type:          "heal",
					Amount:        *event.Amount,
				}

				if event.AbilityID != nil {
					timelineEvent.AbilityName = lookupService.GetAbilityName(*event.AbilityID)
				}
				if event.SourceID != nil {
					timelineEvent.SourceName = lookupService.GetActorName(*event.SourceID)
				}

				// Only include significant heals and those close to death
				if *event.Amount >= 10000 {
					timelineEvents = append(timelineEvents, timelineEvent)
				}
			}
		}

		// Add death event at 0.00s
		deathEvent := &EnhancedTimelineEvent{
			TimeFromDeath: 0.0,
			Type:          "death",
			AbilityName:   "Death",
			SourceName:    "Killing Blow",
			Amount:        0,
		}
		timelineEvents = append(timelineEvents, deathEvent)

		// Sort by time (death first, then most recent)
		sort.Slice(timelineEvents, func(i, j int) bool {
			if timelineEvents[i].Type == "death" && timelineEvents[j].Type != "death" {
				return true
			}
			if timelineEvents[j].Type == "death" && timelineEvents[i].Type != "death" {
				return false
			}
			return timelineEvents[i].TimeFromDeath > timelineEvents[j].TimeFromDeath
		})

		fmt.Printf("    💚 HEALING ATTEMPTS (last 3 seconds):\n")
		fmt.Printf("    ┌──────────┬─────────────────────────────────────────────┬──────────────────┐\n")
		fmt.Printf("    │ Time     │ Ability                                     │ Amount           │\n")
		fmt.Printf("    ├──────────┼─────────────────────────────────────────────┼──────────────────┤\n")

		displayCount := 0
		var totalHealing int64

		for _, event := range timelineEvents {
			if displayCount >= 15 { // Show top 15 most recent events
				break
			}

			timeStr := formatTimeWCLStyle(event.TimeFromDeath, event.Type == "death")

			if event.Type == "death" {
				fmt.Printf("    │ %-8s │ %-43s │ %-16s │\n",
					timeStr,
					color.HiRedString("💀 Death Event"),
					color.HiRedString("-"))
			} else {
				abilityStr := event.AbilityName
				if event.SourceName != "" && event.SourceName != "Unknown" {
					abilityStr = fmt.Sprintf("%s ← %s", event.AbilityName, event.SourceName)
				}
				if len(abilityStr) > 43 {
					abilityStr = abilityStr[:40] + "..."
				}

				amountStr := color.HiGreenString("+%s", models.FormatNumber(int64(event.Amount)))
				totalHealing += int64(event.Amount)

				fmt.Printf("    │ %-8s │ %-43s │ %-16s │\n", timeStr, abilityStr, amountStr)
			}
			displayCount++
		}

		fmt.Printf("    └──────────┴─────────────────────────────────────────────┴──────────────────┘\n")
		if totalHealing > 0 {
			fmt.Printf("    📊 Total Healing Received: %s\n\n", color.HiGreenString(models.FormatNumber(totalHealing)))
		}
	} else {
		fmt.Printf("    ⚠️  No healing events found\n\n")
	}

	// Section 3: Analysis Summary
	fmt.Printf("    🔍 DEATH ANALYSIS SUMMARY:\n")
	if len(damageSources) >= 2 {
		fmt.Printf("    • Top damage source: %s (%s damage)\n",
			color.HiYellowString(damageSources[0].AbilityName),
			color.HiRedString(models.FormatNumber(damageSources[0].Amount)))
		fmt.Printf("    • Second damage source: %s (%s damage)\n",
			color.HiYellowString(damageSources[1].AbilityName),
			color.HiRedString(models.FormatNumber(damageSources[1].Amount)))
	} else if len(damageSources) == 1 {
		fmt.Printf("    • Primary damage source: %s (%s damage)\n",
			color.HiYellowString(damageSources[0].AbilityName),
			color.HiRedString(models.FormatNumber(damageSources[0].Amount)))
	}

	if len(healingEvents) > 0 {
		fmt.Printf("    • Healers attempted to save the player with %d healing spells\n", len(healingEvents))
	}

	fmt.Printf("    • This enhanced analysis combines WCL table data with event timelines\n")

	// Section 2: Healing Timeline (like death CSV timeline but focused on healing)
	if len(healingEvents) > 0 {
		// Convert healing events to timeline format
		var timelineEvents []*EnhancedTimelineEvent

		for _, event := range healingEvents {
			if event.Amount != nil && *event.Amount > 0 {
				timeFromDeath := (event.Timestamp - deathTime) / 1000.0

				timelineEvent := &EnhancedTimelineEvent{
					TimeFromDeath: timeFromDeath,
					Type:          "heal",
					Amount:        *event.Amount,
				}

				if event.AbilityID != nil {
					timelineEvent.AbilityName = lookupService.GetAbilityName(*event.AbilityID)
				}
				if event.SourceID != nil {
					timelineEvent.SourceName = lookupService.GetActorName(*event.SourceID)
				}

				// Only include significant heals and those close to death
				if *event.Amount >= 10000 {
					timelineEvents = append(timelineEvents, timelineEvent)
				}
			}
		}

		// Add death event at 0.00s
		deathEvent := &EnhancedTimelineEvent{
			TimeFromDeath: 0.0,
			Type:          "death",
			AbilityName:   "Death",
			SourceName:    "Killing Blow",
			Amount:        0,
		}
		timelineEvents = append(timelineEvents, deathEvent)

		// Sort by time (death first, then most recent)
		sort.Slice(timelineEvents, func(i, j int) bool {
			if timelineEvents[i].Type == "death" && timelineEvents[j].Type != "death" {
				return true
			}
			if timelineEvents[j].Type == "death" && timelineEvents[i].Type != "death" {
				return false
			}
			return timelineEvents[i].TimeFromDeath > timelineEvents[j].TimeFromDeath
		})

		fmt.Printf("    💚 HEALING ATTEMPTS (last 3 seconds):\n")
		fmt.Printf("    ┌──────────┬─────────────────────────────────────────────┬──────────────────┐\n")
		fmt.Printf("    │ Time     │ Ability                                     │ Amount           │\n")
		fmt.Printf("    ├──────────┼─────────────────────────────────────────────┼──────────────────┤\n")

		displayCount := 0
		var totalHealing int64

		for _, event := range timelineEvents {
			if displayCount >= 15 { // Show top 15 most recent events
				break
			}

			timeStr := formatTimeWCLStyle(event.TimeFromDeath, event.Type == "death")

			if event.Type == "death" {
				fmt.Printf("    │ %-8s │ %-43s │ %-16s │\n",
					timeStr,
					color.HiRedString("💀 Death Event"),
					color.HiRedString("-"))
			} else {
				abilityStr := event.AbilityName
				if event.SourceName != "" && event.SourceName != "Unknown" {
					abilityStr = fmt.Sprintf("%s ← %s", event.AbilityName, event.SourceName)
				}
				if len(abilityStr) > 43 {
					abilityStr = abilityStr[:40] + "..."
				}

				amountStr := color.HiGreenString("+%s", models.FormatNumber(int64(event.Amount)))
				totalHealing += int64(event.Amount)

				fmt.Printf("    │ %-8s │ %-43s │ %-16s │\n", timeStr, abilityStr, amountStr)
			}
			displayCount++
		}

		fmt.Printf("    └──────────┴─────────────────────────────────────────────┴──────────────────┘\n")
		if totalHealing > 0 {
			fmt.Printf("    📊 Total Healing Received: %s\n\n", color.HiGreenString(models.FormatNumber(totalHealing)))
		}
	} else {
		fmt.Printf("    ⚠️  No healing events found\n\n")
	}

	// Section 3: Analysis Summary
	fmt.Printf("    🔍 DEATH ANALYSIS SUMMARY:\n")
	if len(damageSources) >= 2 {
		fmt.Printf("    • Top damage source: %s (%s damage)\n",
			color.HiYellowString(damageSources[0].AbilityName),
			color.HiRedString(models.FormatNumber(damageSources[0].Amount)))
		fmt.Printf("    • Second damage source: %s (%s damage)\n",
			color.HiYellowString(damageSources[1].AbilityName),
			color.HiRedString(models.FormatNumber(damageSources[1].Amount)))
	} else if len(damageSources) == 1 {
		fmt.Printf("    • Primary damage source: %s (%s damage)\n",
			color.HiYellowString(damageSources[0].AbilityName),
			color.HiRedString(models.FormatNumber(damageSources[0].Amount)))
	}

	if len(healingEvents) > 0 {
		fmt.Printf("    • Healers attempted to save the player with %d healing spells\n", len(healingEvents))
	}

	fmt.Printf("    • This enhanced analysis combines WCL table data with event timelines\n")
}

// formatTimeWCLStyle formats time exactly like WCL CSV
func formatTimeWCLStyle(timeFromDeath float64, isDeath bool) string {
	if isDeath {
		return "0.00s"
	} else if timeFromDeath >= 0 {
		return fmt.Sprintf("+%.2fs", timeFromDeath)
	} else {
		return fmt.Sprintf("%.2fs", timeFromDeath)
	}
}

// getKeys helper function to debug table structure
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
