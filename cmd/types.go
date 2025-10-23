package cmd

import (
	"wclogs-cli/api"
)

// TableInfo contains display information for different data types
type TableInfo struct {
	Title       string
	Emoji       string
	DataType    api.DataType
	Description string
}

// tableTypes defines all supported table types and their display info
var tableTypes = map[string]TableInfo{
	"damage": {
		Title:       "DAMAGE TABLE",
		Emoji:       "🗡️",
		DataType:    api.DataTypeDamage,
		Description: "damage done",
	},
	"healing": {
		Title:       "HEALING TABLE",
		Emoji:       "💚",
		DataType:    api.DataTypeHealing,
		Description: "healing done",
	},
	"deaths": {
		Title:       "DEATHS TABLE",
		Emoji:       "💀",
		DataType:    api.DataTypeDeaths,
		Description: "death events",
	},
	"interrupts": {
		Title:       "INTERRUPTS TABLE",
		Emoji:       "🛑",
		DataType:    api.DataTypeInterrupts,
		Description: "interrupts performed",
	},
}

// can add more types as I go along
