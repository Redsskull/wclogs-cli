package cmd

import (
	"testing"
)

func TestTableInfo(t *testing.T) {
	// Test that the TableInfo struct fields are properly accessible
	info := TableInfo{
		DataType:    "DamageDone",
		Description: "damage done",
		Title:       "DAMAGE TABLE",
		Emoji:       "🗡️",
	}

	if info.DataType != "DamageDone" {
		t.Errorf("TableInfo.DataType = %v, expected %v", info.DataType, "DamageDone")
	}

	if info.Description != "damage done" {
		t.Errorf("TableInfo.Description = %v, expected %v", info.Description, "damage done")
	}

	if info.Title != "DAMAGE TABLE" {
		t.Errorf("TableInfo.Title = %v, expected %v", info.Title, "DAMAGE TABLE")
	}

	if info.Emoji != "🗡️" {
		t.Errorf("TableInfo.Emoji = %v, expected %v", info.Emoji, "🗡️")
	}
}

func TestTableTypesMap(t *testing.T) {
	// We already tested the contents in table_handler_test.go
	// Here we just make sure the map exists and is accessible
	if tableTypes == nil {
		t.Error("tableTypes map should not be nil")
	}

	// Test that it contains the expected keys
	expectedKeys := []string{"damage", "healing"}
	
	for _, key := range expectedKeys {
		if _, exists := tableTypes[key]; !exists {
			t.Errorf("tableTypes should contain key '%s'", key)
		}
	}
	
	// Test a few specific values
	damageInfo := tableTypes["damage"]
	if damageInfo.Description != "damage done" {
		t.Errorf("damage Description = %v, expected %v", damageInfo.Description, "damage done")
	}

	healingInfo := tableTypes["healing"]
	if healingInfo.Description != "healing done" {
		t.Errorf("healing Description = %v, expected %v", healingInfo.Description, "healing done")
	}
}

// Test the internal constants if any exist in types.go
func TestInternalConstants(t *testing.T) {
	// Since types.go mostly defines types, there may not be many constants to test
	// But we can verify the structure exists and is accessible
	_ = TableInfo{}
	_ = tableTypes
	
	// If these don't cause compilation errors, the types are properly defined
}