package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"wclogs-cli/api"
	"wclogs-cli/auth"
	"wclogs-cli/display"
	"wclogs-cli/models"
)

func main() {
	// Get credentials from environment variables
	clientID := os.Getenv("WCLOGS_CLIENT_ID")
	clientSecret := os.Getenv("WCLOGS_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set WCLOGS_CLIENT_ID and WCLOGS_CLIENT_SECRET environment variables")
	}

	// Test with a real report
	testReportCode := "Hw9TZc2WyrVKJLCa"
	testFightID := 99

	fmt.Printf("🔍 Testing Warcraft Logs API with report %s, fight %d\n", testReportCode, testFightID)

	// Create auth client
	fmt.Println("🔐 Setting up authentication...")
	authClient := auth.NewClient(clientID, clientSecret)

	// Create API client
	fmt.Println("📡 Creating API client...")
	apiClient := api.NewClient(authClient)

	// Validate parameters
	fmt.Println("✅ Validating parameters...")
	if err := api.ValidateQueryVariables(testReportCode, testFightID); err != nil {
		log.Fatalf("Invalid parameters: %v", err)
	}

	// Execute query
	fmt.Println("🚀 Executing GraphQL query...")
	response, err := apiClient.Query(
		api.DamageTableQuery,
		map[string]any{
			"code":    testReportCode,
			"fightID": testFightID,
		},
	)

	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	// Check if we got data
	if response.Data == nil || response.Data.ReportData == nil || response.Data.ReportData.Report == nil {
		log.Fatal("No report data found")
	}

	rawTable := response.Data.ReportData.Report.Table
	if len(rawTable) == 0 {
		log.Fatal("No table data found")
	}

	fmt.Println("✅ Got table data! Parsing...")

	// Parse the table data
	tableData, err := models.ParseTableData(rawTable)
	if err != nil {
		log.Fatalf("Failed to parse table data: %v", err)
	}

	fmt.Printf("📊 Found %d players in the table\n", len(tableData.Entries))

	// Convert to Player objects
	players := models.GetPlayersFromTable(tableData)

	// Display the damage table with enhanced styling
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🗡️  WARCRAFT LOGS DAMAGE TABLE")
	fmt.Printf("📊 Report: %s | ⚔️  Fight: %d | 👥 Players: %d\n", testReportCode, testFightID, len(players))
	fmt.Println(strings.Repeat("=", 60))

	options := display.DefaultTableOptions()
	options.TopN = 0         // Show all players
	options.UseColors = true // Enable beautiful colors

	display.DisplayDamageTable(players, options)

	fmt.Println("\n🎉 Success! Damage table displayed with colors and formatting!")
}
