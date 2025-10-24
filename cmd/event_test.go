package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/config"
)

// testEventsCmd is a temporary command to explore the Events API
var testEventsCmd = &cobra.Command{
	Use:   "test-events [report-code] [fight-id]",
	Short: "🧪 Test Events API (temporary research command)",
	Long: color.HiMagentaString(`
🧪 TEST EVENTS API

This is a temporary command to research how the Events API works.
It will show us the raw JSON structure so we can understand how to implement
proper death and interrupt analysis.

Examples:
  wclogs test-events Hw9TZc2WyrVKJLCa 99
`) + "\n",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		reportCode := args[0]
		fightID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("fight-id must be a number, got: %s", args[1])
		}

		verbose, _ := cmd.Flags().GetBool("verbose")

		if verbose {
			color.HiBlue("🧪 Testing Events API with report %s, fight %d", reportCode, fightID)
		}

		// Auth setup
		cfg, err := config.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authClient := auth.NewClient(cfg.ClientID, cfg.ClientSecret)
		apiClient := api.NewClient(authClient)

		// Make the test query - using the simple TestEventsQuery
		if verbose {
			color.HiBlue("🚀 Executing simple Events API test query...")
		}

		// Use the simple test query instead of the death events query
		request := api.NewTestEventsRequest(reportCode, fightID)
		response, err := apiClient.Query(request.Query, request.Variables)
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}

		// Check if we have events data
		if response.Data != nil && response.Data.ReportData != nil &&
			response.Data.ReportData.Report != nil && response.Data.ReportData.Report.Events != nil {

			color.HiGreen("✅ Events API query successful!")
			fmt.Printf("\n🧪 RAW EVENTS JSON:\n")

			// Show the raw events data
			eventsJSON := string(response.Data.ReportData.Report.Events.Data)
			if len(eventsJSON) > 1000 {
				fmt.Printf("%s...\n[truncated - showing first 1000 chars of %d total]\n",
					eventsJSON[:1000], len(eventsJSON))
			} else {
				fmt.Printf("%s\n", eventsJSON)
			}
		} else {
			color.HiYellow("⚠️  No events data found")
		}

		// Show the response structure for debugging
		if verbose {
			fmt.Printf("\n🔍 RESPONSE STRUCTURE:\n")
			jsonData, err := json.MarshalIndent(response, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format JSON: %w", err)
			}
			fmt.Printf("%s\n", jsonData)
		}

		return nil
	},
}

func init() {
	testEventsCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.AddCommand(testEventsCmd)
}
