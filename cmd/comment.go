package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var commentMessage string

var commentCmd = &cobra.Command{
	Use:   "comment [issue-identifier]",
	Short: "Add a comment to an issue",
	Long: `Add a comment to a Linear issue.

Examples:
  linear comment ENG-123 -m "This is my comment"
  linear comment ENG-123    # Will prompt for comment text`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		client, err := getLinearClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Fetch issue to get its ID
		issue, err := client.GetIssue(identifier)
		if err != nil {
			fmt.Printf("Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		// Get comment message
		message := commentMessage
		if message == "" {
			fmt.Printf("Adding comment to %s - %s\n", issue.Identifier, issue.Title)
			fmt.Print("Enter comment (press Enter twice to finish):\n")

			reader := bufio.NewReader(os.Stdin)
			var lines []string
			emptyCount := 0
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				if strings.TrimSpace(line) == "" {
					emptyCount++
					if emptyCount >= 1 && len(lines) > 0 {
						break
					}
				} else {
					emptyCount = 0
					lines = append(lines, line)
				}
			}
			message = strings.Join(lines, "")
		}

		message = strings.TrimSpace(message)
		if message == "" {
			fmt.Println("Error: comment cannot be empty")
			os.Exit(1)
		}

		comment, err := client.CreateComment(issue.ID, message)
		if err != nil {
			fmt.Printf("Error creating comment: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Comment added to %s by %s\n", issue.Identifier, comment.User.Name)
	},
}

func init() {
	commentCmd.Flags().StringVarP(&commentMessage, "message", "m", "", "Comment message")
	rootCmd.AddCommand(commentCmd)
}
