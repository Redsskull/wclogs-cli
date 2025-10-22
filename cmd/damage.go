package cmd

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
	"wclogs-cli/models"
	"wclogs-cli/output"
)

// damageCmd represents the damage command
var damageCmd = &cobra.Command{
	Use:   "damage [report-code] [fight-id]",
	Short: "🗡️  Show damage table for a fight",
	Long: color.HiYellowString(`
🗡️  DAMAGE TABLE COMMAND

Display damage done by all players in a specific fight.

The report-code is found in Warcraft Logs URLs:
  https://www.warcraftlogs.com/reports/ABC123XYZ
  Report code: ABC123XYZ

Fight ID is the encounter number (usually 1-20).

Examples:
  wclogs damage ABC123XYZ 5           # Show damage for fight 5
  wclogs damage ABC123XYZ 5 --top 10  # Show top 10 players only
  wclogs damage ABC123XYZ 5 --no-color # Disable colors
`) + "\n",
	Args: cobra.ExactArgs(2), // Must have exactly 2 arguments
	RunE: func(cmd *cobra.Command, args []string) error {
		reportCode := args[0]
		fightID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("fight-id must be a number, got: %s", args[1])
		}

		// Get flag values
		topN, _ := cmd.Flags().GetInt("top")
		noColor, _ := cmd.Flags().GetBool("no-color")
		verbose, _ := cmd.Flags().GetBool("verbose")
		outputPath, _ := cmd.Flags().GetString("output")

		return showDamageTable(reportCode, fightID, topN, noColor, verbose, outputPath)
	},
}

func init() {
	// Register this command with the root command
	rootCmd.AddCommand(damageCmd)

	// Add flags specific to the damage command
	damageCmd.Flags().BoolP("no-color", "n", false, "Disable color output")

	// Note: --top and --verbose are inherited from root command's PersistentFlags
}

// showDamageTable executes the damage table logic
func showDamageTable(reportCode string, fightID int, topN int, noColor bool, verbose bool, outputPath string) error {
	if verbose {
		color.HiBlue("🔍 Fetching damage data for report %s, fight %d", reportCode, fightID)
	}

	// Load credentials from config file (we know it exists thanks to root.go check)
	if verbose {
		color.HiBlue("🔐 Loading configuration...")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	clientID := cfg.ClientID
	clientSecret := cfg.ClientSecret

	// Create auth client
	if verbose {
		color.HiBlue("🔐 Setting up authentication...")
	}
	authClient := auth.NewClient(clientID, clientSecret)

	// Create API client
	if verbose {
		color.HiBlue("📡 Creating API client...")
	}
	apiClient := api.NewClient(authClient)

	// Validate parameters
	if verbose {
		color.HiBlue("✅ Validating parameters...")
	}
	if err := api.ValidateQueryVariables(reportCode, fightID); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Execute query
	if verbose {
		color.HiBlue("🚀 Executing GraphQL query...")
	}
	response, err := apiClient.Query(
		api.DamageTableQuery,
		map[string]any{
			"code":    reportCode,
			"fightID": fightID,
		},
	)

	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Check if we got data
	if response.Data == nil || response.Data.ReportData == nil || response.Data.ReportData.Report == nil {
		return fmt.Errorf("no report data found for code: %s", reportCode)
	}

	rawTable := response.Data.ReportData.Report.Table
	if len(rawTable) == 0 {
		return fmt.Errorf("no table data found for fight %d in report %s", fightID, reportCode)
	}

	if verbose {
		color.HiBlue("✅ Got table data! Parsing...")
	}

	// Parse the table data
	tableData, err := models.ParseTableData(rawTable)
	if err != nil {
		return fmt.Errorf("failed to parse table data: %w", err)
	}

	if verbose {
		color.HiBlue("📊 Found %d players in the table", len(tableData.Entries))
	}

	// Convert to Player objects
	players := models.GetPlayersFromTable(tableData)

	// Calculate total damage
	var totalDamage int64
	for _, player := range players {
		totalDamage += int64(player.Total)
	}

	// Create structured output data
	outputData := &output.OutputData{
		Players:    players,
		ReportCode: reportCode,
		FightID:    fightID,
		Title:      "Damage Table",
		Total:      totalDamage,
	}

	// Handle output (terminal display or file save)
	return output.HandleOutput(outputData, outputPath, topN, noColor, verbose)
}
