# Warcraft Logs CLI - API Usage Examples

This document contains detailed examples of how to use the Warcraft Logs CLI tool for various scenarios.

## Table of Contents
1. [Quick Start](#quick-start)
2. [Damage Analysis](#damage-analysis)
3. [Healing Analysis](#healing-analysis)
4. [Death Analysis](#death-analysis)
5. [Advanced Usage](#advanced-usage)
6. [Export Data](#export-data)
7. [Troubleshooting](#troubleshooting)

## Quick Start

### Setup
```bash
# Install dependencies
go mod tidy

# Configure API credentials (interactive setup)
go run main.go config
```

### Basic Commands
```bash
# Show damage table for report ABC123, fight 5
go run main.go damage ABC123 5

# Show healing table for report ABC123, fight 5
go run main.go healing ABC123 5

# Summary of all deaths in a fight
go run main.go deaths ABC123 5
```

## Damage Analysis

### Basic Usage
```bash
# Display damage done by all players in a fight
go run main.go damage 6qNJmgYBTcyfvpWF 3

# Show top 5 damage dealers
go run main.go damage 6qNJmgYBTcyfvpWF 3 --top 5

# Filter to a specific player
go run main.go damage 6qNJmgYBTcyfvpWF 3 --player "Pherally"
```

### Advanced Damage Analysis
```bash
# Compare multiple players' damage
go run main.go damage 6qNJmgYBTcyfvpWF 3 --player "Pherally" --verbose
go run main.go damage 6qNJmgYBTcyfvpWF 3 --player "Nikkans" --verbose

# Get detailed output with more information
go run main.go damage 6qNJmgYBTcyfvpWF 3 --verbose

# Save damage data to CSV for spreadsheet analysis
go run main.go damage 6qNJmgYBTcyfvpWF 3 --output damage.csv --top 10
```

## Healing Analysis

### Basic Usage
```bash
# Display healing done by all players in a fight
go run main.go healing 6qNJmgYBTcyfvpWF 3

# Show top 3 healers
go run main.go healing 6qNJmgYBTcyfvpWF 3 --top 3

# Filter to a specific healer
go run main.go healing 6qNJmgYBTcyfvpWF 3 --player "Hanahime"
```

### Advanced Healing Analysis
```bash
# Analyze a specific healer's performance
go run main.go healing 6qNJmgYBTcyfvpWF 3 --player "Truxpriest" --top 5

# Export healing data for analysis
go run main.go healing 6qNJmgYBTcyfvpWF 3 --output healing_report.json

# Verbose mode to see API call details
go run main.go healing 6qNJmgYBTcyfvpWF 3 --verbose
```

## Death Analysis

The death analysis is one of the most powerful features of this tool, using the Events API to provide detailed information about deaths during a fight.

### Summary Mode (Default)
```bash
# Get a summary of all deaths in a fight
go run main.go deaths 6qNJmgYBTcyfvpWF 3

# Verbose summary mode to see API progress
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --verbose
```

### Detailed Mode (Player-Specific)
```bash
# Get detailed analysis of a specific player's death(s)
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --player "Tekkyysp"

# Detailed analysis with verbose output
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --player "Tekkyysp" --verbose
```

### Death Analysis Features
The detailed death analysis provides:
- Exact survival time for each death
- What killed the player (with real ability names)
- Damage timeline in 5-second window before death
- Healing received during the critical period
- Defensive abilities used before death

## Advanced Usage

### Finding Report Codes
Report codes can be found in Warcraft Logs URLs:
- URL: `https://www.warcraftlogs.com/reports/Hw9TZc2WyrVKJLCa`
- Report Code: `Hw9TZc2WyrVKJLCa`

### Finding Fight IDs
Use damage or healing commands to identify available fights:
```bash
# This will show fight information including IDs
go run main.go damage Hw9TZc2WyrVKJLCa 1
go run main.go damage Hw9TZc2WyrVKJLCa 2
# etc.
```

### Player Name Matching
The tool supports case-insensitive player name matching:
```bash
# These all work for a player named "Pherally":
go run main.go damage ABC123 5 --player "Pherally"
go run main.go damage ABC123 5 --player "pheralLy"
go run main.go damage ABC123 5 --player "PHERALLY"
```

### Performance Tips
- Use `--top N` to speed up output for large raids
- Use `--output` to save data for later analysis instead of re-running queries
- The lookup service caches ability names automatically to reduce API calls

## Export Data

### CSV Export
```bash
# Export damage data as CSV
go run main.go damage 6qNJmgYBTcyfvpWF 3 --output damage.csv

# Export healing data as CSV  
go run main.go healing 6qNJmgYBTcyfvpWF 3 --output healing.csv

# Export top 10 damage dealers
go run main.go damage 6qNJmgYBTcyfvpWF 3 --top 10 --output top_damage.csv
```

### JSON Export
```bash
# Export detailed death analysis as JSON
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --output deaths.json

# Export specific player's death analysis
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --player "Tekkyysp" --output tekkyysp_analysis.json
```

All exported files are saved to the `saved_reports/` directory.

## Troubleshooting

### Common Issues and Solutions

#### Authentication Failed
```bash
# Reconfigure API credentials
go run main.go config
```

#### Report Not Found
```bash
# Check if the report code is correct
# Ensure the report is public (not private)
go run main.go damage ABC123 5
```

#### Player Not Found
```bash
# Verify the exact player name
# Use damage/healing commands first to see available players
go run main.go damage ABC123 5 --verbose
```

#### Fight ID Not Found
```bash
# Try different fight IDs
# Use damage command to see available fights in verbose mode
go run main.go damage ABC123 5 --verbose
```

### Debug Mode
Add `--verbose` to any command to see detailed API calls and processing steps:
```bash
go run main.go deaths ABC123 5 --verbose
go run main.go damage ABC123 5 --verbose
go run main.go healing ABC123 5 --player "PlayerName" --verbose
```

## Use Cases

### For Raid Leaders
```bash
# Overall performance overview
go run main.go damage 6qNJmgYBTcyfvpWF 3 --top 10
go run main.go healing 6qNJmgYBTcyfvpWF 3 --top 5

# Death investigation
go run main.go deaths 6qNJmgYBTcyfvpWF 3

# Individual performance review
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --player "ProblematicPlayer" --verbose
```

### For Individual Players
```bash
# Personal performance analysis
go run main.go damage 6qNJmgYBTcyfvpWF 3 --player "MyCharacter"
go run main.go healing 6qNJmgYBTcyfvpWF 3 --player "MyCharacter"

# Understanding deaths
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --player "MyCharacter" --verbose
```

### For Data Analysis
```bash
# Export data for external analysis
go run main.go damage 6qNJmgYBTcyfvpWF 3 --output raid_damage.csv
go run main.go healing 6qNJmgYBTcyfvpWF 3 --output raid_healing.csv
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --output death_analysis.json
```