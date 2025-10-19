package auth

import (
	"net/http"
	"time"
)

// TokenResponse represents the OAuth2 token response from Warcraft Logs
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Client handles OAuth2 authentication with Warcraft Logs
type Client struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	ExpiresAt    time.Time
	httpClient   *http.Client
}

// NewClient creates a new auth client with the given credentials
func NewClient(clientID, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}
