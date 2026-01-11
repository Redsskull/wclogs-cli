# wclogs-cli Commands Reference

Complete reference for all available commands in the Warcraft Logs CLI tool.

## 🚀 Quick Installation

```bash
# Native (recommended)
go install github.com/Redsskull/wclogs-cli@latest

# Container (no Go required)
git clone https://github.com/Redsskull/wclogs-cli.git && cd wclogs-cli
make container-run ARGS="config"

**NEW**: Fight ID now supports "last" keyword to automatically select the most recent fight!

**🔍 WCL API Discovery**: We discovered that Warcraft Logs uses `encounterID = 0` to mark trash fights. Our "last" implementation now uses this official WCL logic to match web interface behavior perfectly!
```

For comprehensive documentation, see the [docs](./docs/) directory:
- [API Usage Examples](./docs/api_usage_examples.md) - Usage examples and scenarios
- [GraphQL Queries](./docs/graphql_queries.md) - Technical query documentation
- [Configuration](./docs/configuration.md) - Setup and authentication details

---

## 📋 Command Overview

| Command | Status | Description |
|---------|--------|-------------|
| `config` | ✅ Working | Set up API credentials |
| `damage` | ✅ Working | Show damage tables with player filtering |
| `healing` | ✅ Working | Show healing tables with player filtering |
| `deaths` | ✅ Working | Advanced death analysis with Events API |
| `interrupts` | ✅ Working | Professional interrupt analysis with spell correlation |
| `players` | ✅ Working | List all players with class/role filtering and search |
| `help` | ✅ Working | Show help for commands |
| `completion` | ✅ Working | Generate shell completions |

---

## 🔧 Setup Commands

### `wclogs config`
**Purpose**: Interactive setup of Warcraft Logs API credentials

**Usage**:
```bash
# Native binary
./wclogs config

# Container
make container-run ARGS="config"
```

**What it does**:
- Prompts for Client ID and Client Secret
- Saves credentials to `~/.wclogs.yaml`
- Tests authentication with the API

**Example**:
```bash
$ ./wclogs config
🔧 WARCRAFT LOGS API SETUP

Enter your Client ID: your_client_id_here
Enter your Client Secret: your_client_secret_here

✅ Credentials saved to /home/user/.wclogs.yaml
✅ Authentication test successful!
```

**Prerequisites**: You need API credentials from https://www.warcraftlogs.com/api/clients/

---

## 📊 Table Commands

### `wclogs damage [report-code] [fight-id|last]`
**Purpose**: Display damage done by all players in a fight

**Usage**:
```bash
# Native binary
./wclogs damage <report-code> <fight-id|last> [flags]

# Container
make container-run ARGS="damage <report-code> <fight-id|last> [flags]"
```

**Flags**:
- `--top N` - Show only top N players (default: all)
- `--player "Name"` - Show only specific player
- `--output file.csv` - Save to file (CSV/JSON supported)
- `--no-color` - Disable colored output
- `--verbose` - Show detailed progress

### `wclogs healing [report-code] [fight-id|last]`
**Purpose**: Display healing done by all players in a fight

**Usage**:
```bash
# Native binary
./wclogs healing <report-code> <fight-id|last> [flags]

# Container
make container-run ARGS="healing <report-code> <fight-id|last> [flags]"
```

**Flags**: Same as damage command

---

## 💀 Advanced Analysis Commands

### `wclogs deaths [report-code] [fight-id|last]`
**Purpose**: Advanced death analysis using Events API with real ability names

**Two Modes**:
1. **Summary Mode** (default): Overview of all deaths
2. **Detailed Mode** (`--player` flag): Deep analysis for specific player

**Usage**:
```bash
# Native binary
./wclogs deaths <report-code> <fight-id|last> [flags]

# Container
make container-run ARGS="deaths <report-code> <fight-id|last> [flags]"
```

**Flags**:
- `--player "Name"` - Detailed analysis for specific player
- `--verbose` - Show debug information and API progress
- `--output file.json` - Save analysis to file

**Key Features**:
- Real ability names: Shows "Crystalline Shockwave from Fractillus" not "Ability ID 1226823"
- Damage timeline: Exact damage amounts and sources in 5-second death window
- Friendly fire detection: Shows damage from other players
- Healing context: Shows healing attempts with contextual insights
- Survival analysis: Calculates correct survival times from fight start

---

### `wclogs interrupts [report-code] [fight-id|last]`
**Purpose**: Professional-grade interrupt analysis with spell correlation and missed opportunity tracking

**Two Modes**:
1. **Summary Mode** (default): Overview of all interrupts and spell analysis
2. **Player Mode** (`--player` flag): Detailed analysis for specific player

**Usage**:
```bash
# Native binary
./wclogs interrupts <report-code> <fight-id|last> [flags]

# Container
make container-run ARGS="interrupts <report-code> <fight-id|last> [flags]"
```

**Flags**:
- `--player "Name"` - Detailed analysis for specific player
- `--verbose` - Show debug information, API progress, and cache statistics
- `--output file.json` - Save analysis to file

**Key Features**:
- **Professional WCL-style analysis**: Matches web interface functionality
- **Spell correlation**: Uses `extraAbilityGameID` for perfect interrupt-to-spell mapping
- **Success rate tracking**: Shows stopped vs missed opportunities breakdown
- **Real ability names**: Displays actual spell names, not IDs
- **Player-specific analysis**: Individual interrupt performance with timing data
- **Spell priority analysis**: Shows which spells were interrupted most/least
- **Cache optimization**: Smart ability and actor name lookups

**Examples**:
```bash
# Native binary
./wclogs interrupts YMRqjzC2WPnhwNJd 2                      # Full interrupt overview
./wclogs interrupts YMRqjzC2WPnhwNJd last                   # Last fight interrupt analysis
./wclogs interrupts YMRqjzC2WPnhwNJd 2 --player "BlagZeras"  # Player-specific detailed analysis
./wclogs interrupts YMRqjzC2WPnhwNJd 2 --verbose            # Verbose mode with API progress

# Container
make container-run ARGS="interrupts YMRqjzC2WPnhwNJd 2"
make container-run ARGS="interrupts YMRqjzC2WPnhwNJd last"
make container-run ARGS="interrupts YMRqjzC2WPnhwNJd 2 --player BlagZeras"
make container-run ARGS="interrupts YMRqjzC2WPnhwNJd 2 --verbose"
```

### `wclogs players [report-code] [fight-id]`
**Purpose**: List all players in a report or specific fight with advanced filtering capabilities

**Usage**:
```bash
# Native binary
./wclogs players <report-code> [fight-id] [flags]

# Container
make container-run ARGS="players <report-code> [fight-id] [flags]"
```

**Flags**:
- `--class "ClassName"` - Filter by specific class (e.g., Paladin, Warrior)
- `--role "RoleName"` - Filter by role (Tank, Healer, DPS)
- `--search "PlayerName"` - Search by player name (partial match)
- `--debug` - Show detailed spec icon debugging information
- `--output file.csv` - Save to file (CSV/JSON supported)
- `--top N` - Show only top N players
- `--no-color` - Disable colored output
- `--verbose` - Show detailed progress

**Key Features**:
- **Fight-Specific Filtering**: Optional fight ID parameter to show only fight participants
- **Automatic Role Detection**: Smart role assignment based on class and spec
- **Advanced Filtering**: Filter by class, role, or search by name
- **Class Color Coding**: Each WoW class displayed in appropriate colors
- **Composition Summary**: Shows role and class distribution
- **Export Support**: Save player lists to CSV or JSON
- **Server Information**: Displays player server/realm
- **Smart Actor Filtering**: Automatically filters out NPCs, pets, and incomplete data

**Examples**:
```bash
# Native binary
./wclogs players 6qNJmgYBTcyfvpWF                       # List all players in report
./wclogs players 6qNJmgYBTcyfvpWF 5                     # List players in fight 5
./wclogs players 6qNJmgYBTcyfvpWF last                  # List players in last fight
./wclogs players 6qNJmgYBTcyfvpWF --class "Paladin"     # Filter by class
./wclogs players 6qNJmgYBTcyfvpWF --role "Tank"         # Filter by role
./wclogs players 6qNJmgYBTcyfvpWF --search "Pmpm"       # Search by name
./wclogs players 6qNJmgYBTcyfvpWF --output players.csv  # Export to file
./wclogs players 6qNJmgYBTcyfvpWF 5 --role "DPS" --top 10 # Top 10 DPS in fight 5

# Container
make container-run ARGS="players 6qNJmgYBTcyfvpWF"
make container-run ARGS="players 6qNJmgYBTcyfvpWF 5"
make container-run ARGS="players 6qNJmgYBTcyfvpWF last"
make container-run ARGS="players 6qNJmgYBTcyfvpWF --class Paladin"
make container-run ARGS="players 6qNJmgYBTcyfvpWF --role Tank"
```

**Role Detection Logic**:
The command automatically detects player roles using:
- **Pure Classes**: Hunter, Mage, Rogue, Warlock → DPS
- **Tank Classes**: Death Knight (Blood), Demon Hunter (Vengeance), Warrior (Protection)
- **Hybrid Detection**: Uses spec icons when available for Druid, Monk, Paladin, Priest, Shaman, Evoker

**Use Cases**:
- **Fight Analysis**: See who participated in specific encounters
- **Raid Planning**: Check composition balance before pulls
- **Player Lookup**: Find exact player names for other commands
- **Guild Management**: Export member lists for analysis
- **Role Distribution**: Verify tank/healer/DPS ratios
- **Clean Data**: Automatically filters out pets, NPCs, and incomplete entries

---

## 🌐 Global Flags

All commands support these global flags:

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Save to file (CSV/JSON) |
| `--top` | `-t` | Show top N players |
| `--verbose` | `-v` | Enable verbose output |
| `--help` | `-h` | Show command help |

---

## 🎯 Fight ID Options

All table and analysis commands now support two fight ID formats:
- **Numeric ID**: `5`, `12`, `99` - Specific fight number
- **"last" keyword**: `last` - Last meaningful boss encounter (matches WCL web interface)

### 🔬 How "Last" Fight Resolution Works

Our implementation uses **official Warcraft Logs API logic** discovered through documentation research:

1. **Queries all fights** in the report
2. **Filters out trash fights** using WCL's official rule: `encounterID = 0`
3. **Selects the last remaining fight** (highest ID among boss encounters)

This **perfectly matches** the WCL web interface `fight=last` behavior!

**Discovery Source**: [WCL GraphQL API Documentation](https://www.warcraftlogs.com/v2-api-docs/warcraft/reportfight.doc.html)
> "If the encounterID is 0, the fight is considered a trash fight."

**Examples**:
```bash
./wclogs damage ABC123 5      # Specific fight
./wclogs damage ABC123 last   # Last meaningful boss encounter
./wclogs healing XYZ789 last  # Uses same logic across all commands

# Verbose mode shows the resolution process:
./wclogs damage ABC123 last --verbose
# ⏭️  Skipping fight #38 (Trash) - trash fight (encounterID = 0)
# ✅ Found last meaningful fight #36: Boss Name (KILL, encounterID: 1234)
```

## 🎯 File Output Formats

**Output Location**: All files saved to `saved_reports/` directory

---

## 🔧 Troubleshooting

### Common Errors

**"Authentication failed"**
```bash
# Native binary
./wclogs config

# Container
make container-run ARGS="config"
```

**"Report 'ABC123' not found"**
- Check the report code is correct
- Ensure the report is public (not private)

**"Fight 99 not found"**  
- Check available fights with damage/healing commands first
- Fight IDs start from 1

**"failed to resolve fight ID 'last'"**
- Report may have no boss encounters (only trash fights)
- All fights may have `encounterID = 0`
- Use `--verbose` to see the resolution process

**"Player 'Name' not found"**
- Use exact player name (case-sensitive)
- Check spelling and special characters

### Understanding "Last" Fight Resolution

If `fight=last` gives unexpected results:

1. **Use verbose mode** to see what's happening:
   ```bash
   ./wclogs damage ABC123 last --verbose
   ```

2. **Check fight types**: Our tool skips trash fights (encounterID = 0)
   - This matches WCL web interface behavior
   - Sometimes the "last" meaningful fight isn't chronologically last

3. **Compare with web interface**: 
   - Open `https://www.warcraftlogs.com/reports/ABC123?fight=last&type=damage-done`
   - Our tool should select the same fight the web interface shows

### Debug Mode
Add `--verbose` to any command for detailed debugging:
```bash
# Native binary
./wclogs deaths ABC123 5 --verbose
./wclogs deaths ABC123 last --verbose  # Also works with "last"

# Container
make container-run ARGS="deaths ABC123 5 --verbose"
make container-run ARGS="deaths ABC123 last --verbose"
```
Shows API calls, response sizes, processing steps, and fight resolution when using "last".

---

## ❌ Commands Not Yet Implemented

| Command | Status | Planned |
|---------|--------|---------|
| `timeline` | ❌ Not implemented | Future |
| `boss-abilities` | ❌ Not implemented | Future |

## 📚 API Research Insights

During development, we made several key discoveries about the Warcraft Logs API:

### 🎯 "Last" Fight Logic Discovery
- **Challenge**: WCL web interface `fight=last` didn't match naive "chronologically last" approach
- **Research**: Found official WCL GraphQL API documentation
- **Discovery**: `encounterID = 0` officially marks trash fights
- **Solution**: Filter out trash fights before selecting "last"
- **Result**: Perfect match with web interface behavior

### 🔍 GraphQL API Insights
- **fightIDs parameter**: Only accepts `[Int]`, not strings like `"last"`
- **Client-side resolution**: Web interface does same resolution we implemented
- **encounterID field**: Key to distinguishing boss fights from trash
- **Official documentation**: [WCL API Docs](https://www.warcraftlogs.com/v2-api-docs/)

### 💡 Key Learnings
- Always research official API documentation
- Web interface behavior often reveals underlying API logic
- GraphQL schemas contain crucial implementation details
- Client-side processing sometimes necessary even with robust APIs

For complete usage examples and detailed command information, see the [API Usage Examples](./docs/api_usage_examples.md) in the docs directory.