# wclogs-cli

A fast, terminal-based CLI tool that wraps the Warcraft Logs GraphQL API for power users who need scriptable access to combat data without browser overhead.

## 🎯 Project Status: Days 6 & 7 Complete! ✅

**Current Implementation**: Full player analysis foundation with masterData integration
**Working Commands**: 5 total (damage, healing, deaths, interrupts, players)
**Real Data Tested**: ✅ Korean raid reports with 375+ players
**Player Filtering**: ✅ `--player "Name"` works on all commands

## 🚀 Features

✅ **Multiple Data Types**: Damage, healing, deaths, interrupts  
✅ **Player Analysis**: List all players + filter any command by player name  
✅ **Smart Display**: Proper column headers (DPS vs HPS), class colors, empty data handling  
✅ **Export Options**: CSV and JSON export to `saved_reports/` directory  
✅ **Professional UX**: Color coding, helpful errors, player name validation  
✅ **Clean Architecture**: Zero file explosion - shared handlers for all commands  

## 📋 Requirements

- Go 1.19+
- Warcraft Logs API credentials (free at warcraftlogs.com)

## ⚡ Quick Start

```bash
# Install dependencies
go mod tidy

# Configure (interactive setup)
go run main.go config

# Show all players in a report
go run main.go players Hw9TZc2WyrVKJLCa

# Get damage data for a fight
go run main.go damage Hw9TZc2WyrVKJLCa 99

# Filter to specific player's healing
go run main.go healing Hw9TZc2WyrVKJLCa 99 --player "Sketch"

# Export damage data to CSV
go run main.go damage Hw9TZc2WyrVKJLCa 99 --output damage.csv
```

## 🎮 Real Working Examples

```bash
# List all 375 players with classes and servers
go run main.go players Hw9TZc2WyrVKJLCa

# Show top 5 damage dealers
go run main.go damage Hw9TZc2WyrVKJLCa 99 --top 5

# Focus on specific player's performance
go run main.go damage Hw9TZc2WyrVKJLCa 99 --player "Pmpm"    # 3.1B damage!
go run main.go healing Hw9TZc2WyrVKJLCa 99 --player "Sketch" # 3.4B healing!

# Export for spreadsheet analysis
go run main.go healing Hw9TZc2WyrVKJLCa 99 --output healers.csv
go run main.go players Hw9TZc2WyrVKJLCa --output players.json

# Debug mode for API troubleshooting
go run main.go damage Hw9TZc2WyrVKJLCa 99 --verbose
```

## 🛠️ Commands

| Command | Status | Description | Options |
|---------|--------|-------------|---------|
| `config` | ✅ | Set up API credentials | Interactive setup |
| `players <report>` | ✅ | List all players in report | `--output FILE` |
| `damage <report> <fight>` | ✅ | Damage done table | `--top N`, `--player NAME` |
| `healing <report> <fight>` | ✅ | Healing done table | `--top N`, `--player NAME` |
| `deaths <report> <fight>` | ⚠️ | Death events (basic) | `--player NAME` |
| `interrupts <report> <fight>` | ⚠️ | Interrupt data (basic) | `--player NAME` |

**Legend**: ✅ Fully Working | ⚠️ Basic Implementation | ❌ Not Implemented

**Global Flags**:
- `--top N`: Show only top N players
- `--player "Name"`: Filter to specific player (case-insensitive)
- `--output FILE`: Export to CSV/JSON in `saved_reports/`
- `--no-color`: Disable colored output
- `--verbose`: Show detailed API progress

## 🎯 Key Achievements

### ✅ **Player Analysis Foundation**
```bash
# masterData GraphQL integration - gets all players with classes/servers
go run main.go players Hw9TZc2WyrVKJLCa
# Returns: 375 players with Name, Class, Server, colored by role

# Player filtering with validation and suggestions  
go run main.go damage Hw9TZc2WyrVKJLCa 99 --player "Pmpm"
# ✅ Player 'Pmpm' found in report
# 🎯 Filtered to 1 player(s) matching 'Pmpm'
# Shows: 3,124,207,218 damage (100.0% of filtered view)
```

### ✅ **Smart Display System**
- **Adaptive Headers**: Damage/DPS vs Healing/HPS columns automatically
- **Class Colors**: Evokers/Shamans show as healers (green), DPS classes (red), etc.
- **Empty Data Handling**: "No deaths found (great job!)" instead of confusing empty tables
- **Percentage Calculations**: Proper % calculations, no more NaN% bugs

### ✅ **Professional Architecture**  
- **Zero File Explosion**: All commands use shared `executeTableCommand()` handler
- **Generic Display**: `DisplayTable()` adapts to damage/healing/deaths/interrupts  
- **Clean Error Messages**: Player name suggestions when typos occur
- **Modern Go**: Using `max()` function, clean struct definitions

## ⚠️ **Known Limitations**

### **Deaths & Interrupts Need Events API**
```bash
# These work but show limited data:
go run main.go deaths Hw9TZc2WyrVKJLCa 99    # Shows "No deaths found" 
go run main.go interrupts Hw9TZc2WyrVKJLCa 99 # Shows "No interrupts found"

# Reason: Table API doesn't provide detailed event data
# Solution: Needs Events API integration (planned for Week 2)
```

**What's Missing:**
- **Deaths**: Should show what killed player, damage taken timeline
- **Interrupts**: Should show successful vs missed, interrupt targets
- **Root Cause**: Using `table` dataType instead of `events` API

### **Current Workarounds:**
- Deaths/Interrupts have basic table structure but may show empty data
- Use damage/healing commands for reliable analysis
- Player filtering works on all commands

## 🚀 **Future Work (Week 2)**

```bash
# Planned improvements:
wclogs death-analysis ABC123 5 --player "Name"  # What killed them + timeline
wclogs interrupts ABC123 5                      # Success rate, missed opportunities  
wclogs player-damage ABC123 5 "Name"            # Ability breakdown
wclogs timeline ABC123 5                        # Fight event timeline
```

## ⚙️ **Configuration**

Interactive setup (recommended):
```bash
go run main.go config
# Prompts for API credentials from https://www.warcraftlogs.com/api/clients
```

Manual setup - Create `~/.wclogs.yaml`:
```yaml
client_id: your_warcraft_logs_client_id
client_secret: your_warcraft_logs_client_secret
```

## 🎯 **Proven Capabilities**

This project demonstrates:

✅ **GraphQL Mastery**: Complex nested queries, masterData + table integration  
✅ **OAuth2 Implementation**: Client credentials flow with token management  
✅ **API Integration**: Real-world data from 375+ player Korean raids  
✅ **CLI Design**: Cobra framework, shared handlers, zero code duplication  
✅ **Data Processing**: Player name→ID mapping, filtering, validation  
✅ **User Experience**: Color coding, helpful errors, professional formatting  

## 📊 **Real Performance Data**

**Tested With:**
- **Report**: Hw9TZc2WyrVKJLCa (Korean server raid)
- **Players**: 375 total participants  
- **Fight 99 Highlights**:
  - Pmpm (Mage): 3.1B damage, 7.4M DPS
  - Sketch (Evoker): 3.4B healing, 6.3M HPS
  - 19 healers, 15.5B total healing

**Architecture Handles:**
- Large player datasets (375+ players)
- Korean character encoding  
- Multiple data types simultaneously
- Case-insensitive player matching

---

*Days 6 & 7 Complete! ✅ Player analysis foundation fully implemented.*  
*Next: Events API integration for detailed death/interrupt analysis.*