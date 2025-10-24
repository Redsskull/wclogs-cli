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
}

// can add more types as I go along
