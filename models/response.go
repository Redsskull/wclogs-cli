package models

import "encoding/json"

// GraphQLResponse represents the top-level GraphQL API response
// All GraphQL responses follow this pattern: data + errors
type GraphQLResponse struct {
	Data   *ResponseData  `json:"data,omitempty"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// GraphQLError represents a single GraphQL error
type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLErrorLocation `json:"locations,omitempty"`
	Path      any                    `json:"path,omitempty"`
}

// GraphQLErrorLocation represents where in the query an error occurred
type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// ResponseData represents the "data" field in GraphQL responses
// This is where the actual query results live
type ResponseData struct {
	ReportData *ReportData `json:"reportData,omitempty"`
}

// ReportData represents the reportData field in the API
type ReportData struct {
	Report *Report `json:"report,omitempty"`
}

// Report represents a single Warcraft Logs report
type Report struct {
	Code      string          `json:"code,omitempty"`      // Report code like "ABC123"
	Title     string          `json:"title,omitempty"`     // Report title
	StartTime int64           `json:"startTime,omitempty"` // Unix timestamp
	EndTime   int64           `json:"endTime,omitempty"`   // Unix timestamp
	Fights    []Fight         `json:"fights,omitempty"`    // All fights in this report
	Table     json.RawMessage `json:"table,omitempty"`     // Table data for this report
}

// Fight represents a single encounter/fight within a report
type Fight struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`            // Boss name
	EncounterID     int     `json:"encounterID"`     // Encounter ID
	StartTime       int64   `json:"startTime"`       // Fight start (relative to report start)
	EndTime         int64   `json:"endTime"`         // Fight end (relative to report start)
	Kill            bool    `json:"kill"`            // true if boss was killed
	Difficulty      int     `json:"difficulty"`      // Difficulty (10N, 25H, etc)
	FightPercentage float64 `json:"fightPercentage"` // Boss health % when fight ended
}
