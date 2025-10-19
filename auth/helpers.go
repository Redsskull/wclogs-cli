package auth

import "time"

// IsTokenValid checks if the current access token is still valid
func (c *Client) IsTokenValid() bool {
	return c.AccessToken != "" && time.Now().Before(c.ExpiresAt)
}

// GetAuthHeader returns the Authorization header value for API requests
func (c *Client) GetAuthHeader() string {
	return "Bearer " + c.AccessToken
}

// EnsureValidToken gets a new token if the current one is expired
func (c *Client) EnsureValidToken() error {
	if !c.IsTokenValid() {
		return c.GetAccessToken()
	}
	return nil
}
