# Warcraft Logs CLI

A fast, terminal-based CLI tool for analyzing Warcraft Logs data using GraphQL.

## Features

- **Damage/Healing Tables**: Professional display with DPS/HPS calculations
- **Advanced Death Analysis**: Real ability names and 5-second damage timelines using Events API
- **Player Filtering**: Filter results to specific players with case-insensitive matching
- **Export Options**: CSV and JSON export capabilities
- **Smart Caching**: Ability name lookup with performance optimization

## Quick Start

```bash
# Install dependencies
go mod tidy

# Configure API credentials (interactive setup)
go run main.go config

# Get damage table for a report
go run main.go damage <report-code> <fight-id>

# Get healing table
go run main.go healing <report-code> <fight-id>

# Advanced death analysis
go run main.go deaths <report-code> <fight-id>
```

## Documentation

For detailed usage examples and API information, see the [docs](./docs/) directory:

- [API Usage Examples](./docs/api_usage_examples.md) - Practical examples of all commands
- [GraphQL Queries](./docs/graphql_queries.md) - Technical query documentation
- [Configuration](./docs/configuration.md) - Setup and authentication details

## Commands

- `wclogs config` - Set up API credentials
- `wclogs damage [report] [fight]` - Show damage table
- `wclogs healing [report] [fight]` - Show healing table  
- `wclogs deaths [report] [fight]` - Advanced death analysis
- `wclogs help` - Show help for commands

## Requirements

- Go 1.19+
- Warcraft Logs API credentials (free at warcraftlogs.com)

## Configuration

Get API credentials at [Warcraft Logs API Clients](https://www.warcraftlogs.com/api/clients) and configure with:
```bash
go run main.go config
```