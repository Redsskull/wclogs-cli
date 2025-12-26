package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetAccessToken performs the OAuth2 client credentials flow
func (c *Client) GetAccessToken() error {
	// Create Basic Auth header
	credentials := c.ClientID + ":" + c.ClientSecret
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	// Create HTTP request
	req, err := http.NewRequest("POST",
		"https://www.warcraftlogs.com/oauth/token",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Authorization", "Basic "+encoded)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// ake the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed with status %d", resp.StatusCode)
	}

	// Parse the JSON response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Store the token and expiration
	c.AccessToken = tokenResp.AccessToken
	c.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}
