# wclogs-cli Commands Reference

Complete reference for all available commands in the Warcraft Logs CLI tool.

## 🚀 Quick Installation

```bash
# Native (recommended)
go install github.com/Redsskull/wclogs-cli@latest

# Container (no Go required)
git clone https://github.com/Redsskull/wclogs-cli.git && cd wclogs-cli
make container-run ARGS="config"
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

### `wclogs damage [report-code] [fight-id]`
**Purpose**: Display damage done by all players in a fight

**Usage**:
```bash
# Native binary
./wclogs damage <report-code> <fight-id> [flags]

# Container
make container-run ARGS="damage <report-code> <fight-id> [flags]"
```

**Flags**:
- `--top N` - Show only top N players (default: all)
- `--player "Name"` - Show only specific player
- `--output file.csv` - Save to file (CSV/JSON supported)
- `--no-color` - Disable colored output
- `--verbose` - Show detailed progress

### `wclogs healing [report-code] [fight-id]`
**Purpose**: Display healing done by all players in a fight

**Usage**:
```bash
# Native binary
./wclogs healing <report-code> <fight-id> [flags]

# Container
make container-run ARGS="healing <report-code> <fight-id> [flags]"
```

**Flags**: Same as damage command

---

## 💀 Advanced Analysis Commands

### `wclogs deaths [report-code] [fight-id]`
**Purpose**: Advanced death analysis using Events API with real ability names

**Two Modes**:
1. **Summary Mode** (default): Overview of all deaths
2. **Detailed Mode** (`--player` flag): Deep analysis for specific player

**Usage**:
```bash
# Native binary
./wclogs deaths <report-code> <fight-id> [flags]

# Container
make container-run ARGS="deaths <report-code> <fight-id> [flags]"
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

### `wclogs interrupts [report-code] [fight-id]`
**Purpose**: Professional-grade interrupt analysis with spell correlation and missed opportunity tracking

**Two Modes**:
1. **Summary Mode** (default): Overview of all interrupts and spell analysis
2. **Player Mode** (`--player` flag): Detailed analysis for specific player

**Usage**:
```bash
# Native binary
./wclogs interrupts <report-code> <fight-id> [flags]

# Container
make container-run ARGS="interrupts <report-code> <fight-id> [flags]"
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
./wclogs interrupts YMRqjzC2WPnhwNJd 2 --player "BlagZeras"  # Player-specific detailed analysis
./wclogs interrupts YMRqjzC2WPnhwNJd 2 --verbose            # Verbose mode with API progress

# Container
make container-run ARGS="interrupts YMRqjzC2WPnhwNJd 2"
make container-run ARGS="interrupts YMRqjzC2WPnhwNJd 2 --player BlagZeras"
make container-run ARGS="interrupts YMRqjzC2WPnhwNJd 2 --verbose"
```

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

**"Player 'Name' not found"**
- Use exact player name (case-sensitive)
- Check spelling and special characters

### Debug Mode
Add `--verbose` to any command for detailed debugging:
```bash
# Native binary
./wclogs deaths ABC123 5 --verbose

# Container
make container-run ARGS="deaths ABC123 5 --verbose"
```
Shows API calls, response sizes, and processing steps.

---

## ❌ Commands Not Yet Implemented

| Command | Status | Planned |
|---------|--------|---------|
| `players` | ❌ Missing | Future |
| `timeline` | ❌ Not implemented | Future |
| `boss-abilities` | ❌ Not implemented | Future |

For complete usage examples and detailed command information, see the [API Usage Examples](./docs/api_usage_examples.md) in the docs directory.