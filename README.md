# wclogs-cli

A fast, terminal-based CLI tool that wraps the Warcraft Logs GraphQL API for power users who need scriptable access to combat data without browser overhead.

## 🎯 Project Status: Day 9 Complete! ✅

**Current Implementation**: Production-ready death analysis with Events API integration
**Working Commands**: 3 core commands (damage, healing, deaths) with advanced analysis
**Major Breakthrough**: Real ability names and detailed death timelines
**Events API Mastery**: 5-second damage timeline showing exact combat events

## 🚀 Features

### ✅ **Core Table Commands**
- **Damage/Healing Tables**: Professional display with DPS/HPS calculations
- **Player Filtering**: `--player "Name"` for focused analysis  
- **Export Options**: CSV and JSON export to `saved_reports/` directory
- **Smart Display**: Class colors, percentage calculations, top N filtering

### ✅ **Advanced Death Analysis** (Events API)
- **Production-Ready**: Real ability names ("Crystalline Shockwave from Fractillus")
- **Damage Timeline**: Shows exact damage sources in 5-second death window
- **Two-Mode System**: Summary for overview, detailed for specific player analysis
- **Friendly Fire Detection**: Reveals damage from other players vs boss
- **Healing Context**: Shows healing attempts with "healers tried hard!" insights
- **Smart Caching**: Ability name lookup with performance optimization

### ✅ **Technical Achievements**
- **Events API Integration**: Complex JSON parsing and event timeline analysis
- **GameData API**: Ability name resolution and actor lookup
- **GraphQL Mastery**: Multiple API endpoints, nested queries, error handling
- **Clean Architecture**: Services layer with lookup caching and shared handlers

## 📋 Requirements

- Go 1.19+
- Warcraft Logs API credentials (free at warcraftlogs.com)

## ⚡ Quick Start

```bash
# Install dependencies
go mod tidy

# Configure (interactive setup)
go run main.go config

# Get damage/healing tables
go run main.go damage 6qNJmgYBTcyfvpWF 3 --top 5
go run main.go healing 6qNJmgYBTcyfvpWF 3 --player "Hanahime"

# Advanced death analysis (NEW!)
go run main.go deaths 6qNJmgYBTcyfvpWF 3                    # Summary mode
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --player "Tekkyysp" # Detailed analysis

# Export data
go run main.go damage 6qNJmgYBTcyfvpWF 3 --output damage.csv
```

## 🎮 Real Working Examples

### Basic Table Analysis
```bash
# Top damage dealers
go run main.go damage 6qNJmgYBTcyfvpWF 3 --top 5

# Specific player performance
go run main.go healing 6qNJmgYBTcyfvpWF 3 --player "Hanahime"  # 1.67B healing!

# Export for spreadsheet analysis
go run main.go damage 6qNJmgYBTcyfvpWF 3 --output damage.csv
```

### Advanced Death Analysis (NEW!)
```bash
# Fight overview - who died and when
go run main.go deaths 6qNJmgYBTcyfvpWF 3

# Detailed player death investigation  
go run main.go deaths 6qNJmgYBTcyfvpWF 3 --player "Tekkyysp" --verbose

# Shows: "Killed by Crystalline Shockwave from Fractillus"
# Shows: 18.8M damage in 5 seconds with exact sources and amounts
# Shows: Healing attempts and defensive ability usage
```

## 🛠️ Commands

| Command | Status | Description | Key Features |
|---------|--------|-------------|--------------|
| `config` | ✅ | Set up API credentials | Interactive OAuth2 setup |
| `damage <report> <fight>` | ✅ | Damage table with DPS | Player filtering, top N, export |
| `healing <report> <fight>` | ✅ | Healing table with HPS | Player filtering, top N, export |
| `deaths <report> <fight>` | ✅ | **Advanced death analysis** | Real ability names, damage timeline, Events API |
| `interrupts <report> <fight>` | ❌ | Interrupt analysis | Coming in Day 11 |
| `players <report>` | ❌ | List all players | Coming in Day 11 |

**Legend**: ✅ Production Ready | ❌ Not Yet Implemented

### Deaths Command (Advanced Features)
- ✅ **Real ability names**: "Crystalline Shockwave from Fractillus" 
- ✅ **Damage timeline**: Shows 12-18M damage in 5-second death window
- ✅ **Two modes**: Summary (all deaths) vs Detailed (--player specific)
- ✅ **Friendly fire detection**: Reveals damage from other players
- ✅ **Healing analysis**: Shows healing attempts with context

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