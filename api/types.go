package api

import (
	"net/http"

	"wclogs-cli/auth"
)

// GraphQLRequest represents a GraphQL query request
type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL API response
type GraphQLResponse struct {
	Data   any            `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLErrorLocation `json:"locations,omitempty"`
	Path      []any                  `json:"path,omitempty"`
}

type GraphQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
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
