package models

import (
	"testing"
)

func TestNewGraphQLResponse(t *testing.T) {
	response := NewGraphQLResponse()
	
	if response.Data == nil {
		t.Error("NewGraphQLResponse() Data should not be nil")
	}
	
	if response.Errors == nil {
		t.Error("NewGraphQLResponse() Errors should not be nil")
	}
	
	if len(response.Errors) != 0 {
		t.Errorf("NewGraphQLResponse() Errors should be empty, got %d", len(response.Errors))
	}
}

func TestNewPlayer(t *testing.T) {
	player := NewPlayer("TestPlayer", "Warrior", 1000.0, "icon_url")
	
	if player.Name != "TestPlayer" {
		t.Errorf("NewPlayer() Name = %v, expected %v", player.Name, "TestPlayer")
	}
	
	if player.Class != "Warrior" {
		t.Errorf("NewPlayer() Class = %v, expected %v", player.Class, "Warrior")
	}
	
	if player.Total != 1000.0 {
		t.Errorf("NewPlayer() Total = %v, expected %v", player.Total, 1000.0)
	}
	
	if player.Icon != "icon_url" {
		t.Errorf("NewPlayer() Icon = %v, expected %v", player.Icon, "icon_url")
	}
}

func TestNewTableData(t *testing.T) {
	tableData := NewTableData()
	
	if tableData.Entries == nil {
		t.Error("NewTableData() Entries should not be nil")
	}
	
	if len(tableData.Entries) != 0 {
		t.Errorf("NewTableData() Entries should be empty, got %d", len(tableData.Entries))
	}
}

func TestGraphQLResponseIsValid(t *testing.T) {
	tests := []struct {
		name     string
		response GraphQLResponse
		expected bool
	}{
		{
			name: "valid response",
			response: GraphQLResponse{
				Data:   &ResponseData{},
				Errors: []GraphQLError{},
			},
			expected: true,
		},
		{
			name: "nil data",
			response: GraphQLResponse{
				Data:   nil,
				Errors: []GraphQLError{},
			},
			expected: false,
		},
		{
			name: "with errors",
			response: GraphQLResponse{
				Data: &ResponseData{},
				Errors: []GraphQLError{
					{Message: "test error"},
				},
			},
			expected: false,
		},
		{
			name: "nil data and errors",
			response: GraphQLResponse{
				Data:   nil,
				Errors: []GraphQLError{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.IsValid()
			if result != tt.expected {
				t.Errorf("GraphQLResponse.IsValid() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGraphQLResponseHasErrors(t *testing.T) {
	tests := []struct {
		name     string
		response GraphQLResponse
		expected bool
	}{
		{
			name: "no errors",
			response: GraphQLResponse{
				Errors: []GraphQLError{},
			},
			expected: false,
		},
		{
			name: "with errors",
			response: GraphQLResponse{
				Errors: []GraphQLError{
					{Message: "test error"},
				},
			},
			expected: true,
		},
		{
			name: "multiple errors",
			response: GraphQLResponse{
				Errors: []GraphQLError{
					{Message: "test error 1"},
					{Message: "test error 2"},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.HasErrors()
			if result != tt.expected {
				t.Errorf("GraphQLResponse.HasErrors() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGraphQLResponseFirstError(t *testing.T) {
	tests := []struct {
		name     string
		response GraphQLResponse
		expected string
	}{
		{
			name: "no errors",
			response: GraphQLResponse{
				Errors: []GraphQLError{},
			},
			expected: "",
		},
		{
			name: "one error",
			response: GraphQLResponse{
				Errors: []GraphQLError{
					{Message: "test error"},
				},
			},
			expected: "test error",
		},
		{
			name: "multiple errors",
			response: GraphQLResponse{
				Errors: []GraphQLError{
					{Message: "first error"},
					{Message: "second error"},
				},
			},
			expected: "first error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.response.FirstError()
			if result != tt.expected {
				t.Errorf("GraphQLResponse.FirstError() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestNewPlayerFromEntry(t *testing.T) {
	entry := &PlayerEntry{
		Name:      "TestPlayer",
		Type:      "Warrior",
		Total:     1000.0,
		Icon:      "icon_url",
		ItemLevel: 380,
	}

	player := NewPlayerFromEntry(entry)

	if player.Name != "TestPlayer" {
		t.Errorf("NewPlayerFromEntry() Name = %v, expected %v", player.Name, "TestPlayer")
	}

	if player.Class != "Warrior" {
		t.Errorf("NewPlayerFromEntry() Class = %v, expected %v", player.Class, "Warrior")
	}

	if player.Total != 1000.0 {
		t.Errorf("NewPlayerFromEntry() Total = %v, expected %v", player.Total, 1000.0)
	}

	if player.Icon != "icon_url" {
		t.Errorf("NewPlayerFromEntry() Icon = %v, expected %v", player.Icon, "icon_url")
	}

	if player.ItemLevel != 380 {
		t.Errorf("NewPlayerFromEntry() ItemLevel = %v, expected %v", player.ItemLevel, 380)
	}

	// Note: The DPS field is set during the model creation but not calculated from the entry
	// The actual calculation might happen elsewhere
}

func TestPlayerLookup(t *testing.T) {
	// Create test PlayerInfo objects
	player1 := &PlayerInfo{ID: 1, Name: "Player1", Class: "Warrior", Server: "Area 52", Icon: "warrior_icon"}
	player2 := &PlayerInfo{ID: 2, Name: "Player2", Class: "Mage", Server: "Area 52", Icon: "mage_icon"}

	// Create a PlayerLookup manually since there's no constructor in the provided types
	lookup := &PlayerLookup{
		playersByName: map[string]*PlayerInfo{
			"player1": player1,
			"player2": player2,
		},
		playersByID: map[int]*PlayerInfo{
			1: player1,
			2: player2,
		},
	}

	// Test that the lookup data is stored properly
	if len(lookup.playersByName) != 2 {
		t.Errorf("PlayerLookup should have 2 players in playersByName, got %d", len(lookup.playersByName))
	}

	if len(lookup.playersByID) != 2 {
		t.Errorf("PlayerLookup should have 2 players in playersByID, got %d", len(lookup.playersByID))
	}

	// Test that entries can be retrieved
	if lookup.playersByName["player1"] != player1 {
		t.Error("PlayerLookup playersByName mapping is incorrect")
	}

	if lookup.playersByID[1] != player1 {
		t.Error("PlayerLookup playersByID mapping is incorrect")
	}
}

func TestEventStruct(t *testing.T) {
	// Test the Event struct with some sample data
	sourceID := 123
	targetID := 456
	abilityID := 789
	amount := 50000

	event := Event{
		Timestamp: 1000.5,
		Type:      "damage",
		SourceID:  &sourceID,
		TargetID:  &targetID,
		AbilityID: &abilityID,
		Amount:    &amount,
	}

	if event.Timestamp != 1000.5 {
		t.Errorf("Event Timestamp = %v, expected %v", event.Timestamp, 1000.5)
	}

	if event.Type != "damage" {
		t.Errorf("Event Type = %v, expected %v", event.Type, "damage")
	}

	if event.SourceID == nil || *event.SourceID != 123 {
		t.Errorf("Event SourceID = %v, expected %v", event.SourceID, 123)
	}

	if event.TargetID == nil || *event.TargetID != 456 {
		t.Errorf("Event TargetID = %v, expected %v", event.TargetID, 456)
	}

	if event.AbilityID == nil || *event.AbilityID != 789 {
		t.Errorf("Event AbilityID = %v, expected %v", event.AbilityID, 789)
	}

	if event.Amount == nil || *event.Amount != 50000 {
		t.Errorf("Event Amount = %v, expected %v", event.Amount, 50000)
	}
}

func TestDeathEventStruct(t *testing.T) {
	ability := &EventAbility{Name: "Crystalline Shockwave", GameID: 1226823, Type: 1, Icon: "icon"}
	source := &EventActor{Name: "Fractillus", ID: 24, Type: "NPC", Icon: "boss_icon"}

	deathEvent := DeathEvent{
		PlayerID:      123,
		PlayerName:    "TestPlayer",
		Timestamp:     1000.5,
		KillingAbility: ability,
		KillingSource:  source,
		Overkill:      10000,
	}

	if deathEvent.PlayerID != 123 {
		t.Errorf("DeathEvent PlayerID = %v, expected %v", deathEvent.PlayerID, 123)
	}

	if deathEvent.PlayerName != "TestPlayer" {
		t.Errorf("DeathEvent PlayerName = %v, expected %v", deathEvent.PlayerName, "TestPlayer")
	}

	if deathEvent.Timestamp != 1000.5 {
		t.Errorf("DeathEvent Timestamp = %v, expected %v", deathEvent.Timestamp, 1000.5)
	}

	if deathEvent.KillingAbility == nil || deathEvent.KillingAbility.Name != "Crystalline Shockwave" {
		t.Errorf("DeathEvent KillingAbility.Name = %v, expected %v", 
			deathEvent.KillingAbility.Name, "Crystalline Shockwave")
	}

	if deathEvent.KillingSource == nil || deathEvent.KillingSource.Name != "Fractillus" {
		t.Errorf("DeathEvent KillingSource.Name = %v, expected %v", 
			deathEvent.KillingSource.Name, "Fractillus")
	}

	if deathEvent.Overkill != 10000 {
		t.Errorf("DeathEvent Overkill = %v, expected %v", deathEvent.Overkill, 10000)
	}
}