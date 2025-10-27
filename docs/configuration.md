# Configuration and Authentication

This document explains how the Warcraft Logs CLI handles configuration and authentication.

## Table of Contents
1. [API Credentials Setup](#api-credentials-setup)
2. [Configuration File](#configuration-file)
3. [OAuth2 Authentication](#oauth2-authentication)
4. [Token Management](#token-management)
5. [Troubleshooting](#troubleshooting)

## API Credentials Setup

### Getting Warcraft Logs API Credentials

1. Visit [Warcraft Logs API Clients](https://www.warcraftlogs.com/api/clients)
2. Click "Create a new client"
3. Fill in the client details:
   - **Name**: wclogs-cli (or your preferred name)
   - **Description**: CLI tool for Warcraft Logs analysis
   - **Website**: (optional)
4. Copy the Client ID and Client Secret

### Interactive Setup

The simplest way to configure the tool is using the interactive setup:

```bash
go run main.go config
```

This will:
- Prompt for your Client ID
- Prompt for your Client Secret
- Test the credentials with the API
- Save them to `~/.wclogs.yaml`

### Manual Configuration

You can also create the configuration file manually:

Create `~/.wclogs.yaml` with the following content:

```yaml
client_id: your_client_id_here
client_secret: your_client_secret_here
```

## Configuration File

### File Location
The configuration file is stored at `~/.wclogs.yaml` in your home directory.

### Format
```yaml
client_id: "your_client id string"
client_secret: "your client secret string"
```

### Security
The configuration file is created with read/write permissions only for the owner (0600).

### Validation
The tool validates that both `client_id` and `client_secret` are present and non-empty.

## OAuth2 Authentication

### Authentication Flow

The tool uses the OAuth2 Client Credentials flow:

1. **Token Request**: Exchanges Client ID and Client Secret for an access token
2. **Token Usage**: Uses the access token for all subsequent API requests
3. **Token Refresh**: Automatically refreshes the token when it expires

### Implementation

The authentication is handled by the `auth.Client` type:

```go
type Client struct {
    ClientID     string
    ClientSecret string
    AccessToken  string
    ExpiresAt    time.Time
    httpClient   *http.Client
}
```

### Authentication Process

1. **Credentials Encoding**: Client ID and Client Secret are combined and Base64 encoded
2. **Token Request**: Sends POST request to `https://www.warcraftlogs.com/oauth/token`
3. **Authorization Header**: Uses Basic auth with the encoded credentials
4. **Response Handling**: Parses the JSON response to extract the access token

### API Request Authorization

All API requests include the Bearer token in the Authorization header:

```
Authorization: Bearer your_access_token_here
```

## Token Management

### Token Validation

The tool checks if the current token is still valid before each API request:

```go
func (c *Client) IsTokenValid() bool {
    return c.AccessToken != "" && time.Now().Before(c.ExpiresAt)
}
```

### Automatic Token Refresh

Before making any API request, the tool ensures a valid token:

```go
func (c *Client) EnsureValidToken() error {
    if !c.IsTokenValid() {
        return c.GetAccessToken()
    }
    return nil
}
```

### Token Expiration

- Tokens have a standard expiration time provided by the API (typically 1 hour)
- The `ExpiresAt` field stores the exact expiration time
- The token is automatically refreshed if it's expired or will expire soon

## Configuration Commands

### Setup Command
```bash
# Interactive configuration setup
go run main.go config
```

### Verification
The config command:
- Prompts for Client ID and Client Secret
- Tests the credentials by making a sample API call
- Saves valid credentials to the config file
- Shows success or error message

## Troubleshooting

### Common Issues

#### Invalid Credentials
```
Error: authentication failed with status 401
```

**Solution**: Check your Client ID and Client Secret. Run `go run main.go config` again.

#### Configuration File Not Found
```
Error: config file not found at /home/user/.wclogs.yaml
```

**Solution**: Run `go run main.go config` to set up your credentials.

#### Token Expiration
If you see authentication errors during long-running operations:
- The tool should automatically refresh the token
- If not, the issue is likely in the token refresh logic

### Debugging Authentication

Use the `--verbose` flag to see authentication details:

```bash
go run main.go damage ABC123 5 --verbose
```

This will show:
- Configuration loading
- Authentication setup
- Token refresh if needed
- API request details

### Testing Configuration

You can test your configuration by running any command:

```bash
# This will test authentication and show if there are any issues
go run main.go damage ABC123 5
```

### Environment Variables (Future Enhancement)

While not currently implemented, you could consider supporting environment variables:

```bash
export WCLOGS_CLIENT_ID="your_client_id"
export WCLOGS_CLIENT_SECRET="your_client_secret"
```

## Security Best Practices

1. **Secure Storage**: Configuration is stored with limited permissions (0600)
2. **No Echo**: Client Secret is not echoed to the terminal during input
3. **Token Security**: Access tokens are stored in memory only and not persisted
4. **HTTPS**: All API communication uses HTTPS
5. **Short-Lived Tokens**: Tokens automatically expire and are refreshed

## Configuration Validation

The application validates configuration at startup:
- Checks that the config file exists
- Verifies that required fields are present
- Tests that credentials work with the API

This validation prevents most configuration-related errors before attempting API operations.