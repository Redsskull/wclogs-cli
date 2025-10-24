# wclogs-cli Commands Reference

Complete reference for all available commands in the Warcraft Logs CLI tool.

---

## 📋 Command Overview

| Command | Status | Description |
|---------|--------|-------------|
| `config` | ✅ Working | Set up API credentials |
| `damage` | ✅ Working | Show damage tables with player filtering |
| `healing` | ✅ Working | Show healing tables with player filtering |
| `deaths` | ✅ Working | Advanced death analysis with Events API |
| `help` | ✅ Working | Show help for commands |
| `completion` | ✅ Working | Generate shell completions |

---

## 🔧 Setup Commands

### `wclogs config`
**Purpose**: Interactive setup of Warcraft Logs API credentials

**Usage**:
```bash
wclogs config
```

**What it does**:
- Prompts for Client ID and Client Secret
- Saves credentials to `~/.wclogs.yaml`
- Tests authentication with the API

**Example**:
```bash
$ wclogs config
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
wclogs damage <report-code> <fight-id> [flags]
```

**Flags**:
- `--top N` - Show only top N players (default: all)
- `--player "Name"` - Show only specific player
- `--output file.csv` - Save to file (CSV/JSON supported)
- `--no-color` - Disable colored output
- `--verbose` - Show detailed progress

**Examples**:
```bash
# Basic damage table
wclogs damage 6qNJmgYBTcyfvpWF 3

# Top 5 DPS players
wclogs damage 6qNJmgYBTcyfvpWF 3 --top 5

# Specific player only
wclogs damage 6qNJmgYBTcyfvpWF 3 --player "Pherally"

# Save to CSV file
wclogs damage 6qNJmgYBTcyfvpWF 3 --output damage_report.csv

# Save to JSON file
wclogs damage 6qNJmgYBTcyfvpWF 3 --output damage_data.json
```

**Sample Output**:
```
🗡️ DAMAGE TABLE 🗡️

Player Name   Class            Damage        DPS   % Total
==========================================================
Pherally      Warrior   1,639,721,988  5,226,588     35.1%
Nikkans       Paladin   1,540,555,415  4,893,184     32.9%
Rach          Paladin   1,495,729,438  4,756,985     32.0%
==========================================================
📊 Showing top 3 of 20 players | Total Damage: 4,676,006,841
```

---

### `wclogs healing [report-code] [fight-id]`
**Purpose**: Display healing done by all players in a fight

**Usage**:
```bash
wclogs healing <report-code> <fight-id> [flags]
```

**Flags**: Same as damage command

**Examples**:
```bash
# Basic healing table
wclogs healing 6qNJmgYBTcyfvpWF 3

# Top 3 healers
wclogs healing 6qNJmgYBTcyfvpWF 3 --top 3

# Specific healer analysis
wclogs healing 6qNJmgYBTcyfvpWF 3 --player "Hanahime"
```

**Sample Output**:
```
💚 HEALING TABLE 💚

Player Name   Class           Healing        HPS   % Total
==========================================================
Hanahime      Monk      1,673,061,580  5,318,549     37.3%
Truxpriest    Priest    1,445,207,171  4,624,529     32.2%
Hejblx        Evoker    1,366,204,988  4,360,261     30.5%
==========================================================
📊 Showing top 3 of 20 players | Total Healing: 4,484,473,739
```

---

## 💀 Advanced Analysis Commands

### `wclogs deaths [report-code] [fight-id]`
**Purpose**: Advanced death analysis using Events API with real ability names

**Two Modes**:
1. **Summary Mode** (default): Overview of all deaths
2. **Detailed Mode** (`--player` flag): Deep analysis for specific player

**Usage**:
```bash
wclogs deaths <report-code> <fight-id> [flags]
```

**Flags**:
- `--player "Name"` - Detailed analysis for specific player
- `--verbose` - Show debug information and API progress
- `--output file.json` - Save analysis to file

**Examples**:

#### Summary Mode (Default)
```bash
wclogs deaths 6qNJmgYBTcyfvpWF 3
```

**Sample Output**:
```
💀 DEATH ANALYSIS SUMMARY 💀
Fight: Fractillus (Duration: 5m15.349s)
Result: SUCCESS ✅
Deaths: 9

📅 DEATH TIMELINE:
  • 84s: Disfatbidge
  • 182s: Tekkyysp
  • 248s: Tekkyysp
  • 312s: White, Amberlotrev, Hanahime, Shaepeshift, Bræt (5 players)
  • 315s: Nikkans

⚔️  TOP KILLING ABILITIES:
  • Crystalline Shockwave: 5 deaths
  • Crystalline Shockwave: 2 deaths
  • Null Explosion: 1 deaths
  • Nexus Shrapnel: 1 deaths

💡 TIP: Use --player "PlayerName" for detailed death analysis of a specific player
```

#### Detailed Mode (Player-Specific)
```bash
wclogs deaths 6qNJmgYBTcyfvpWF 3 --player "Tekkyysp"
```

**Sample Output**:
```
💀 DETAILED DEATH ANALYSIS: Tekkyysp 💀
Fight: Fractillus (Duration: 5m15.349s)
Deaths: 2

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
💀 Death #1
  ⏱️  Survival Time: 3m2.087s
  ⚔️  Killed by: Crystalline Shockwave from Fractillus
  📈 Events Around Death:
    • -5.0s: 423,292 damage from Rach (Crusading Strikes)
    • -4.9s: 909,791 damage from Shibawar (Shield Slam)
    • -4.9s: 1,759,236 damage from Shibawar (Ire of Devotion)
    • -4.8s: 6,516,897 damage from Bræt (Earth Shock)
    📊 Total damage in window: 12,522,701 (17 events)
  💚 Healing Analysis:
    • Total healing: 3,179,019 (healers tried hard!)
  🛡️  Defensive Analysis:
    • Used 4 defensive abilities

📊 INSIGHTS:
• Tekkyysp died 2 times - focus on mechanics and survival
```

**Key Features**:
- ✅ **Real ability names**: Shows "Crystalline Shockwave from Fractillus" not "Ability ID 1226823"
- ✅ **Damage timeline**: Exact damage amounts and sources in 5-second death window
- ✅ **Friendly fire detection**: Shows damage from other players
- ✅ **Healing context**: Shows healing attempts with contextual insights
- ✅ **Survival analysis**: Calculates correct survival times from fight start

---

## 🌐 Global Flags

All commands support these global flags:

| Flag | Short | Description | Example |
|------|-------|-------------|---------|
| `--output` | `-o` | Save to file (CSV/JSON) | `--output report.csv` |
| `--top` | `-t` | Show top N players | `--top 5` |
| `--verbose` | `-v` | Enable verbose output | `--verbose` |
| `--help` | `-h` | Show command help | `--help` |

---

## 🎯 File Output Formats

### CSV Output
```bash
wclogs damage ABC123 5 --output damage.csv
```
Creates structured CSV with columns: Name, Class, Damage, DPS, Percentage

### JSON Output  
```bash
wclogs deaths ABC123 5 --output deaths.json
```
Creates structured JSON with all analysis data for programmatic use

**Output Location**: All files saved to `saved_reports/` directory

---

## ❌ Commands Not Yet Implemented

| Command | Status | Planned |
|---------|--------|---------|
| `interrupts` | ❌ Not working | Day 11 |
| `players` | ❌ Missing | Day 11 |
| `timeline` | ❌ Not implemented | Day 12 |
| `boss-abilities` | ❌ Not implemented | Day 12 |

---

## 🔧 Troubleshooting

### Common Errors

**"Authentication failed"**
```bash
# Run config setup
wclogs config
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
wclogs deaths ABC123 5 --verbose
```
Shows API calls, response sizes, and processing steps.

---

## 💡 Usage Tips

### Finding Report Codes
Report codes are in Warcraft Logs URLs:
`https://www.warcraftlogs.com/reports/Hw9TZc2WyrVKJLCa` → Code: `Hw9TZc2WyrVKJLCa`

### Finding Fight IDs  
Use damage/healing commands to see available fights, then use specific fight ID for death analysis.

### Player Name Filtering
All table commands support `--player "Name"` for focused analysis:
```bash
wclogs damage ABC123 5 --player "Tankadin"
wclogs healing ABC123 5 --player "Healbot" 
wclogs deaths ABC123 5 --player "Dpswarrior"
```

### Performance Tips
- Death analysis caches ability names automatically
- Use `--top N` for faster results with large raids
- JSON output is faster than CSV for large datasets

---

## 📚 Examples by Use Case

### Raid Leader Analysis
```bash
# Quick overview of fight performance
wclogs damage 6qNJmgYBTcyfvpWF 3 --top 10
wclogs healing 6qNJmgYBTcyfvpWF 3 --top 5

# Death investigation
wclogs deaths 6qNJmgYBTcyfvpWF 3

# Individual player review  
wclogs deaths 6qNJmgYBTcyfvpWF 3 --player "Strugglingdps"
```

### Personal Performance Review
```bash
# My damage performance
wclogs damage 6qNJmgYBTcyfvpWF 3 --player "Mycharacter"

# How did I die?
wclogs deaths 6qNJmgYBTcyfvpWF 3 --player "Mycharacter" --verbose
```

### Data Export for Spreadsheets
```bash
# Export all data for analysis
wclogs damage 6qNJmgYBTcyfvpWF 3 --output raid_damage.csv
wclogs healing 6qNJmgYBTcyfvpWF 3 --output raid_healing.csv
wclogs deaths 6qNJmgYBTcyfvpWF 3 --output death_analysis.json
```

---

## 🎯 Command Success Matrix

| Command | Basic Usage | Player Filter | File Output | Verbose Mode |
|---------|-------------|---------------|-------------|--------------|
| `damage` | ✅ | ✅ | ✅ | ✅ |
| `healing` | ✅ | ✅ | ✅ | ✅ |
| `deaths` | ✅ | ✅ | ✅ | ✅ |
| `config` | ✅ | N/A | N/A | N/A |

**Legend**: ✅ Working | ❌ Not implemented | N/A Not applicable