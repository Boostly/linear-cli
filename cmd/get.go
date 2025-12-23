package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [issue-identifier]",
	Short: "Get details of a specific issue",
	Long:  `Fetch and display details of a Linear issue by its identifier (e.g., ENG-123).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		client, err := getLinearClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		issue, err := client.GetIssue(identifier)
		if err != nil {
			fmt.Printf("Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s: %s\n", issue.Identifier, issue.Title)
		fmt.Printf("State: %s | Team: %s", issue.State.Name, issue.Team.Name)
		if issue.Priority > 0 {
			fmt.Printf(" | Priority: %d", issue.Priority)
		}
		fmt.Println()

		if issue.Assignee.Name != "" {
			fmt.Printf("Assignee: %s\n", issue.Assignee.Name)
		}

		if issue.BranchName != "" {
			fmt.Printf("Branch: %s\n", issue.BranchName)
		}

		if issue.Parent.ID != "" {
			fmt.Printf("Parent: %s - %s\n", issue.Parent.Identifier, issue.Parent.Title)
		}

		fmt.Printf("URL: %s\n", issue.URL)

		if issue.Description != "" {
			fmt.Printf("\nDescription:\n%s\n", issue.Description)
		}

		if len(issue.Children.Nodes) > 0 {
			fmt.Printf("\nSub-issues (%d):\n", len(issue.Children.Nodes))
			for _, child := range issue.Children.Nodes {
				fmt.Printf("  - %s: %s [%s]\n", child.Identifier, child.Title, child.State.Name)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
