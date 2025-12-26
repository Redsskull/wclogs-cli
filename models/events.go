package models

import (
	"encoding/json"
	"fmt"
)

// ParseEventsJSON parses raw JSON event data into Event structs
func ParseEventsJSON(data interface{}) ([]*Event, error) {
	// Convert to JSON bytes first
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Parse as array of raw events
	var rawEvents []map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &rawEvents); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	var events []*Event
	for _, rawEvent := range rawEvents {
		event, err := parseRawEvent(rawEvent)
		if err != nil {
			// Skip malformed events rather than failing entirely
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

// parseRawEvent converts a raw event map to an Event struct
func parseRawEvent(raw map[string]interface{}) (*Event, error) {
	event := &Event{}

	// Parse timestamp (required)
	if ts, ok := raw["timestamp"].(float64); ok {
		event.Timestamp = ts
	} else {
		return nil, fmt.Errorf("missing or invalid timestamp")
	}

	// Parse type (required)
	if eventType, ok := raw["type"].(string); ok {
		event.Type = eventType
	} else {
		return nil, fmt.Errorf("missing or invalid event type")
	}

	// Parse optional fields
	if sourceID, ok := raw["sourceID"].(float64); ok {
		id := int(sourceID)
		event.SourceID = &id
	}

	if targetID, ok := raw["targetID"].(float64); ok {
		id := int(targetID)
		event.TargetID = &id
	}

	if abilityID, ok := raw["abilityGameID"].(float64); ok {
		id := int(abilityID)
		event.AbilityID = &id
	}

	if amount, ok := raw["amount"].(float64); ok {
		amt := int(amount)
		event.Amount = &amt
	}

	if hitType, ok := raw["hitType"].(float64); ok {
		ht := int(hitType)
		event.HitType = &ht
	}

	if overkill, ok := raw["overkill"].(float64); ok {
		ok := int(overkill)
		event.Overkill = &ok
	}

	if tick, ok := raw["tick"].(bool); ok {
		event.Tick = &tick
	}

	// Death-specific fields
	if killerID, ok := raw["killerID"].(float64); ok {
		id := int(killerID)
		event.KillerID = &id
	}

	if killingAbilityID, ok := raw["killingAbilityGameID"].(float64); ok {
		id := int(killingAbilityID)
		event.KillingAbilityGameID = &id
	}

	// Parse ability information if present
	if abilityData, ok := raw["ability"].(map[string]interface{}); ok {
		event.Ability = parseEventAbility(abilityData)
	}

	// Parse source information if present
	if sourceData, ok := raw["source"].(map[string]interface{}); ok {
		event.Source = parseEventActor(sourceData)
	}

	// Parse target information if present
	if targetData, ok := raw["target"].(map[string]interface{}); ok {
		event.Target = parseEventActor(targetData)
	}

	return event, nil
}

// parseEventAbility parses ability information from raw data
func parseEventAbility(raw map[string]interface{}) *EventAbility {
	ability := &EventAbility{}

	if name, ok := raw["name"].(string); ok {
		ability.Name = name
	}

	if gameID, ok := raw["gameID"].(float64); ok {
		ability.GameID = int(gameID)
	}

	if abilityType, ok := raw["type"].(float64); ok {
		ability.Type = int(abilityType)
	}

	if icon, ok := raw["icon"].(string); ok {
		ability.Icon = icon
	}

	return ability
}

// parseEventActor parses actor information from raw data
func parseEventActor(raw map[string]any) *EventActor {
	actor := &EventActor{}

	if name, ok := raw["name"].(string); ok {
		actor.Name = name
	}

	if id, ok := raw["id"].(float64); ok {
		actor.ID = int(id)
	}

	if actorType, ok := raw["type"].(string); ok {
		actor.Type = actorType
	}

	if icon, ok := raw["icon"].(string); ok {
		actor.Icon = icon
	}

	return actor
}

// ParseCastEventsJSON parses cast events from JSON data
func ParseCastEventsJSON(data any) ([]*Event, error) {
	events, err := ParseEventsJSON(data)
	if err != nil {
		return nil, err
	}

	// Filter to only cast events
	var castEvents []*Event
	for _, event := range events {
		if event.Type == "cast" {
			castEvents = append(castEvents, event)
		}
	}

	return castEvents, nil
}

// ParseInterruptEventsJSON parses interrupt events from JSON data
func ParseInterruptEventsJSON(data interface{}) ([]*Event, error) {
	events, err := ParseEventsJSON(data)
	if err != nil {
		return nil, err
	}

	// Filter to only interrupt events
	var interruptEvents []*Event
	for _, event := range events {
		if event.Type == "interrupt" {
			interruptEvents = append(interruptEvents, event)
		}
	}

	return interruptEvents, nil
}

// ParseBeginCastEventsJSON parses begin cast events from JSON data
func ParseBeginCastEventsJSON(data interface{}) ([]*Event, error) {
	events, err := ParseEventsJSON(data)
	if err != nil {
		return nil, err
	}

	// Filter to only begin cast events
	var beginCastEvents []*Event
	for _, event := range events {
		if event.Type == "begincast" {
			beginCastEvents = append(beginCastEvents, event)
		}
	}

	return beginCastEvents, nil
}

// ParseDeathEventsJSON parses death events from JSON data
func ParseDeathEventsJSON(data interface{}) ([]*Event, error) {
	events, err := ParseEventsJSON(data)
	if err != nil {
		return nil, err
	}

	// Filter to only death events
	var deathEvents []*Event
	for _, event := range events {
		if event.Type == "death" {
			deathEvents = append(deathEvents, event)
		}
	}

	return deathEvents, nil
}

// ParseDamageTakenEventsJSON parses damage taken events from JSON data
func ParseDamageTakenEventsJSON(data interface{}) ([]*Event, error) {
	events, err := ParseEventsJSON(data)
	if err != nil {
		return nil, err
	}

	// Filter to only damage events
	var damageEvents []*Event
	for _, event := range events {
		if event.Type == "damage" {
			damageEvents = append(damageEvents, event)
		}
	}

	return damageEvents, nil
}

// ParseHealingEventsJSON parses healing events from JSON data
func ParseHealingEventsJSON(data interface{}) ([]*Event, error) {
	events, err := ParseEventsJSON(data)
	if err != nil {
		return nil, err
	}

	// Filter to only healing events
	var healingEvents []*Event
	for _, event := range events {
		if event.Type == "heal" {
			healingEvents = append(healingEvents, event)
		}
	}

	return healingEvents, nil
}

// GetEventsByType filters events by their type
func GetEventsByType(events []*Event, eventType string) []*Event {
	var filtered []*Event
	for _, event := range events {
		if event.Type == eventType {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// GetEventsInTimeRange filters events by timestamp range
func GetEventsInTimeRange(events []*Event, startTime, endTime float64) []*Event {
	var filtered []*Event
	for _, event := range events {
		if event.Timestamp >= startTime && event.Timestamp <= endTime {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// GetEventsByActor filters events by source or target actor ID
func GetEventsByActor(events []*Event, actorID int, checkSource, checkTarget bool) []*Event {
	var filtered []*Event
	for _, event := range events {
		if checkSource && event.SourceID != nil && *event.SourceID == actorID {
			filtered = append(filtered, event)
		}
		if checkTarget && event.TargetID != nil && *event.TargetID == actorID {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// GetEventsByAbility filters events by ability ID
func GetEventsByAbility(events []*Event, abilityID int) []*Event {
	var filtered []*Event
	for _, event := range events {
		if event.AbilityID != nil && *event.AbilityID == abilityID {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// ValidateEventData performs basic validation on event data
func ValidateEventData(events []*Event) ([]string, []string) {
	var warnings []string
	var errors []string

	if len(events) == 0 {
		warnings = append(warnings, "No events found")
		return warnings, errors
	}

	// Check for required fields
	missingTimestamps := 0
	missingTypes := 0

	for _, event := range events {
		if event.Timestamp == 0 {
			missingTimestamps++
		}
		if event.Type == "" {
			missingTypes++
		}
	}

	if missingTimestamps > 0 {
		warnings = append(warnings, fmt.Sprintf("%d events missing timestamps", missingTimestamps))
	}

	if missingTypes > 0 {
		errors = append(errors, fmt.Sprintf("%d events missing event types", missingTypes))
	}

	return warnings, errors
}
