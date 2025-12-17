package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search issues by text",
	Long: `Search for issues in your Linear workspace by text.

Examples:
  linear search "bug fix"
  linear search authentication
  linear search "login error" -d   # Include descriptions in output`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		client, err := getLinearClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		issues, err := client.SearchIssues(query)
		if err != nil {
			fmt.Printf("Error searching issues: %v\n", err)
			os.Exit(1)
		}

		displaySearchResults(issues, query)
	},
}

func displaySearchResults(issues []Issue, query string) {
	if len(issues) == 0 {
		fmt.Printf("No issues found matching '%s'.\n", query)
		return
	}

	fmt.Printf("Found %d issues matching '%s':\n\n", len(issues), query)

	for i, issue := range issues {
		fmt.Printf("%d. [%s] %s\n", i+1, issue.Identifier, issue.Title)
		fmt.Printf("   Team: %s | State: %s", issue.Team.Name, issue.State.Name)
		if issue.Assignee.Name != "" {
			fmt.Printf(" | Assignee: %s", issue.Assignee.Name)
		}
		if issue.Priority > 0 {
			fmt.Printf(" | Priority: %d", issue.Priority)
		}

		if showDescription && issue.Description != "" {
			fmt.Printf("\n   Description: %s", issue.Description)
		}
		fmt.Printf("\n\n")
	}
}

func init() {
	searchCmd.Flags().BoolVarP(&showDescription, "description", "d", false, "Show issue descriptions")
	rootCmd.AddCommand(searchCmd)
}
