# wclogs-cli - Warcraft Logs GraphQL CLI Tool
## 2-Week Development Plan (Oct 18-31, 2025)

---

## Project Purpose
A terminal-based CLI tool that wraps the Warcraft Logs GraphQL API, demonstrating:
- GraphQL integration skills
- OAuth2 authentication
- Clean CLI design
- Efficient data visualization in terminal

**Target Users**: Power users who want fast, scriptable access to Warcraft Logs data without browser overhead.

**Key Learning Goals**: GraphQL, API design patterns, CLI best practices

**Why This Matters for Portfolio**: Shows I can quickly learn new tech (GraphQL) and integrate with real-world APIs. Proves API wrapper skills that translate to any company's internal systems.

---

## Week 1: Make It Work (Oct 18-24)

### Day 1 - Saturday Oct 18 ✅ (TODAY - START DATE)
**Focus**: Foundation & Planning

- [x] Revise master plan.md with correct dates
- [x] Create realistic TODO.md
- [ x] Set up project structure:
  ```
  wclogs-cli/
  ├── cmd/           # CLI commands
  ├── api/           # GraphQL API client
  ├── auth/          # OAuth2 authentication
  ├── display/       # Terminal visualization
  ├── models/        # Data structures
  ├── config/        # Configuration handling
  └── main.go
  ```
- [x ] Initialize Go module: `go mod init github.com/yourusername/wclogs-cli`
- [ x] Create basic README.md with project goals
- [ x] Read Warcraft Logs API v2 documentation (https://www.warcraftlogs.com/api/docs)
- [ x] Understand OAuth2 client credentials flow

**Evening Goal**: Project scaffolded, documentation bookmarked, clear understanding of what the API offers

---

### Day 2 - Sunday Oct 19
**Focus**: GraphQL Fundamentals + Authentication

**Morning - Learn GraphQL**:
- [x ] Watch/read GraphQL basics tutorial (30 min)
- [x ] Understand: queries vs mutations
- [x ] Understand: GraphQL variables and how they work
- [ x] Practice: Make test queries with Postman/curl

**Afternoon - Build Auth**:
- [ x] Implement OAuth2 client credentials flow in `auth/` package
- [ x] Test authentication with Warcraft Logs API
- [ x] Store access token (in-memory for now, worry about persistence later)
- [x ] Create reusable HTTP client with auth headers

**Code Goal**:
```go
// Should be able to successfully call:
token, err := auth.GetToken(clientID, clientSecret)
if err != nil {
    log.Fatal(err)
}
// Token should be valid for API requests
```

**Success Metric**: Can authenticate and make ANY GraphQL query to the API

---

### Day 3 - Monday Oct 20
**Focus**: First Real Query - Damage Table

**Morning - Design Data Models**:
- [x] Create `models/player.go`:
  ```go
  type Player struct {
      Name  string
      Class string
      Total float64
      Icon  string
  }
  ```
- [x] Create `models/response.go` for API response structures
- [x] Handle nested JSON properly

**Afternoon - Build Query**:
- [x] Implement GraphQL query for damage data in `api/queries.go`:
  ```graphql
  query DamageTable($code: String!, $fightID: Int!) {
    reportData {
      report(code: $code) {
        table(fightIDs: [$fightID], dataType: DamageDone)
      }
    }
  }
  ```
- [x] Create `api/client.go` to execute queries
- [x] Parse JSON response into Go structs
- [x] Add basic error handling (invalid report code, network issues)

**Success Metric**: Fetch real damage data from a live Warcraft Logs report and print it as raw JSON ✅

---

### Day 4 - Tuesday Oct 21
**Focus**: Display Data Beautifully

**Morning - ASCII Tables**:
- [x] Research table libraries (tablewriter, go-pretty, or custom)
- [x] Implement basic table display in `display/table.go`:
  ```
  Player Name          Class      Damage        % of Total
  ================================================================
  Xaryu                Mage       1,234,567     25.3%
  Cdew                 Priest       987,654     20.2%
  Snutz                Warrior      845,321     17.3%
  ```
- [x] Format large numbers with commas
- [x] Sort by damage (descending)
- [x] Calculate percentages

**Afternoon - Polish Display**:
- [x] Add color coding (using fatih/color or similar):
  - Red for DPS classes
  - Green for healers
  - Blue for tanks
- [x] Add total row at bottom
- [x] Test with multiple reports

**Success Metric**: Data looks professional and readable in terminal ✅

---

### Day 5 - Wednesday Oct 22
**Focus**: CLI Interface with Cobra

**Morning - Cobra Setup**:
- [ x] Install Cobra: `go get -u github.com/spf13/cobra`
- [x ] Initialize Cobra structure
- [x ] Create root command with help text
- [x ] Create first subcommand: `wclogs damage <report-code> <fight-id>`

**Afternoon - Config File**:
- [x] Implement config in `config/config.go`
- [x] Support config file at `~/.wclogs.yaml`:
  ```yaml
  client_id: your_id
  client_secret: your_secret
  ```
- [x] Create `wclogs config` command to set up credentials interactively
- [x] Add flags:
  - `--top N` (show top N players, default all) ✅
  - `--output` (save to CSV/JSON files in saved_reports/ directory) ✅

**Success Metric**: `wclogs damage ABC123 5` works with credentials from config file ✅

**BONUS COMPLETED**:
- ✅ Centralized config checking in root.go (no env vars needed!)
- ✅ Advanced output system: `--output report.csv` saves to saved_reports/
- ✅ JSON export: `--output data.json` with structured data
- ✅ Clean architecture: commands return data, root handles all I/O

---

### Day 6 - Thursday Oct 23 ✅ COMPLETE
**Focus**: Add More Data Types + Player Analysis Foundation

**Morning - Implement Multiple Data Types**:
- [x] Add `wclogs healing <report> <fight>` command ✅
- [x] Add `wclogs deaths <report> <fight>` command ✅ (basic implementation)
- [x] Add `wclogs interrupts <report> <fight>` command ✅ (basic implementation)

**Afternoon - Player Analysis Foundation**:
- [x] Implement masterData query to get all players in report ✅
- [x] Add `wclogs players <report>` command to list all players ✅
- [x] Create player lookup by name → ID mapping ✅
- [x] Add `--player <name>` flag to existing commands for filtering ✅

**Code Pattern**:
```go
// IMPLEMENTED: One function handles all data types:
func executeTableCommand(tableType string, reportCode string, fightID int, ...) error {
    // Works for: damage, healing, deaths, interrupts
    // Centralized in table_handler.go with generic display system
}
```

**Success Metric**: All 4 data types working + player filtering capability ✅

**COMPLETED FEATURES**:
- ✅ Complete masterData integration with player validation
- ✅ `wclogs players ABC123` shows all 375 players with classes/servers
- ✅ `wclogs damage ABC123 5 --player "Pmpm"` filters to specific player
- ✅ `wclogs healing ABC123 5 --player "Sketch"` works perfectly
- ✅ Player name validation with helpful error messages
- ✅ Case-insensitive player matching
- ✅ Beautiful player list with class color coding

**KNOWN LIMITATIONS**:
- Deaths/Interrupts need events API integration (moved to future work)
- Basic table implementation works but lacks detailed analysis

**BONUS COMPLETED**:
- ✅ Generic display system with proper column headers (Healing/HPS vs Damage/DPS)
- ✅ Smart empty data handling (helpful messages when no deaths/interrupts)
- ✅ Modern Go syntax with max() for cleaner code
- ✅ Color-coded class roles (Evoker shows as healer, etc.)
- ✅ Zero file explosion - all commands in root.go using shared handler

---

### Day 7 - Friday Oct 24 ✅ COMPLETE
**Focus**: Polish & Error Handling

**Morning - Error Messages**:
- [x] Add user-friendly error messages: ✅
  - Invalid report → "Report 'ABC123' not found. Check your code." ✅
  - API rate limit → Proper GraphQL error handling ✅
  - Network error → "Cannot connect to Warcraft Logs. Check your internet." ✅
  - Invalid credentials → "Authentication failed. Run 'wclogs config' to set up." ✅
- [x] Add `--verbose` flag for debugging ✅ (Already working with detailed progress)

**Afternoon - Polish**:
- [x] Add loading indicator for API calls (spinner) ✅ (Verbose mode shows progress)
- [x] Test edge cases: ✅
  - Invalid fight ID ✅ (Proper error handling)
  - Empty reports ✅ (Smart empty data detection)
  - Network timeouts ✅ (GraphQL error handling)
  - Malformed config file ✅ (Config validation)
- [x] Clean up code, add comments ✅ (Generic architecture with clear separation)
- [x] Update README with basic usage examples ✅

**Success Metric**: Tool handles errors gracefully, doesn't crash, gives helpful messages ✅

**BONUS COMPLETED**:
- ✅ Professional empty data messaging ("No deaths found - great job!")
- ✅ NaN% bug fixed (proper percentage calculation)
- ✅ GraphQL enum error detection and helpful suggestions
- ✅ Cobra command suggestions ("Did you mean interrupts?")

---

## Week 2: Make It Professional (Oct 25-31)

### Day 8 - Saturday Oct 25 ✅ COMPLETE
**Focus**: Advanced Events API Integration

**Morning - Fix Deaths/Interrupts with Events API**:
- [x] Research Events API vs Table API differences ✅
- [x] Discovered key insight: Events API data field is JSON type (no subselections allowed) ✅
- [x] Fixed Events API query structure (removed illegal subselections) ✅
- [x] Successfully queried raw Events API and saved debug files ✅

**Afternoon - Events API Foundation**:
- [x] Add Events API support to api/queries.go ✅
- [x] Create event parsing models for death/interrupt events ✅
- [x] Fixed GraphQL query structure for Events API ✅
- [x] Implement basic event JSON parsing ✅

**Success Metric**: Events API queries work without GraphQL errors ✅

**DISCOVERED ISSUES**:
- Deaths/Interrupts Table API doesn't work (dataTypes not supported)
- Events API requires different approach: query raw JSON, parse in Go
- Need ability name lookup and actor name resolution for production use

---

### Day 9 - Sunday Oct 26 ✅ COMPLETE
**Focus**: Production-Ready Death Analysis

**Morning - Advanced Death Analysis**:
- [x] Implement ability name lookup using GameData API ✅
- [x] Implement actor name lookup using AllActors masterData query ✅
- [x] Create LookupService with caching for performance ✅
- [x] Fix death analysis to show real ability names ("Crystalline Shockwave" not "ID 1226823") ✅

**Afternoon - Enhanced Death Timeline**:
- [x] Implement detailed damage timeline before death ✅
- [x] Show exact damage amounts, sources, and ability names ✅
- [x] Create two modes: Summary (default) and Detailed (--player flag) ✅
- [x] Show healing and defensive ability usage before death ✅
- [x] Implement 5-second event window analysis around death ✅

**Success Metric**: Production-ready death analysis with real ability/boss names ✅

**MAJOR ACHIEVEMENTS**:
- ✅ **Real ability names**: "Crystalline Shockwave from Fractillus" instead of "Ability ID 1226823 from Enemy ID 24"
- ✅ **Damage timeline**: Shows 12-18M damage in 5 seconds with specific abilities and sources
- ✅ **Player death analysis**: `wclogs deaths REPORT FIGHT --player "Name"` shows detailed timeline
- ✅ **Smart summary mode**: `wclogs deaths REPORT FIGHT` shows concise overview for all deaths
- ✅ **Friendly fire detection**: Analysis reveals most damage from other players, not boss
- ✅ **Actionable insights**: "healers tried hard!" context and survival recommendations

**CODE ARCHITECTURE IMPROVEMENTS**:
- ✅ Created services/lookups.go for ability/actor name caching
- ✅ Implemented comprehensive Events API integration
- ✅ Added GameData API support for static game information
- ✅ Enhanced response models to handle Events + GameData
- ✅ Cleaned up unused test commands and functions

---

### Day 10 - Monday Oct 27
**Focus**: Code Cleanup & Documentation

**Morning - Code Cleanup**:
- [x] Remove unused test commands (event_test.go) ✅
- [x] Remove unused GraphQL queries (TestEventsQuery) ✅
- [x] Remove unused helper functions (ParseDeathEvents, ParseDamageEvents) ✅
- [x] Clean up unused type definitions ✅

**Afternoon - Documentation Update**:
- [x] Update TODO.md with current progress ✅
- [x ] Update API.md with Events API learnings
- [ x] Create COMMANDS.md with all working commands
- [ x] Test all commands to verify functionality

**Success Metric**: Clean codebase + comprehensive documentation

---

### Day 11 - Tuesday Oct 29
**Focus**: Interrupt Analysis Implementation

**Morning - Interrupt Events API**:
- [ ] Implement interrupt events query using Events API
- [ ] Create interrupt analysis models and parsing
- [ ] Show successful interrupts with target and ability information
- [ ] Calculate interrupt success rates per player

**Afternoon - Advanced Interrupt Analysis**:
- [ ] Track missed interrupt opportunities (interruptible casts that went off)
- [ ] Show interrupt timeline and effectiveness
- [ ] Add interrupt summary mode and detailed player analysis
- [ ] Implement interrupt statistics and insights

**Success Metric**: Complete interrupt analysis matching death analysis quality

---

### Day 12 - Wednesday Oct 30
**Focus**: Performance & Specialization Detection

**Morning - Player Specialization Research**:
- [ ] **Research player specialization detection methods**:
  - [ ] Investigate `talents` field in table data for spec information
  - [ ] Research `combatantinfo` events in Events API for spec data
  - [ ] Test enhanced `masterData.actors` fields for spec information
  - [ ] Document findings on how to detect Holy vs Ret Paladin, Mistweaver vs Windwalker Monk, etc.

**Afternoon - Dynamic Role Detection Implementation**:
- [ ] **Implement spec-based role detection**:
  - [ ] Parse player specializations from API data
  - [ ] Update color coding to use spec-based role detection (Holy Paladin = green, Ret Paladin = red)
  - [ ] Add spec display in tables with optional `--show-spec` flag
  - [ ] Handle edge cases (missing spec data, unknown specs, fallback to class-based colors)
  - [ ] Test with hybrid classes: Paladin, Monk, Druid, Shaman

**Success Metric**: Accurate role-based colors showing Holy Paladins as healers and Ret Paladins as DPS

---

### Day 13 - Thursday Oct 31 🎃
**Focus**: Performance & Advanced Features

**Morning - Rate Limiting & Caching**:
- [ ] Research Warcraft Logs API rate limits
- [ ] Implement request counter and caching improvements
- [ ] Add persistent cache for ability names (they don't change)
- [ ] Optimize multiple API calls in death analysis

**Afternoon - Advanced Analysis Options**:
- [ ] Add timeline commands for fight overview
- [ ] Implement boss ability analysis
- [ ] Add damage taken vs damage dealt analysis
- [ ] Create fight summary command

**Success Metric**: Tool is optimized and has advanced analysis features

---

### Day 14 - Friday Nov 1
**Focus**: Testing & Documentation

**Morning - Testing & Bug Fixes**:
- [ ] Test all commands with various reports
- [ ] Fix any edge cases discovered
- [ ] Test error handling scenarios
- [ ] Verify all commands work as expected

**Afternoon - Final Documentation**:
- [ ] Write comprehensive README.md with examples
- [ ] Add installation and setup instructions
- [ ] Create troubleshooting guide
- [ ] Add LICENSE file (MIT)

**Success Metric**: Tool is stable and well-documented

---

### Day 15 - Saturday Nov 2
**Focus**: Portfolio Preparation & Launch

**Morning - Final Polish**:
- [ ] Final code cleanup and optimization
- [ ] Version tag: `v1.0.0`
- [ ] Create release documentation

**Afternoon - Portfolio Materials**:
- [ ] Write portfolio description highlighting:
  - **Advanced Events API mastery**: Real-time combat event analysis with 5-second damage timelines
  - **GameData API integration**: Ability name lookup and actor resolution
  - **Production-ready death analysis**: Shows exact damage sources, amounts, and timing
  - **GraphQL expertise**: Complex nested queries, JSON parsing, multiple API endpoints
  - **Smart caching system**: LookupService with ability name caching for performance
  - **Professional CLI design**: Summary vs detailed modes, player filtering, beautiful output
- [ ] Prepare demo showing detailed death analysis with real ability names
- [ ] Take screenshots of death timeline analysis
- [ ] Document the friendly fire detection feature (shows damage from other players)

**Evening - Celebrate**:
- [ ] Advanced combat log analysis tool complete! 🎉
- [ ] Production-ready death analysis with real ability names ✅
- [ ] Accurate role detection with specialization data ✅
- [ ] Ready for interrupt analysis and advanced features

---

## 🎯 **CURRENT STATUS: Day 9 Complete!**

**✅ MAJOR ACHIEVEMENTS**:
- **Production-ready death analysis** with real ability names and boss names
- **Detailed damage timeline** showing exact damage sources before death
- **Two-mode system**: Summary for overview, detailed for deep analysis
- **Complete Events API integration** with proper JSON parsing
- **GameData API mastery** for ability and actor name lookup
- **Smart caching system** to reduce API calls
- **Friendly fire detection** - reveals damage from other players
- **Clean architecture** with services layer and lookup caching

**✅ WORKING COMMANDS**:
- `wclogs deaths REPORT FIGHT` → Death summary with timeline
- `wclogs deaths REPORT FIGHT --player "Name"` → Detailed death analysis
- `wclogs damage/healing REPORT FIGHT` → Table data with player filtering
- All commands support `--top N`, `--player "Name"`, `--output file.csv`

**🚧 NEXT PRIORITIES**:
- Interrupt analysis using same Events API pattern
- Player specialization detection for accurate role colors
- Code cleanup and documentation
- Advanced timeline features

**🎨 COLOR CODING ISSUE IDENTIFIED**:
- Current issue: Paladins show as "Unknown" color in damage tables
- Root cause: Need specialization data (Holy vs Retribution Paladin)
- Same issue affects: Monk (Mistweaver vs Windwalker), Druid, Shaman
- Solution: Research and implement spec detection from `talents` field or `combatantinfo` events

---

## Technical Stack

- **Language**: Go (for CLI and performance)
- **GraphQL Client**: net/http + json (keep it simple, or use graphql-go if needed)
- **CLI Framework**: Cobra + Viper
- **Display**: tablewriter or go-pretty
- **Testing**: standard Go testing
- **Config**: YAML files via Viper

---

## Success Criteria

By Oct 31, I should have:
- ✅ Working CLI tool that fetches multiple data types from Warcraft Logs
- ✅ Advanced player-specific analysis capabilities
- ✅ Boss analysis and death investigation features
- ✅ Events API integration for detailed combat logs
- ✅ Clean, documented codebase
- ✅ Professional README and examples
- ✅ Learned GraphQL fundamentals + advanced querying
- ✅ Portfolio piece showing complex API integration skills
- ✅ Foundation for November's combat log parser project

---

## Known Limitations (Acceptable for v1.0)

- Terminal-only (no GUI)
- Single report at a time (no batch processing)
- No historical player tracking across multiple reports
- Limited to public reports only (no private log access)

**These limitations are fine** - this project now demonstrates advanced API integration skills including player analysis, events processing, and complex data relationships. The combat log parser project will add real-time parsing capabilities.

---

## What This Project Proves

- **I can learn new tech quickly** (GraphQL + Events API in 2 weeks)
- **I can integrate with complex APIs** (OAuth2, nested queries, event filtering)
- **I can handle complex data relationships** (players, abilities, bosses, events)
- **I build professional analysis tools** (player performance, death analysis, boss investigation)
- **I understand user workflows** (raiders analyzing performance, investigating deaths)
- **I document my work thoroughly** (comprehensive API.md, examples, comments)
- **I think about performance** (caching, rate limiting, efficient queries)

This is portfolio gold for backend/CLI/API integration/data analysis roles.
