package services

import (
	"fmt"
	"sync"

	"wclogs-cli/api"
)

// LookupService provides caching for ability and actor name lookups
type LookupService struct {
	apiClient    *api.Client
	abilityCache map[int]string // ability ID -> name
	actorCache   map[int]string // actor ID -> name
	cacheMutex   sync.RWMutex
}

// NewLookupService creates a new lookup service with caching
func NewLookupService(apiClient *api.Client) *LookupService {
	return &LookupService{
		apiClient:    apiClient,
		abilityCache: make(map[int]string),
		actorCache:   make(map[int]string),
	}
}

// GetAbilityName returns the ability name for the given ID, with caching
func (ls *LookupService) GetAbilityName(abilityID int) string {
	if abilityID == 0 {
		return "Unknown Ability"
	}

	// Check cache first (read lock)
	ls.cacheMutex.RLock()
	if name, exists := ls.abilityCache[abilityID]; exists {
		ls.cacheMutex.RUnlock()
		return name
	}
	ls.cacheMutex.RUnlock()

	// Not in cache, fetch from API
	name := ls.fetchAbilityName(abilityID)

	// Store in cache (write lock)
	ls.cacheMutex.Lock()
	ls.abilityCache[abilityID] = name
	ls.cacheMutex.Unlock()

	return name
}

// fetchAbilityName fetches ability name from the API
func (ls *LookupService) fetchAbilityName(abilityID int) string {
	request := api.NewAbilityLookupRequest(abilityID)
	response, err := ls.apiClient.Query(request.Query, request.Variables)
	if err != nil {
		return fmt.Sprintf("Ability ID %d", abilityID)
	}

	if response.Data == nil || response.Data.GameData == nil ||
		response.Data.GameData.Ability == nil {
		return fmt.Sprintf("Ability ID %d", abilityID)
	}

	if response.Data.GameData.Ability.Name == "" {
		return fmt.Sprintf("Ability ID %d", abilityID)
	}

	return response.Data.GameData.Ability.Name
}

// LoadActorsFromReport loads all actors (players, NPCs, pets) from report into cache
func (ls *LookupService) LoadActorsFromReport(reportCode string) error {
	request := api.NewAllActorsRequest(reportCode)
	response, err := ls.apiClient.Query(request.Query, request.Variables)
	if err != nil {
		return fmt.Errorf("failed to fetch actors: %w", err)
	}

	if response.Data == nil || response.Data.ReportData == nil ||
		response.Data.ReportData.Report == nil ||
		response.Data.ReportData.Report.MasterData == nil {
		return fmt.Errorf("no actor data found")
	}

	// Load all actors into cache
	ls.cacheMutex.Lock()
	defer ls.cacheMutex.Unlock()

	for _, actor := range response.Data.ReportData.Report.MasterData.Actors {
		ls.actorCache[actor.ID] = actor.Name
	}

	return nil
}

// GetActorName returns the actor name for the given ID
func (ls *LookupService) GetActorName(actorID int) string {
	if actorID == -1 {
		return "Environment"
	}

	ls.cacheMutex.RLock()
	defer ls.cacheMutex.RUnlock()

	if name, exists := ls.actorCache[actorID]; exists {
		return name
	}

	return fmt.Sprintf("Unknown Actor (ID %d)", actorID)
}

// GetPlayerLookup returns a map of player IDs to names (for backwards compatibility)
func (ls *LookupService) GetPlayerLookup() map[int]string {
	ls.cacheMutex.RLock()
	defer ls.cacheMutex.RUnlock()

	// Filter to only players
	playerLookup := make(map[int]string)
	for id, name := range ls.actorCache {
		// We'll need to differentiate players from NPCs somehow
		// For now, return all actors - the calling code can filter if needed
		playerLookup[id] = name
	}

	return playerLookup
}

// PreloadAbilities fetches multiple ability names in advance to reduce API calls
func (ls *LookupService) PreloadAbilities(abilityIDs []int) {
	var toFetch []int

	// Check which abilities we don't have cached
	ls.cacheMutex.RLock()
	for _, id := range abilityIDs {
		if id != 0 {
			if _, exists := ls.abilityCache[id]; !exists {
				toFetch = append(toFetch, id)
			}
		}
	}
	ls.cacheMutex.RUnlock()

	// Fetch missing abilities (could be optimized with batch requests)
	for _, abilityID := range toFetch {
		name := ls.fetchAbilityName(abilityID)

		ls.cacheMutex.Lock()
		ls.abilityCache[abilityID] = name
		ls.cacheMutex.Unlock()
	}
}

// GetCacheStats returns information about cached entries for debugging
func (ls *LookupService) GetCacheStats() (int, int) {
	ls.cacheMutex.RLock()
	defer ls.cacheMutex.RUnlock()

	return len(ls.abilityCache), len(ls.actorCache)
}

// FormatKillingInfo returns a formatted string for what killed the player
func (ls *LookupService) FormatKillingInfo(killerID *int, abilityID *int) (string, string) {
	var abilityName, sourceName string

	if abilityID != nil {
		abilityName = ls.GetAbilityName(*abilityID)
	} else {
		abilityName = "Unknown Ability"
	}

	if killerID != nil {
		sourceName = ls.GetActorName(*killerID)
	} else {
		sourceName = "Unknown Source"
	}

	return abilityName, sourceName
}
