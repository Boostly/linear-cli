package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var commentsCmd = &cobra.Command{
	Use:   "comments [issue-identifier]",
	Short: "View comments on an issue",
	Long:  `Fetch and display all comments on a Linear issue by its identifier (e.g., ENG-123).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		client, err := getLinearClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		comments, err := client.GetComments(identifier)
		if err != nil {
			fmt.Printf("Error fetching comments: %v\n", err)
			os.Exit(1)
		}

		if len(comments) == 0 {
			fmt.Printf("No comments on %s\n", identifier)
			return
		}

		fmt.Printf("Comments on %s (%d):\n\n", identifier, len(comments))

		for i, comment := range comments {
			timestamp := formatTimestamp(comment.CreatedAt)
			fmt.Printf("--- %s at %s ---\n", comment.User.Name, timestamp)
			fmt.Printf("%s\n", comment.Body)
			if i < len(comments)-1 {
				fmt.Println()
			}
		}
	},
}

func formatTimestamp(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	return t.Local().Format("Jan 2, 2006 3:04 PM")
}

func init() {
	rootCmd.AddCommand(commentsCmd)
}
