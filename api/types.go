package api

import (
	"fmt"
	"net/http"

	"wclogs-cli/auth"
)

// GraphQLRequest represents a GraphQL query request
type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// QueryVariables represents the variables we pass to GraphQL queries
type QueryVariables struct {
	Code    string `json:"code"`    // Report code like "ABC123"
	FightID int    `json:"fightID"` // Fight ID like 5
}

// Client handles GraphQL API requests to Warcraft Logs
type Client struct {
	authClient *auth.Client
	httpClient *http.Client
	endpoint   string
}

// NewClient creates a new GraphQL API client
func NewClient(authClient *auth.Client) *Client {
	return &Client{
		authClient: authClient,
		httpClient: &http.Client{},
		endpoint:   "https://www.warcraftlogs.com/api/v2/client",
	}
}

// ValidateQueryVariables checks if the query variables are valid
func ValidateQueryVariables(code string, fightID int) error {
	if code == "" {
		return fmt.Errorf("report code cannot be empty")
	}

	if len(code) < 6 {
		return fmt.Errorf("report code '%s' is too short (must be at least 6 characters)", code)
	}

	if fightID <= 0 {
		return fmt.Errorf("fight ID must be greater than 0, got: %d", fightID)
	}

	return nil
}

// DataType represents the different types of combat data we can fetch
type DataType string

const (
	DataTypeDamage     DataType = "DamageDone"
	DataTypeHealing    DataType = "Healing"
	DataTypeDeaths     DataType = "Deaths"
	DataTypeInterrupts DataType = "Interrupts"
)

// EventHostilityType represents the hostility type for filtering events
type EventHostilityType string

const (
	EventHostilityFriendly EventHostilityType = "Friendlies" // Fixed!
	EventHostilityHostile  EventHostilityType = "Enemies"
	EventHostilityAll      EventHostilityType = "All"
)
