package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"wclogs-cli/models"
)

// Query executes a GraphQL query
func (c *Client) Query(query string, variables map[string]any) (*models.GraphQLResponse, error) {
	// Ensure we have a valid auth token
	if err := c.authClient.EnsureValidToken(); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Prepare GraphQL request
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.authClient.GetAuthHeader())

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Parse response
	var gqlResp models.GraphQLResponse // ADD models. prefix
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for GraphQL errors
	if len(gqlResp.Errors) > 0 {
		return &gqlResp, fmt.Errorf("GraphQL error: %s", gqlResp.Errors[0].Message)
	}

	return &gqlResp, nil
}
