package services

import (
	"testing"
)

func TestNewLookupService(t *testing.T) {
	// Note: We can't easily test with a real API client due to complex dependencies
	// Instead, we'll test the internal functionality directly
	lookupService := NewLookupService(nil) // We pass nil because we're only testing cache functionality
	
	if lookupService == nil {
		t.Error("NewLookupService() should not return nil")
	}
	
	if lookupService.abilityCache == nil {
		t.Error("NewLookupService() should initialize abilityCache")
	}
	
	if lookupService.actorCache == nil {
		t.Error("NewLookupService() should initialize actorCache")
	}
	
	// We just ensure the function runs without error
}

func TestGetActorName(t *testing.T) {
	// Test the GetActorName method directly using internal cache manipulation
	lookupService := NewLookupService(nil)
	
	// Pre-populate actor cache for testing
	lookupService.cacheMutex.Lock()
	lookupService.actorCache[1] = "Test Actor"
	lookupService.actorCache[2] = "Another Actor"
	lookupService.cacheMutex.Unlock()
	
	// Test known actor
	actorName := lookupService.GetActorName(1)
	if actorName != "Test Actor" {
		t.Errorf("GetActorName(1) = %v, expected %v", actorName, "Test Actor")
	}
	
	// Test environment actor
	envName := lookupService.GetActorName(-1)
	if envName != "Environment" {
		t.Errorf("GetActorName(-1) = %v, expected %v", envName, "Environment")
	}
	
	// Test unknown actor
	unknownName := lookupService.GetActorName(999)
	expectedUnknown := "Unknown Actor (ID 999)"
	if unknownName != expectedUnknown {
		t.Errorf("GetActorName(999) = %v, expected %v", unknownName, expectedUnknown)
	}
}

func TestFormatKillingInfo(t *testing.T) {
	// Test the nil case which should return default values
	lookupService := NewLookupService(nil)
	
	nilAbilityName, nilSourceName := lookupService.FormatKillingInfo(nil, nil)
	if nilAbilityName != "Unknown Ability" {
		t.Errorf("FormatKillingInfo(nil, nil) abilityName = %v, expected %v", nilAbilityName, "Unknown Ability")
	}
	if nilSourceName != "Unknown Source" {
		t.Errorf("FormatKillingInfo(nil, nil) sourceName = %v, expected %v", nilSourceName, "Unknown Source")
	}
}

func TestGetCacheStats(t *testing.T) {
	lookupService := NewLookupService(nil)
	
	// Initially empty
	abilityCount, actorCount := lookupService.GetCacheStats()
	if abilityCount != 0 {
		t.Errorf("Initial ability cache count = %d, expected %d", abilityCount, 0)
	}
	if actorCount != 0 {
		t.Errorf("Initial actor cache count = %d, expected %d", actorCount, 0)
	}
	
	// Add items to cache
	lookupService.cacheMutex.Lock()
	lookupService.abilityCache[1] = "Test Ability"
	lookupService.abilityCache[2] = "Another Ability"
	lookupService.actorCache[1] = "Test Actor"
	lookupService.actorCache[2] = "Another Actor"
	lookupService.actorCache[3] = "Third Actor"
	lookupService.cacheMutex.Unlock()
	
	// Check counts after adding
	abilityCount, actorCount = lookupService.GetCacheStats()
	if abilityCount != 2 {
		t.Errorf("After adding, ability cache count = %d, expected %d", abilityCount, 2)
	}
	if actorCount != 3 {
		t.Errorf("After adding, actor cache count = %d, expected %d", actorCount, 3)
	}
}

func TestGetPlayerLookup(t *testing.T) {
	lookupService := NewLookupService(nil)
	
	// Pre-populate actor cache
	lookupService.cacheMutex.Lock()
	lookupService.actorCache[1] = "Player1"
	lookupService.actorCache[2] = "Player2"
	lookupService.actorCache[100] = "NPC1"
	lookupService.cacheMutex.Unlock()
	
	playerLookup := lookupService.GetPlayerLookup()
	
	if len(playerLookup) == 0 {
		t.Error("GetPlayerLookup() should return some players")
	}
	
	// Check that the actors we added are in the lookup
	if name, exists := playerLookup[1]; !exists || name != "Player1" {
		t.Errorf("GetPlayerLookup()[1] = %v, expected %v", name, "Player1")
	}
	
	if name, exists := playerLookup[2]; !exists || name != "Player2" {
		t.Errorf("GetPlayerLookup()[2] = %v, expected %v", name, "Player2")
	}
}

func TestCacheThreadSafety(t *testing.T) {
	lookupService := NewLookupService(nil)
	
	// Test concurrent access to cache
	done := make(chan bool)
	
	// Goroutine 1: Add to ability cache
	go func() {
		for i := 0; i < 100; i++ {
			lookupService.cacheMutex.Lock()
			lookupService.abilityCache[i] = "Ability " + string(rune('0'+(i%10)))
			lookupService.cacheMutex.Unlock()
		}
		done <- true
	}()
	
	// Goroutine 2: Add to actor cache
	go func() {
		for i := 0; i < 100; i++ {
			lookupService.cacheMutex.Lock()
			lookupService.actorCache[i+100] = "Actor " + string(rune('0'+(i%10)))
			lookupService.cacheMutex.Unlock()
		}
		done <- true
	}()
	
	// Goroutine 3: Read cache stats
	go func() {
		for i := 0; i < 100; i++ {
			_, _ = lookupService.GetCacheStats()
		}
		done <- true
	}()
	
	// Wait for all goroutines to complete
	<-done
	<-done
	<-done
	
	// Check final state
	abilityCount, actorCount := lookupService.GetCacheStats()
	if abilityCount == 0 {
		t.Error("Ability cache should have some entries after concurrent operations")
	}
	if actorCount == 0 {
		t.Error("Actor cache should have some entries after concurrent operations")
	}
}