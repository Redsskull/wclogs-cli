package services

import (
	"encoding/json"
	"fmt"
	"sort"

	"wclogs-cli/api"
	"wclogs-cli/models"

	"github.com/fatih/color"
)

// InterruptCorrelator handles the complex logic of matching cast events with interrupt events
type InterruptCorrelator struct {
	apiClient     *api.Client
	lookupService *LookupService
	reportCode    string
	fightID       int
	fightInfo     *models.FightInfo
	verbose       bool
}

// NewInterruptCorrelator creates a new interrupt correlator
func NewInterruptCorrelator(apiClient *api.Client, lookupService *LookupService, reportCode string, fightID int, fightInfo *models.FightInfo, verbose bool) *InterruptCorrelator {
	return &InterruptCorrelator{
		apiClient:     apiClient,
		lookupService: lookupService,
		reportCode:    reportCode,
		fightID:       fightID,
		fightInfo:     fightInfo,
		verbose:       verbose,
	}
}

// AnalyzeInterrupts performs comprehensive interrupt analysis
func (ic *InterruptCorrelator) AnalyzeInterrupts() (*models.InterruptAnalysisResult, error) {
	if ic.verbose {
		color.HiBlue("🔍 Starting comprehensive interrupt correlation analysis...")
	}

	// Create the analysis result
	analysis := models.NewInterruptAnalysisResult(ic.fightInfo)

	// Step 1: Fetch all interrupt events
	interrupts, err := ic.fetchInterruptEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch interrupt events: %w", err)
	}

	if ic.verbose {
		color.HiGreen("✅ Found %d interrupt events", len(interrupts))
	}

	// Step 2: Fetch enemy cast completion events (casts that went through)
	completedCasts, err := ic.fetchCastCompletionEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cast completion events: %w", err)
	}

	if ic.verbose {
		color.HiGreen("✅ Found %d enemy cast completion events", len(completedCasts))
	}

	// Step 3: Correlate the events
	if ic.verbose {
		color.HiBlue("🔗 Correlating interrupts with completed casts...")
	}

	err = ic.correlateEvents(analysis, interrupts, completedCasts)
	if err != nil {
		return nil, fmt.Errorf("failed to correlate events: %w", err)
	}

	// Step 5: Calculate rates and statistics
	analysis.CalculateRates()

	if ic.verbose {
		color.HiGreen("📊 Analysis complete: %d total casts, %d stopped (%.1f%%)",
			analysis.TotalCasts, analysis.TotalStops, analysis.OverallRate)
	}

	return analysis, nil
}

// fetchInterruptEvents retrieves all interrupt events for the fight
func (ic *InterruptCorrelator) fetchInterruptEvents() ([]*models.InterruptEventDetail, error) {
	var allInterrupts []*models.InterruptEventDetail
	var startTime *float64

	for {
		request := api.NewInterruptEventsRequest(ic.reportCode, ic.fightID, nil, startTime)
		response, err := ic.apiClient.Query(request.Query, request.Variables)
		if err != nil {
			return nil, fmt.Errorf("interrupt events API call failed: %w", err)
		}

		if response.Data == nil || response.Data.ReportData == nil ||
			response.Data.ReportData.Report == nil ||
			response.Data.ReportData.Report.Events == nil {
			break
		}

		// Parse interrupt events from JSON
		pageInterrupts, err := ic.parseInterruptEvents(response.Data.ReportData.Report.Events.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse interrupt events: %w", err)
		}

		allInterrupts = append(allInterrupts, pageInterrupts...)

		// Check for pagination
		if response.Data.ReportData.Report.Events.NextPageTimestamp == nil {
			break
		}
		startTime = response.Data.ReportData.Report.Events.NextPageTimestamp
	}

	return allInterrupts, nil
}

// fetchCastCompletionEvents retrieves all enemy cast completion events
func (ic *InterruptCorrelator) fetchCastCompletionEvents() ([]*models.CastCompleteEvent, error) {
	var allCompletedCasts []*models.CastCompleteEvent
	var startTime *float64

	for {
		request := api.NewAllCastEventsRequest(ic.reportCode, ic.fightID, api.EventHostilityHostile, startTime)
		response, err := ic.apiClient.Query(request.Query, request.Variables)
		if err != nil {
			return nil, fmt.Errorf("cast completion events API call failed: %w", err)
		}

		if response.Data == nil || response.Data.ReportData == nil ||
			response.Data.ReportData.Report == nil ||
			response.Data.ReportData.Report.Events == nil {
			break
		}

		// Parse cast completion events from JSON
		pageCompletedCasts, err := ic.parseCastCompleteEvents(response.Data.ReportData.Report.Events.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cast completion events: %w", err)
		}

		allCompletedCasts = append(allCompletedCasts, pageCompletedCasts...)

		// Check for pagination
		if response.Data.ReportData.Report.Events.NextPageTimestamp == nil {
			break
		}
		startTime = response.Data.ReportData.Report.Events.NextPageTimestamp
	}

	return allCompletedCasts, nil
}

// correlateEvents matches interrupts with completed casts to build analysis
func (ic *InterruptCorrelator) correlateEvents(analysis *models.InterruptAnalysisResult, interrupts []*models.InterruptEventDetail, completedCasts []*models.CastCompleteEvent) error {
	// Sort all events by timestamp for efficient correlation
	sort.Slice(interrupts, func(i, j int) bool {
		return interrupts[i].Timestamp < interrupts[j].Timestamp
	})
	sort.Slice(completedCasts, func(i, j int) bool {
		return completedCasts[i].Timestamp < completedCasts[j].Timestamp
	})

	// Create a map to track which casts were interrupted
	interruptedCasts := make(map[string]bool)

	// Process all interrupts - use the interrupted spell information from extraAbilityGameID
	for _, interrupt := range interrupts {
		spellName := interrupt.InterruptedSpellName
		spellID := interrupt.InterruptedSpellID

		// Fallback to interrupt ability name if interrupted spell is unknown
		if spellName == "" {
			spellName = interrupt.InterruptSpellName
			spellID = interrupt.InterruptSpellID
		}
		if spellName == "" {
			spellName = "Unknown Spell"
		}

		analysis.AddPlayerInterrupt(
			interrupt.InterrupterName,
			interrupt.InterrupterID,
			spellName,
			spellID,
			interrupt.Timestamp,
			interrupt.TargetName,
		)
	}

	// Process all completed casts - these are missed interrupts (excluding interrupted ones)
	for _, cast := range completedCasts {
		// Only consider interruptable spells from enemies
		if cast.SpellName != "" && models.IsInterruptableSpell(cast.SpellName, 2000) {
			castKey := fmt.Sprintf("%d-%d-%.3f", cast.CasterID, cast.SpellID, cast.Timestamp)

			// Skip casts that were actually interrupted
			if interruptedCasts[castKey] {
				continue
			}

			targetID := 0
			targetName := "Unknown"
			if cast.TargetID != nil {
				targetID = *cast.TargetID
				targetName = cast.TargetName
				if targetName == "" {
					targetName = ic.lookupService.GetActorName(targetID)
				}
			}

			analysis.AddMissedCast(
				cast.CasterName,
				cast.CasterID,
				targetName,
				targetID,
				cast.SpellName,
				cast.SpellID,
				cast.Timestamp,
				ic.fightInfo.StartTime,
			)
		}
	}

	return nil
}

// findInterruptedSpell attempts to find which spell cast was interrupted by correlating timing and targets
func (ic *InterruptCorrelator) findInterruptedSpell(interrupt *models.InterruptEventDetail, completedCasts []*models.CastCompleteEvent) *models.CastCompleteEvent {
	// Look for cast events from the same target (who was interrupted) within a reasonable time window
	// Interrupts typically happen during the cast, so we look slightly before the interrupt timestamp

	const lookbackWindow = 5000.0 // Look back 5 seconds maximum

	var candidateCasts []*models.CastCompleteEvent

	// Find all casts from the interrupted target in the time window
	for _, cast := range completedCasts {
		if cast.CasterID == interrupt.TargetID {
			// Cast must be within the lookback window before the interrupt
			if cast.Timestamp >= (interrupt.Timestamp-lookbackWindow) && cast.Timestamp <= interrupt.Timestamp {
				// Only consider interruptable spells
				if models.IsInterruptableSpell(cast.SpellName, 2000) {
					candidateCasts = append(candidateCasts, cast)
				}
			}
		}
	}

	// If we found candidates, return the most recent one (closest to interrupt time)
	if len(candidateCasts) > 0 {
		// Sort by timestamp descending (most recent first)
		sort.Slice(candidateCasts, func(i, j int) bool {
			return candidateCasts[i].Timestamp > candidateCasts[j].Timestamp
		})

		return candidateCasts[0]
	}

	return nil
}

// parseInterruptEvents parses interrupt events from raw JSON data
func (ic *InterruptCorrelator) parseInterruptEvents(data interface{}) ([]*models.InterruptEventDetail, error) {
	// Parse raw JSON to extract extraAbilityGameID
	rawEvents, err := ic.parseRawInterruptJSON(data)
	if err != nil {
		return nil, err
	}

	var interrupts []*models.InterruptEventDetail
	for _, raw := range rawEvents {
		if raw["type"] == "interrupt" {
			interrupt := &models.InterruptEventDetail{
				Timestamp:     raw["timestamp"].(float64),
				InterrupterID: int(raw["sourceID"].(float64)),
				TargetID:      int(raw["targetID"].(float64)),
			}

			// Get names from lookup service
			interrupt.InterrupterName = ic.lookupService.GetActorName(interrupt.InterrupterID)
			interrupt.TargetName = ic.lookupService.GetActorName(interrupt.TargetID)

			// Extract interrupt ability (what was used to interrupt)
			if abilityID, ok := raw["abilityGameID"].(float64); ok {
				interrupt.InterruptSpellID = int(abilityID)
				interrupt.InterruptSpellName = ic.lookupService.GetAbilityName(interrupt.InterruptSpellID)
			}

			// Extract interrupted spell (what was being cast)
			if extraAbilityID, ok := raw["extraAbilityGameID"].(float64); ok {
				interrupt.InterruptedSpellID = int(extraAbilityID)
				interrupt.InterruptedSpellName = ic.lookupService.GetAbilityName(interrupt.InterruptedSpellID)
			}

			interrupts = append(interrupts, interrupt)
		}
	}

	return interrupts, nil
}

// parseRawInterruptJSON parses raw JSON data into map format to access all fields
func (ic *InterruptCorrelator) parseRawInterruptJSON(data interface{}) ([]map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal interrupt data: %w", err)
	}

	var rawEvents []map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &rawEvents); err != nil {
		return nil, fmt.Errorf("failed to unmarshal interrupt data: %w", err)
	}

	return rawEvents, nil
}

// parseCastCompleteEvents parses cast completion events from raw JSON data
func (ic *InterruptCorrelator) parseCastCompleteEvents(data interface{}) ([]*models.CastCompleteEvent, error) {
	events, err := models.ParseEventsJSON(data)
	if err != nil {
		return nil, err
	}

	var completedCasts []*models.CastCompleteEvent
	for _, event := range events {
		if event.Type == "cast" {
			completed := &models.CastCompleteEvent{
				Timestamp: event.Timestamp,
				CasterID:  *event.SourceID,
			}

			if event.TargetID != nil {
				completed.TargetID = event.TargetID
				completed.TargetName = ic.lookupService.GetActorName(*event.TargetID)
			}

			// Get names from lookup service
			completed.CasterName = ic.lookupService.GetActorName(completed.CasterID)

			if event.AbilityID != nil {
				completed.SpellID = *event.AbilityID
				completed.SpellName = ic.lookupService.GetAbilityName(completed.SpellID)
			}

			completedCasts = append(completedCasts, completed)
		}
	}

	return completedCasts, nil
}
