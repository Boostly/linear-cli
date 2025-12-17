package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/Boostly/linear-cli/cmd"
)

func main() {
	// Load environment variables from .env file
	// Try local .env first, then ~/.config/linear/.env
	err := godotenv.Load()
	if err != nil {
		home, _ := os.UserHomeDir()
		if home != "" {
			configPath := home + "/.config/linear/.env"
			err = godotenv.Load(configPath)
		}
	}

	// Check if we have the required API key
	if os.Getenv("LINEAR_API_KEY") == "" {
		fmt.Println("Error: LINEAR_API_KEY environment variable not set")
		fmt.Println("Please create a .env file with your Linear API key:")
		fmt.Println("LINEAR_API_KEY=your_api_key_here")
		os.Exit(1)
	}

	// Execute cobra CLI
	cmd.Execute()
}
