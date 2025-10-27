package api

import (
	"testing"
)

func TestNewTableRequest(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		fightID  int
		dataType DataType
		expectedQuery string
	}{
		{
			name:     "damage request",
			code:     "ABC123",
			fightID:  5,
			dataType: DataTypeDamage,
			expectedQuery: DamageTableQuery,
		},
		{
			name:     "healing request",
			code:     "XYZ789",
			fightID:  3,
			dataType: DataTypeHealing,
			expectedQuery: HealingTableQuery,
		},
		{
			name:     "unknown data type defaults to damage",
			code:     "DEF456",
			fightID:  1,
			dataType: DataType("InvalidType"), // Invalid DataType
			expectedQuery: DamageTableQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := NewTableRequest(tt.code, tt.fightID, tt.dataType)

			if request.Query != tt.expectedQuery {
				t.Errorf("NewTableRequest() Query = %v, expected %v", request.Query, tt.expectedQuery)
			}

			if variables, ok := request.Variables["code"].(string); !ok || variables != tt.code {
				t.Errorf("NewTableRequest() code variable = %v, expected %v", request.Variables["code"], tt.code)
			}

			if variables, ok := request.Variables["fightID"].(int); !ok || variables != tt.fightID {
				t.Errorf("NewTableRequest() fightID variable = %v, expected %v", request.Variables["fightID"], tt.fightID)
			}
		})
	}
}

func TestNewMasterDataRequest(t *testing.T) {
	code := "ABC123"
	request := NewMasterDataRequest(code)

	if request.Query != MasterDataQuery {
		t.Errorf("NewMasterDataRequest() Query = %v, expected %v", request.Query, MasterDataQuery)
	}

	if variables, ok := request.Variables["code"].(string); !ok || variables != code {
		t.Errorf("NewMasterDataRequest() code variable = %v, expected %v", request.Variables["code"], code)
	}
}

func TestNewFightInfoRequest(t *testing.T) {
	code := "XYZ789"
	request := NewFightInfoRequest(code)

	if request.Query != FightInfoQuery {
		t.Errorf("NewFightInfoRequest() Query = %v, expected %v", request.Query, FightInfoQuery)
	}

	if variables, ok := request.Variables["code"].(string); !ok || variables != code {
		t.Errorf("NewFightInfoRequest() code variable = %v, expected %v", request.Variables["code"], code)
	}
}

func TestNewAllActorsRequest(t *testing.T) {
	code := "DEF456"
	request := NewAllActorsRequest(code)

	if request.Query != AllActorsQuery {
		t.Errorf("NewAllActorsRequest() Query = %v, expected %v", request.Query, AllActorsQuery)
	}

	if variables, ok := request.Variables["code"].(string); !ok || variables != code {
		t.Errorf("NewAllActorsRequest() code variable = %v, expected %v", request.Variables["code"], code)
	}
}

func TestNewAbilityLookupRequest(t *testing.T) {
	abilityID := 12345
	request := NewAbilityLookupRequest(abilityID)

	if request.Query != SingleAbilityLookupQuery {
		t.Errorf("NewAbilityLookupRequest() Query = %v, expected %v", request.Query, SingleAbilityLookupQuery)
	}

	if variables, ok := request.Variables["abilityID"].(int); !ok || variables != abilityID {
		t.Errorf("NewAbilityLookupRequest() abilityID variable = %v, expected %v", request.Variables["abilityID"], abilityID)
	}
}

func TestNewDeathEventsRequest(t *testing.T) {
	code := "ABC123"
	fightID := 5
	playerID := 123

	// Test with player ID
	request := NewDeathEventsRequest(code, fightID, &playerID)

	if request.Query != DeathEventsQuery {
		t.Errorf("NewDeathEventsRequest() Query = %v, expected %v", request.Query, DeathEventsQuery)
	}

	if variables, ok := request.Variables["code"].(string); !ok || variables != code {
		t.Errorf("NewDeathEventsRequest() code variable = %v, expected %v", request.Variables["code"], code)
	}

	if variables, ok := request.Variables["fightID"].(int); !ok || variables != fightID {
		t.Errorf("NewDeathEventsRequest() fightID variable = %v, expected %v", request.Variables["fightID"], fightID)
	}

	if variables, ok := request.Variables["playerID"].(int); !ok || variables != playerID {
		t.Errorf("NewDeathEventsRequest() playerID variable = %v, expected %v", request.Variables["playerID"], playerID)
	}

	// Test without player ID
	requestNoPlayer := NewDeathEventsRequest(code, fightID, nil)

	if _, exists := requestNoPlayer.Variables["playerID"]; exists {
		t.Error("NewDeathEventsRequest() should not have playerID variable when playerID is nil")
	}
}

func TestNewHealingReceivedRequest(t *testing.T) {
	code := "TEST123"
	fightID := 7
	playerID := 456
	startTime := 1000.0
	endTime := 2000.0

	request := NewHealingReceivedRequest(code, fightID, playerID, startTime, endTime)

	if request.Query != HealingReceivedBeforeDeathQuery {
		t.Errorf("NewHealingReceivedRequest() Query = %v, expected %v", request.Query, HealingReceivedBeforeDeathQuery)
	}

	if variables, ok := request.Variables["code"].(string); !ok || variables != code {
		t.Errorf("NewHealingReceivedRequest() code variable = %v, expected %v", request.Variables["code"], code)
	}

	if variables, ok := request.Variables["fightID"].(int); !ok || variables != fightID {
		t.Errorf("NewHealingReceivedRequest() fightID variable = %v, expected %v", request.Variables["fightID"], fightID)
	}

	if variables, ok := request.Variables["playerID"].(int); !ok || variables != playerID {
		t.Errorf("NewHealingReceivedRequest() playerID variable = %v, expected %v", request.Variables["playerID"], playerID)
	}

	if variables, ok := request.Variables["startTime"].(float64); !ok || variables != startTime {
		t.Errorf("NewHealingReceivedRequest() startTime variable = %v, expected %v", request.Variables["startTime"], startTime)
	}

	if variables, ok := request.Variables["endTime"].(float64); !ok || variables != endTime {
		t.Errorf("NewHealingReceivedRequest() endTime variable = %v, expected %v", request.Variables["endTime"], endTime)
	}
}

func TestNewDamageTakenRequest(t *testing.T) {
	code := "DAMAGE456"
	fightID := 3
	playerID := 789
	startTime := 500.0
	endTime := 1500.0

	request := NewDamageTakenRequest(code, fightID, playerID, startTime, endTime)

	if request.Query != DamageTakenBeforeDeathQuery {
		t.Errorf("NewDamageTakenRequest() Query = %v, expected %v", request.Query, DamageTakenBeforeDeathQuery)
	}

	if variables, ok := request.Variables["code"].(string); !ok || variables != code {
		t.Errorf("NewDamageTakenRequest() code variable = %v, expected %v", request.Variables["code"], code)
	}

	if variables, ok := request.Variables["fightID"].(int); !ok || variables != fightID {
		t.Errorf("NewDamageTakenRequest() fightID variable = %v, expected %v", request.Variables["fightID"], fightID)
	}

	if variables, ok := request.Variables["playerID"].(int); !ok || variables != playerID {
		t.Errorf("NewDamageTakenRequest() playerID variable = %v, expected %v", request.Variables["playerID"], playerID)
	}

	if variables, ok := request.Variables["startTime"].(float64); !ok || variables != startTime {
		t.Errorf("NewDamageTakenRequest() startTime variable = %v, expected %v", request.Variables["startTime"], startTime)
	}

	if variables, ok := request.Variables["endTime"].(float64); !ok || variables != endTime {
		t.Errorf("NewDamageTakenRequest() endTime variable = %v, expected %v", request.Variables["endTime"], endTime)
	}
}

func TestNewDefensiveAbilitiesRequest(t *testing.T) {
	code := "DEFENSE789"
	fightID := 9
	playerID := 999
	startTime := 200.0
	endTime := 1200.0

	request := NewDefensiveAbilitiesRequest(code, fightID, playerID, startTime, endTime)

	if request.Query != DefensiveAbilitiesBeforeDeathQuery {
		t.Errorf("NewDefensiveAbilitiesRequest() Query = %v, expected %v", request.Query, DefensiveAbilitiesBeforeDeathQuery)
	}

	if variables, ok := request.Variables["code"].(string); !ok || variables != code {
		t.Errorf("NewDefensiveAbilitiesRequest() code variable = %v, expected %v", request.Variables["code"], code)
	}

	if variables, ok := request.Variables["fightID"].(int); !ok || variables != fightID {
		t.Errorf("NewDefensiveAbilitiesRequest() fightID variable = %v, expected %v", request.Variables["fightID"], fightID)
	}

	if variables, ok := request.Variables["playerID"].(int); !ok || variables != playerID {
		t.Errorf("NewDefensiveAbilitiesRequest() playerID variable = %v, expected %v", request.Variables["playerID"], playerID)
	}

	if variables, ok := request.Variables["startTime"].(float64); !ok || variables != startTime {
		t.Errorf("NewDefensiveAbilitiesRequest() startTime variable = %v, expected %v", request.Variables["startTime"], startTime)
	}

	if variables, ok := request.Variables["endTime"].(float64); !ok || variables != endTime {
		t.Errorf("NewDefensiveAbilitiesRequest() endTime variable = %v, expected %v", request.Variables["endTime"], endTime)
	}
}

func TestGraphQLRequestStructure(t *testing.T) {
	request := &GraphQLRequest{
		Query: "test query",
		Variables: map[string]any{
			"test": "value",
		},
	}

	if request.Query != "test query" {
		t.Errorf("GraphQLRequest Query = %v, expected %v", request.Query, "test query")
	}

	if request.Variables == nil {
		t.Error("GraphQLRequest Variables should not be nil")
	}

	if val, ok := request.Variables["test"].(string); !ok || val != "value" {
		t.Errorf("GraphQLRequest Variables should contain test = value, got %v", request.Variables["test"])
	}
}

