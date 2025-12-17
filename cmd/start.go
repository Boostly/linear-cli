package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	startNoBranch bool
	startBranch   string
)

var startCmd = &cobra.Command{
	Use:   "start [issue-identifier]",
	Short: "Start working on an issue",
	Long: `Start working on a Linear issue. This command:

1. Assigns the issue to you (if not already assigned)
2. Sets the status to "In Progress"
3. Creates and checks out a git branch

Examples:
  linear start ENG-123
  linear start ENG-123 --no-branch    # Skip branch creation
  linear start ENG-123 -b "my-branch" # Custom branch name`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		client, err := getLinearClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Get current user
		user, err := client.GetCurrentUser()
		if err != nil {
			fmt.Printf("Error getting current user: %v\n", err)
			os.Exit(1)
		}

		// Fetch issue with branch info
		issue, err := client.GetIssue(identifier)
		if err != nil {
			fmt.Printf("Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Starting work on %s - %s\n\n", issue.Identifier, issue.Title)

		// Build update input
		updateInput := IssueUpdateInput{}
		needsUpdate := false

		// Step 1: Assign to current user if not assigned
		if issue.Assignee.ID == "" || issue.Assignee.ID != user.ID {
			fmt.Printf("Assigning to %s...\n", user.Name)
			updateInput.AssigneeID = user.ID
			needsUpdate = true
		} else {
			fmt.Printf("Already assigned to %s\n", user.Name)
		}

		// Step 2: Set status to "In Progress" if not already started
		if issue.State.Type != "started" {
			inProgressState, err := client.GetInProgressState(issue.Team.ID)
			if err != nil {
				fmt.Printf("Warning: %v\n", err)
			} else {
				fmt.Printf("Setting status to '%s'...\n", inProgressState.Name)
				updateInput.StateID = inProgressState.ID
				needsUpdate = true
			}
		} else {
			fmt.Printf("Status already '%s'\n", issue.State.Name)
		}

		// Apply updates if needed
		if needsUpdate {
			_, err = client.UpdateIssue(issue.ID, updateInput)
			if err != nil {
				fmt.Printf("Error updating issue: %v\n", err)
				os.Exit(1)
			}
		}

		// Step 3: Create and checkout git branch
		if !startNoBranch {
			branchName := startBranch
			if branchName == "" {
				branchName = issue.BranchName
			}
			if branchName == "" {
				// Fallback: generate from identifier and title
				branchName = generateBranchName(issue.Identifier, issue.Title)
			}

			fmt.Printf("\nCreating branch '%s'...\n", branchName)

			// Check if we're in a git repo
			if err := exec.Command("git", "rev-parse", "--git-dir").Run(); err != nil {
				fmt.Println("Warning: not in a git repository, skipping branch creation")
			} else {
				// Create and checkout branch
				gitCmd := exec.Command("git", "checkout", "-b", branchName)
				gitCmd.Stdout = os.Stdout
				gitCmd.Stderr = os.Stderr
				if err := gitCmd.Run(); err != nil {
					// Branch might already exist, try just checkout
					gitCmd = exec.Command("git", "checkout", branchName)
					gitCmd.Stdout = os.Stdout
					gitCmd.Stderr = os.Stderr
					if err := gitCmd.Run(); err != nil {
						fmt.Printf("Error creating/checking out branch: %v\n", err)
					}
				}
			}
		}

		fmt.Printf("\nReady to work on %s!\n", issue.Identifier)
		if issue.URL != "" {
			fmt.Printf("View issue: %s\n", issue.URL)
		}
	},
}

func generateBranchName(identifier, title string) string {
	// Convert to lowercase and replace spaces with hyphens
	branch := strings.ToLower(identifier + "-" + title)
	// Remove special characters
	branch = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		if r == ' ' {
			return '-'
		}
		return -1
	}, branch)
	// Limit length
	if len(branch) > 50 {
		branch = branch[:50]
	}
	// Remove trailing hyphens
	branch = strings.TrimRight(branch, "-")
	return branch
}

func init() {
	startCmd.Flags().BoolVar(&startNoBranch, "no-branch", false, "Skip git branch creation")
	startCmd.Flags().StringVarP(&startBranch, "branch", "b", "", "Custom branch name (overrides Linear's suggestion)")
	rootCmd.AddCommand(startCmd)
}
