package main

import (
	"fmt"
	"log"
	"os"

	"wclogs-cli/api"
	"wclogs-cli/auth"
)

func main() {
	// Get credentials
	clientID := os.Getenv("WCL_CLIENT_ID")
	clientSecret := os.Getenv("WCL_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set WCL_CLIENT_ID and WCL_CLIENT_SECRET environment variables")
	}

	// Create auth client
	authClient := auth.NewClient(clientID, clientSecret)

	// Create API client
	apiClient := api.NewClient(authClient)

	// Test 1: Simple connection test
	fmt.Println("🚀 Testing GraphQL API connection...")

	query1 := `query { reportData { __typename } }`
	resp1, err := apiClient.Query(query1, nil)
	if err != nil {
		log.Fatalf("❌ Connection test failed: %v", err)
	}

	fmt.Printf("✅ Connection successful: %+v\n", resp1.Data)

	// Test 2: Get report info
	fmt.Println("\n📊 Testing report query...")

	query2 := `query($code: String!) {
		reportData {
			report(code: $code) {
				title
				startTime
				endTime
			}
		}
	}`

	variables := map[string]interface{}{
		"code": "AvxwLnYgm3qczQpG", // Your working report code
	}

	resp2, err := apiClient.Query(query2, variables)
	if err != nil {
		log.Fatalf("❌ Report query failed: %v", err)
	}

	fmt.Printf("✅ Report data: %+v\n", resp2.Data)

	fmt.Println("\n🎉 Both OAuth2 and GraphQL are working perfectly!")
}
