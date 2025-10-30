package auth

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test_id", "test_secret")

	if client.ClientID != "test_id" {
		t.Errorf("NewClient() ClientID = %v, expected %v", client.ClientID, "test_id")
	}

	if client.ClientSecret != "test_secret" {
		t.Errorf("NewClient() ClientSecret = %v, expected %v", client.ClientSecret, "test_secret")
	}

	if client.httpClient == nil {
		t.Error("NewClient() httpClient should not be nil")
	}

	// Check if the client has a default timeout
	if client.httpClient.Timeout == 0 {
		t.Error("NewClient() httpClient should have a timeout")
	}
}

func TestIsTokenValid(t *testing.T) {
	tests := []struct {
		name          string
		accessToken   string
		expiresAt     time.Time
		expectedValid bool
	}{
		{
			name:          "valid token",
			accessToken:   "valid_token",
			expiresAt:     time.Now().Add(1 * time.Hour), // Expires in 1 hour
			expectedValid: true,
		},
		{
			name:          "empty token",
			accessToken:   "",
			expiresAt:     time.Now().Add(1 * time.Hour),
			expectedValid: false,
		},
		{
			name:          "expired token",
			accessToken:   "valid_token",
			expiresAt:     time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
			expectedValid: false,
		},
		{
			name:          "empty token and expired",
			accessToken:   "",
			expiresAt:     time.Now().Add(-1 * time.Hour),
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				AccessToken: tt.accessToken,
				ExpiresAt:   tt.expiresAt,
			}

			result := client.IsTokenValid()
			if result != tt.expectedValid {
				t.Errorf("IsTokenValid() = %v, expected %v", result, tt.expectedValid)
			}
		})
	}
}

func TestGetAuthHeader(t *testing.T) {
	client := &Client{
		AccessToken: "test_token_123",
	}

	expected := "Bearer test_token_123"
	result := client.GetAuthHeader()

	if result != expected {
		t.Errorf("GetAuthHeader() = %v, expected %v", result, expected)
	}
}

func TestEnsureValidToken(t *testing.T) {
	// This test would require mocking the GetAccessToken method
	// For now, we'll test the logic path when token is valid

	client := &Client{
		AccessToken: "valid_token",
		ExpiresAt:   time.Now().Add(1 * time.Hour), // Token is valid
	}

	// Should not need to get a new token since it's still valid
	err := client.EnsureValidToken()
	if err != nil {
		t.Errorf("EnsureValidToken() returned error when token was valid: %v", err)
	}

	if client.AccessToken != "valid_token" {
		t.Errorf("EnsureValidToken() changed token when it should not have")
	}
}

func TestTokenResponse(t *testing.T) {
	// Just verify the TokenResponse structure works as expected
	tokenResp := TokenResponse{
		AccessToken: "test_token",
		TokenType:   "Bearer",
		ExpiresIn:   3600, // 1 hour
	}

	if tokenResp.AccessToken != "test_token" {
		t.Errorf("TokenResponse AccessToken = %v, expected %v", tokenResp.AccessToken, "test_token")
	}

	if tokenResp.TokenType != "Bearer" {
		t.Errorf("TokenResponse TokenType = %v, expected %v", tokenResp.TokenType, "Bearer")
	}

	if tokenResp.ExpiresIn != 3600 {
		t.Errorf("TokenResponse ExpiresIn = %v, expected %v", tokenResp.ExpiresIn, 3600)
	}
}
