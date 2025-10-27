package cmd

import (
	"testing"
	"wclogs-cli/models"
)

func TestFilterPlayersByName(t *testing.T) {
	// Create test players
	players := []*models.Player{
		{Name: "Pherally", Class: "Warrior", Total: 1000.0},
		{Name: "Hanahime", Class: "Monk", Total: 900.0},
		{Name: "Nikkans", Class: "Paladin", Total: 800.0},
		{Name: "TEKKYysp", Class: "Shaman", Total: 700.0}, // Test case-insensitive matching
	}

	tests := []struct {
		name     string
		target   string
		expected int
	}{
		{
			name:     "exact match",
			target:   "Pherally",
			expected: 1,
		},
		{
			name:     "case insensitive match",
			target:   "nikkans",
			expected: 1,
		},
		{
			name:     "case insensitive match with mixed case",
			target:   "tekkYysp",
			expected: 1,
		},
		{
			name:     "no match",
			target:   "NonExistent",
			expected: 0,
		},
		{
			name:     "partial match should not match",
			target:   "Phera",
			expected: 0, // Should not match "Pherally" because it's exact match only
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterPlayersByName(players, tt.target)
			if len(result) != tt.expected {
				t.Errorf("filterPlayersByName() returned %d players, expected %d", len(result), tt.expected)
			}
			
			if tt.expected > 0 && len(result) > 0 {
				// Verify the matched player has the correct name (case-insensitive)
				found := false
				for _, player := range players {
					if player.Name == result[0].Name && 
						toLowerCase(player.Name) == toLowerCase(tt.target) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("filterPlayersByName() returned unexpected player: %v", result[0].Name)
				}
			}
		})
	}
}

// Helper function to convert string to lowercase
func toLowerCase(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

func TestTableTypesInfo(t *testing.T) {
	// Test that the tableTypes map is properly defined
	damageInfo, exists := tableTypes["damage"]
	if !exists {
		t.Error("tableTypes should contain 'damage' key")
	} else {
		if damageInfo.Description != "damage done" {
			t.Errorf("damage Description = %v, expected %v", damageInfo.Description, "damage done")
		}
		if damageInfo.Title != "DAMAGE TABLE" {
			t.Errorf("damage Title = %v, expected %v", damageInfo.Title, "DAMAGE TABLE")
		}
		if damageInfo.Emoji != "🗡️" {
			t.Errorf("damage Emoji = %v, expected %v", damageInfo.Emoji, "🗡️")
		}
	}

	healingInfo, exists := tableTypes["healing"]
	if !exists {
		t.Error("tableTypes should contain 'healing' key")
	} else {
		if healingInfo.Description != "healing done" {
			t.Errorf("healing Description = %v, expected %v", healingInfo.Description, "healing done")
		}
		if healingInfo.Title != "HEALING TABLE" {
			t.Errorf("healing Title = %v, expected %v", healingInfo.Title, "HEALING TABLE")
		}
		if healingInfo.Emoji != "💚" {
			t.Errorf("healing Emoji = %v, expected %v", healingInfo.Emoji, "💚")
		}
	}

	// Test non-existent type
	_, exists = tableTypes["nonexistent"]
	if exists {
		t.Error("tableTypes should not contain 'nonexistent' key")
	}
}

func TestCreateTableHandler(t *testing.T) {
	// Test that the function returns a valid handler
	handler := createTableHandler("damage")
	
	// The handler is a function, so we can't directly check its contents
	// But we can at least verify it's not nil
	if handler == nil {
		t.Error("createTableHandler() should not return nil")
	}
	
	// Test with different table types
	damageHandler := createTableHandler("damage")
	healingHandler := createTableHandler("healing")
	
	// Both should be valid functions (though we can't compare functions directly)
	if damageHandler == nil || healingHandler == nil {
		t.Error("createTableHandler() should return valid handlers for known table types")
	}
}

// Test the internal TableInfo struct
func TestTableInfoStruct(t *testing.T) {
	info := TableInfo{
		DataType:    "DamageDone",
		Description: "damage done",
		Title:       "DAMAGE TABLE",
		Emoji:       "🗡️",
	}
	
	if info.DataType != "DamageDone" {
		t.Errorf("TableInfo DataType = %v, expected %v", info.DataType, "DamageDone")
	}
	
	if info.Description != "damage done" {
		t.Errorf("TableInfo Description = %v, expected %v", info.Description, "damage done")
	}
	
	if info.Title != "DAMAGE TABLE" {
		t.Errorf("TableInfo Title = %v, expected %v", info.Title, "DAMAGE TABLE")
	}
	
	if info.Emoji != "🗡️" {
		t.Errorf("TableInfo Emoji = %v, expected %v", info.Emoji, "🗡️")
	}
}