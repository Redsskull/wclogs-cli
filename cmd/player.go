package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gookit/color"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/models"
	"wclogs-cli/output"
)

// ExecutePlayersCommand handles the players command with full filtering support
func ExecutePlayersCommand(reportCode, fightIDStr, classFilter, roleFilter, searchFilter, outputPath string, topN int, noColor, verbose, debug bool) error {
	if verbose {
		if fightIDStr != "" {
			color.Cyan.Printf("🔍 Fetching player list for report %s, fight %s\n", reportCode, fightIDStr)
		} else {
			color.Cyan.Printf("🔍 Fetching player list for report %s\n", reportCode)
		}
		if classFilter != "" {
			color.Cyan.Printf("   - Class filter: %s\n", classFilter)
		}
		if roleFilter != "" {
			color.Cyan.Printf("   - Role filter: %s\n", roleFilter)
		}
		if searchFilter != "" {
			color.Cyan.Printf("   - Search filter: %s\n", searchFilter)
		}
	}

	// Auth logic
	if verbose {
		color.Cyan.Printf("🔐 Loading configuration...\n")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// API client setup
	if verbose {
		color.Cyan.Printf("🔐 Setting up authentication...\n")
	}
	authClient := auth.NewClient(cfg.ClientID, cfg.ClientSecret)

	if verbose {
		color.Cyan.Printf("📡 Creating API client...\n")
	}
	apiClient := api.NewClient(authClient)

	// Validation
	if verbose {
		color.Cyan.Printf("✅ Validating parameters...\n")
	}
	if reportCode == "" {
		return fmt.Errorf("report code cannot be empty")
	}

	if len(reportCode) < 6 {
		return fmt.Errorf("report code '%s' is too short (must be at least 6 characters)", reportCode)
	}

	// Query execution - different approach for fight-specific vs report-wide
	var actualPlayers []models.Actor

	if fightIDStr != "" {
		// Fight-specific player lookup
		if verbose {
			color.Cyan.Printf("🚀 Resolving fight ID and fetching fight participants...\n")
		}

		// Resolve fight ID (handles both numbers and "last" keyword)
		fightID, err := resolveFightID(reportCode, fightIDStr, verbose)
		if err != nil {
			return fmt.Errorf("failed to resolve fight ID '%s': %w", fightIDStr, err)
		}

		actualPlayers, err = getFightParticipants(apiClient, reportCode, fightID, verbose)
		if err != nil {
			return fmt.Errorf("failed to get fight participants: %w", err)
		}
	} else {
		// Report-wide player lookup
		if verbose {
			color.Cyan.Printf("🚀 Executing masterData GraphQL query...\n")
		}

		request := api.NewMasterDataRequest(reportCode)
		response, err := apiClient.Query(request.Query, request.Variables)
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}

		// Validate response structure
		if response.Data == nil || response.Data.ReportData == nil || response.Data.ReportData.Report == nil {
			return fmt.Errorf("no report data found for code: %s", reportCode)
		}

		if response.Data.ReportData.Report.MasterData == nil {
			return fmt.Errorf("no master data found for report %s", reportCode)
		}

		masterData := response.Data.ReportData.Report.MasterData
		if len(masterData.Actors) == 0 {
			return fmt.Errorf("no actors found in report %s", reportCode)
		}

		if verbose {
			color.Cyan.Printf("📊 Found %d total actors in the report\n", len(masterData.Actors))
		}

		// Filter to only actual players (those with valid class information)
		actualPlayers = filterToActualPlayers(masterData.Actors, verbose)
	}

	if len(actualPlayers) == 0 {
		if fightIDStr != "" {
			return fmt.Errorf("no players found in fight %s of report %s", fightIDStr, reportCode)
		} else {
			return fmt.Errorf("no players found in report %s", reportCode)
		}
	}

	if verbose {
		color.Green.Printf("✅ Found %d actual players\n", len(actualPlayers))
	}

	// Create enhanced player list with role detection
	players := createEnhancedPlayerListWithDebug(actualPlayers, debug)

	// Apply filters
	filteredPlayers := applyFilters(players, classFilter, roleFilter, searchFilter)

	if len(filteredPlayers) == 0 {
		return fmt.Errorf("no players match the specified filters")
	}

	if verbose {
		color.Green.Printf("🎯 Filtered to %d player(s)\n", len(filteredPlayers))
	}

	// Apply top N limit
	if topN > 0 && topN < len(filteredPlayers) {
		filteredPlayers = filteredPlayers[:topN]
		if verbose {
			color.Cyan.Printf("📊 Limited to top %d players\n", topN)
		}
	}

	// Handle output
	if outputPath != "" {
		// File output
		return outputPlayersToFile(filteredPlayers, reportCode, outputPath, verbose)
	} else {
		// Terminal output
		displayPlayersInTerminal(filteredPlayers, reportCode, fightIDStr, classFilter, roleFilter, searchFilter, noColor)
		return nil
	}
}

// EnhancedPlayer represents a player with additional computed fields
type EnhancedPlayer struct {
	*models.PlayerInfo
	Role string `json:"role"`
}

// createEnhancedPlayerList creates enhanced player objects with role detection
func createEnhancedPlayerList(actors []models.Actor) []*EnhancedPlayer {
	return createEnhancedPlayerListWithDebug(actors, false)
}

// createEnhancedPlayerListWithDebug creates enhanced player objects with optional debug output
func createEnhancedPlayerListWithDebug(actors []models.Actor, debug bool) []*EnhancedPlayer {
	players := make([]*EnhancedPlayer, 0, len(actors))

	for _, actor := range actors {
		player := &EnhancedPlayer{
			PlayerInfo: &models.PlayerInfo{
				ID:     actor.ID,
				Name:   actor.Name,
				Class:  actor.SubType, // SubType contains the class name
				Server: actor.Server,
				Icon:   actor.Icon,
			},
			Role: detectRoleWithDebug(actor.SubType, actor.Icon, debug),
		}
		players = append(players, player)
	}

	// Sort by name
	sort.Slice(players, func(i, j int) bool {
		return strings.ToLower(players[i].Name) < strings.ToLower(players[j].Name)
	})

	return players
}

// getFightParticipants gets players who participated in a specific fight with full metadata
// Uses hybrid approach: table data for participants + MasterData for complete info
func getFightParticipants(apiClient *api.Client, reportCode string, fightID int, verbose bool) ([]models.Actor, error) {
	if verbose {
		color.Cyan.Printf("🔍 Fetching fight participants using hybrid query strategy...\n")
	}

	// Step 1: Get fight participants from damage table
	if verbose {
		color.Cyan.Printf("   📊 Getting fight participants from damage data...\n")
	}
	damageRequest := api.NewTableRequest(reportCode, fightID, api.DataTypeDamage)
	damageResponse, err := apiClient.Query(damageRequest.Query, damageRequest.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fight data: %w", err)
	}

	if damageResponse.Data == nil || damageResponse.Data.ReportData == nil || damageResponse.Data.ReportData.Report == nil {
		return nil, fmt.Errorf("no fight data found")
	}

	rawTable := damageResponse.Data.ReportData.Report.Table
	if len(rawTable) == 0 {
		return nil, fmt.Errorf("no damage data found for fight %d", fightID)
	}

	// Parse damage table to get participant IDs
	tableData, err := models.ParseTableData(rawTable)
	if err != nil {
		return nil, fmt.Errorf("failed to parse table data: %w", err)
	}

	// Create map of fight participant IDs
	participantIDs := make(map[int]bool)
	for _, entry := range tableData.Entries {
		participantIDs[entry.ID] = true
	}

	if verbose {
		color.Cyan.Printf("   👥 Found %d participants in fight %d\n", len(participantIDs), fightID)
	}

	// Step 2: Get full MasterData for complete metadata
	if verbose {
		color.Cyan.Printf("   🔍 Fetching full player metadata from MasterData...\n")
	}
	masterRequest := api.NewMasterDataRequest(reportCode)
	masterResponse, err := apiClient.Query(masterRequest.Query, masterRequest.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch master data: %w", err)
	}

	if masterResponse.Data == nil || masterResponse.Data.ReportData == nil ||
		masterResponse.Data.ReportData.Report == nil ||
		masterResponse.Data.ReportData.Report.MasterData == nil {
		return nil, fmt.Errorf("no master data found")
	}

	// Step 3: Cross-reference to get fight participants with full metadata
	allActors := masterResponse.Data.ReportData.Report.MasterData.Actors
	var fightActors []models.Actor

	for _, actor := range allActors {
		// Only include actors who participated in this fight and have valid class data
		if participantIDs[actor.ID] && actor.SubType != "" && strings.TrimSpace(actor.SubType) != "" {
			fightActors = append(fightActors, actor)
		}
	}

	if verbose {
		color.Green.Printf("✅ Found %d players in fight %d with complete metadata\n", len(fightActors), fightID)
		color.Cyan.Printf("   🎯 Includes: servers, detailed spec icons, and full actor data\n")
	}

	return fightActors, nil
}

// filterToActualPlayers filters actors to only include actual players with valid class data
func filterToActualPlayers(actors []models.Actor, verbose bool) []models.Actor {
	var players []models.Actor
	filteredCount := 0

	for _, actor := range actors {
		// Filter out actors without valid class information
		if actor.SubType == "" || strings.TrimSpace(actor.SubType) == "" {
			filteredCount++
			if verbose && filteredCount <= 5 { // Show first few filtered actors for debugging
				color.Gray.Printf("   Filtering out: %s (no class info)\n", actor.Name)
			}
			continue
		}

		// Filter out obvious non-player entities based on name patterns
		name := strings.ToLower(actor.Name)
		if strings.Contains(name, "mirror image") ||
			strings.Contains(name, "spirit wolf") ||
			strings.Contains(name, "army of the dead") ||
			strings.HasPrefix(name, "unknown") {
			filteredCount++
			if verbose && filteredCount <= 5 {
				color.Gray.Printf("   Filtering out: %s (pet/summon)\n", actor.Name)
			}
			continue
		}

		players = append(players, actor)
	}

	if verbose && filteredCount > 5 {
		color.Gray.Printf("   ... and %d more actors filtered out\n", filteredCount-5)
	}

	return players
}

// detectRole determines the player's role based on class and sometimes icon/spec
func detectRole(class, icon string) string {
	return detectRoleWithDebug(class, icon, false)
}

// detectRoleWithDebug determines the player's role with optional debug output
func detectRoleWithDebug(class, icon string, debug bool) string {
	if debug {
		color.Gray.Printf("   🔍 Role detection: %s (icon: %s)\n", class, icon)
	}
	class = strings.ToLower(class)
	iconLower := strings.ToLower(icon)

	switch class {
	// Pure DPS classes
	case "hunter", "mage", "rogue", "warlock":
		return "DPS"

	// Tank classes (can also be DPS/Healer)
	case "death knight", "deathknight":
		// Blood = Tank, Frost/Unholy = DPS
		if strings.Contains(iconLower, "blood") ||
			strings.Contains(iconLower, "deathknight_blood") ||
			strings.Contains(iconLower, "spell_deathknight_bloodpresence") {
			return "Tank"
		}
		if strings.Contains(iconLower, "frost") || strings.Contains(iconLower, "unholy") ||
			strings.Contains(iconLower, "deathknight_frost") || strings.Contains(iconLower, "deathknight_unholy") {
			return "DPS"
		}
		// If we can't determine spec, return Unknown for better debugging
		return "Unknown"

	case "demon hunter", "demonhunter":
		// Vengeance = Tank, Havoc = DPS
		if strings.Contains(iconLower, "vengeance") ||
			strings.Contains(iconLower, "demonhunter_vengeance") ||
			strings.Contains(iconLower, "ability_demonhunter_vengeance") {
			return "Tank"
		}
		if strings.Contains(iconLower, "havoc") ||
			strings.Contains(iconLower, "demonhunter_havoc") ||
			strings.Contains(iconLower, "ability_demonhunter_havoc") {
			return "DPS"
		}
		// If we can't determine spec, return Unknown for better debugging
		return "Unknown"

	case "warrior":
		// Protection = Tank, Arms/Fury = DPS
		if strings.Contains(iconLower, "protection") ||
			strings.Contains(iconLower, "warrior_protection") ||
			strings.Contains(iconLower, "ability_warrior_defensivestance") {
			return "Tank"
		}
		if strings.Contains(iconLower, "arms") || strings.Contains(iconLower, "fury") ||
			strings.Contains(iconLower, "warrior_arms") || strings.Contains(iconLower, "warrior_fury") {
			return "DPS"
		}
		return "DPS" // Default to DPS for warriors if unclear

	// Hybrid classes (need spec detection)
	case "druid":
		if strings.Contains(iconLower, "guardian") ||
			strings.Contains(iconLower, "druid_guardian") ||
			strings.Contains(iconLower, "ability_druid_maul") {
			return "Tank"
		} else if strings.Contains(iconLower, "restoration") ||
			strings.Contains(iconLower, "druid_restoration") ||
			strings.Contains(iconLower, "spell_nature_healingtouch") {
			return "Healer"
		}
		return "DPS" // Balance or Feral

	case "monk":
		if strings.Contains(iconLower, "brewmaster") ||
			strings.Contains(iconLower, "monk_brewmaster") ||
			strings.Contains(iconLower, "spell_monk_brewmaster_spec") {
			return "Tank"
		} else if strings.Contains(iconLower, "mistweaver") ||
			strings.Contains(iconLower, "monk_mistweaver") ||
			strings.Contains(iconLower, "spell_monk_mistweaver_spec") {
			return "Healer"
		}
		return "DPS" // Windwalker

	case "paladin":
		if strings.Contains(iconLower, "protection") ||
			strings.Contains(iconLower, "paladin_protection") ||
			strings.Contains(iconLower, "ability_paladin_shieldoftemplar") {
			return "Tank"
		} else if strings.Contains(iconLower, "holy") ||
			strings.Contains(iconLower, "paladin_holy") ||
			strings.Contains(iconLower, "spell_holy_holybolt") {
			return "Healer"
		}
		return "DPS" // Retribution

	case "priest":
		// Shadow = DPS, Holy/Discipline = Healer
		if strings.Contains(iconLower, "shadow") ||
			strings.Contains(iconLower, "priest_shadow") ||
			strings.Contains(iconLower, "spell_shadow_shadowwordpain") {
			return "DPS"
		}
		return "Healer" // Holy or Discipline

	case "shaman":
		// Restoration = Healer, Enhancement/Elemental = DPS
		if strings.Contains(iconLower, "restoration") ||
			strings.Contains(iconLower, "shaman_restoration") ||
			strings.Contains(iconLower, "spell_nature_healingwavegreater") {
			return "Healer"
		}
		return "DPS" // Enhancement or Elemental

	case "evoker":
		// Preservation = Healer, Devastation/Augmentation = DPS
		if strings.Contains(iconLower, "preservation") ||
			strings.Contains(iconLower, "evoker_preservation") {
			return "Healer"
		}
		return "DPS"

	default:
		if debug {
			color.Gray.Printf("   ❓ Unknown class: %s\n", class)
		}
		return "Unknown"
	}
}

// applyFilters applies all filtering criteria to the player list
func applyFilters(players []*EnhancedPlayer, classFilter, roleFilter, searchFilter string) []*EnhancedPlayer {
	filtered := make([]*EnhancedPlayer, 0, len(players))

	for _, player := range players {
		// Class filter
		if classFilter != "" && !strings.EqualFold(player.Class, classFilter) {
			continue
		}

		// Role filter
		if roleFilter != "" && !strings.EqualFold(player.Role, roleFilter) {
			continue
		}

		// Search filter (checks name)
		if searchFilter != "" && !strings.Contains(strings.ToLower(player.Name), strings.ToLower(searchFilter)) {
			continue
		}

		filtered = append(filtered, player)
	}

	return filtered
}

// displayPlayersInTerminal shows the player list in a beautiful terminal format
func displayPlayersInTerminal(players []*EnhancedPlayer, reportCode, fightIDStr, classFilter, roleFilter, searchFilter string, noColor bool) {
	// Build title with applied filters
	title := fmt.Sprintf("PLAYERS IN REPORT %s", reportCode)
	if fightIDStr != "" {
		title += fmt.Sprintf(" (Fight %s)", fightIDStr)
	}
	if classFilter != "" || roleFilter != "" || searchFilter != "" {
		filters := make([]string, 0, 3)
		if classFilter != "" {
			filters = append(filters, fmt.Sprintf("Class: %s", classFilter))
		}
		if roleFilter != "" {
			filters = append(filters, fmt.Sprintf("Role: %s", roleFilter))
		}
		if searchFilter != "" {
			filters = append(filters, fmt.Sprintf("Search: %s", searchFilter))
		}
		title += fmt.Sprintf(" (%s)", strings.Join(filters, ", "))
	}

	// Header
	fmt.Printf("\n👥 %s 👥\n", color.Cyan.Sprint(title))
	fmt.Printf("%s\n\n", color.Gray.Sprintf("Found %d players:", len(players)))

	// Table headers
	if noColor {
		fmt.Printf("%-3s %-20s %-12s %-8s %-20s\n", "#", "NAME", "CLASS", "ROLE", "SERVER")
		fmt.Printf("%-3s %-20s %-12s %-8s %-20s\n", "---", "--------------------", "------------", "--------", "--------------------")
	} else {
		color.White.Printf("%-3s %-20s %-12s %-8s %-20s\n", "#", "NAME", "CLASS", "ROLE", "SERVER")
		color.Gray.Printf("%-3s %-20s %-12s %-8s %-20s\n", "---", "--------------------", "------------", "--------", "--------------------")
	}

	// Player list with class and role colors
	for i, player := range players {
		if noColor {
			fmt.Printf("%-3d %-20s %-12s %-8s %-20s\n",
				i+1,
				player.Name,
				player.Class,
				player.Role,
				player.Server)
		} else {
			classColor := getClassColor(player.Class)
			roleColor := getRoleColor(player.Role)
			fmt.Printf("%-3d %-20s %s %s %-20s\n",
				i+1,
				player.Name,
				classColor.Sprint(fmt.Sprintf("%-12s", player.Class)),
				roleColor.Sprint(fmt.Sprintf("%-8s", player.Role)),
				player.Server)
		}
	}

	// Footer with usage hints
	if !noColor {
		fmt.Printf("\n%s\n", color.Green.Sprint("✅ Use these exact names with --player flag"))
		if len(players) > 0 {
			fmt.Printf("%s\n", color.Yellow.Sprintf("Example: wclogs damage %s 5 --player \"%s\"", reportCode, players[0].Name))
		}
	}

	// Show role/class distribution
	showDistribution(players, noColor)
}

// showDistribution displays a summary of class and role distribution
func showDistribution(players []*EnhancedPlayer, noColor bool) {
	classCounts := make(map[string]int)
	roleCounts := make(map[string]int)

	for _, player := range players {
		classCounts[player.Class]++
		roleCounts[player.Role]++
	}

	if !noColor {
		fmt.Printf("\n📊 %s\n", color.Cyan.Sprint("COMPOSITION SUMMARY"))
	} else {
		fmt.Printf("\nCOMPOSITION SUMMARY\n")
	}

	// Show role distribution
	fmt.Printf("Roles: ")
	roleOrder := []string{"Tank", "Healer", "DPS", "Unknown"}
	roleParts := make([]string, 0, len(roleOrder))
	for _, role := range roleOrder {
		if count, exists := roleCounts[role]; exists && count > 0 {
			if noColor {
				roleParts = append(roleParts, fmt.Sprintf("%s: %d", role, count))
			} else {
				roleColor := getRoleColor(role)
				roleParts = append(roleParts, roleColor.Sprint(fmt.Sprintf("%s: %d", role, count)))
			}
		}
	}
	fmt.Printf("%s\n", strings.Join(roleParts, ", "))

	// Show most common classes
	type classCount struct {
		class string
		count int
	}
	var sortedClasses []classCount
	for class, count := range classCounts {
		sortedClasses = append(sortedClasses, classCount{class, count})
	}
	sort.Slice(sortedClasses, func(i, j int) bool {
		return sortedClasses[i].count > sortedClasses[j].count
	})

	fmt.Printf("Top Classes: ")
	classParts := make([]string, 0, 3)
	for i, cc := range sortedClasses {
		if i >= 3 { // Show top 3 classes
			break
		}
		if noColor {
			classParts = append(classParts, fmt.Sprintf("%s: %d", cc.class, cc.count))
		} else {
			classColor := getClassColor(cc.class)
			classParts = append(classParts, classColor.Sprint(fmt.Sprintf("%s: %d", cc.class, cc.count)))
		}
	}
	fmt.Printf("%s\n", strings.Join(classParts, ", "))
}

// outputPlayersToFile saves the player list to a file
func outputPlayersToFile(players []*EnhancedPlayer, reportCode, outputPath string, verbose bool) error {
	if verbose {
		color.Cyan.Printf("💾 Saving %d players to file: %s\n", len(players), outputPath)
	}

	// Create output data structure
	outputData := &output.PlayersOutputData{
		ReportCode: reportCode,
		Count:      len(players),
		Players:    make([]*models.PlayerInfo, len(players)),
	}

	// Convert enhanced players to basic player info for output
	for i, player := range players {
		outputData.Players[i] = &models.PlayerInfo{
			ID:     player.ID,
			Name:   player.Name,
			Class:  player.Class,
			Server: player.Server,
			Icon:   player.Icon,
		}
	}

	return output.HandlePlayersOutput(outputData, outputPath, verbose)
}

// getClassColor returns the appropriate color function for each class using official WoW RGB values
func getClassColor(class string) color.RGBColor {
	switch class {
	case "Death Knight", "DeathKnight":
		return color.RGB(196, 30, 58) // Death Knight: Official Red
	case "Demon Hunter", "DemonHunter":
		return color.RGB(163, 48, 201) // Demon Hunter: Official Purple
	case "Druid":
		return color.RGB(255, 124, 10) // Druid: Official Orange
	case "Hunter":
		return color.RGB(170, 211, 114) // Hunter: Official Green
	case "Mage":
		return color.RGB(63, 199, 235) // Mage: Official Light Blue
	case "Monk":
		return color.RGB(0, 255, 152) // Monk: Official Jade Green
	case "Paladin":
		return color.RGB(244, 140, 186) // Paladin: Official Pink
	case "Priest":
		return color.RGB(255, 255, 255) // Priest: Official White
	case "Rogue":
		return color.RGB(255, 244, 104) // Rogue: Official Yellow
	case "Shaman":
		return color.RGB(0, 112, 221) // Shaman: Official Blue
	case "Warlock":
		return color.RGB(135, 136, 238) // Warlock: Official Purple
	case "Warrior":
		return color.RGB(198, 155, 109) // Warrior: Official Brown
	case "Evoker":
		return color.RGB(51, 147, 127) // Evoker: Teal-Green (estimated)
	default:
		return color.RGB(255, 255, 255) // Default: White
	}
}

// getRoleColor returns the appropriate color function for each role
func getRoleColor(role string) color.RGBColor {
	switch role {
	case "Tank":
		return color.RGB(74, 144, 226) // Tank: Blue
	case "Healer":
		return color.RGB(80, 200, 120) // Healer: Emerald Green
	case "DPS":
		return color.RGB(255, 107, 107) // DPS: Coral Red
	default:
		return color.RGB(153, 153, 153) // Unknown: Gray
	}
}
