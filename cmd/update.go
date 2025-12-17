package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	updateIssueID     string
	updateTitle       string
	updateDescription string
	updatePriority    int
	updateState       string
	updateAssignee    string
	updatePrioritySet bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing issue",
	Long: `Update an existing issue in your Linear workspace.

Requires an issue identifier (e.g., ENG-123) or issue ID.

Examples:
  linear-tui update -i ENG-123 -t "New title"
  linear-tui update -i ENG-123 -s "In Progress"
  linear-tui update -i ENG-123 -a "john@example.com" -p 2
  linear-tui update -i ENG-123 -d "Updated description"`,
	Run: func(cmd *cobra.Command, args []string) {
		if updateIssueID == "" {
			fmt.Println("Error: issue identifier is required (-i, --issue)")
			os.Exit(1)
		}

		// Check if any update flags were provided
		titleChanged := cmd.Flags().Changed("title")
		descChanged := cmd.Flags().Changed("description")
		stateChanged := cmd.Flags().Changed("state")
		assigneeChanged := cmd.Flags().Changed("assignee")
		priorityChanged := cmd.Flags().Changed("priority")

		if !titleChanged && !descChanged && !stateChanged && !assigneeChanged && !priorityChanged {
			fmt.Println("Error: at least one field to update is required")
			fmt.Println("Use -t (title), -d (description), -s (state), -a (assignee), or -p (priority)")
			os.Exit(1)
		}

		client, err := getLinearClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Fetch the issue to get its ID and team ID
		issue, err := client.GetIssue(updateIssueID)
		if err != nil {
			fmt.Printf("Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		input := IssueUpdateInput{}

		if titleChanged {
			input.Title = updateTitle
		}
		if descChanged {
			input.Description = updateDescription
		}
		if priorityChanged {
			input.Priority = &updatePriority
		}

		// Resolve state name to ID if provided
		if stateChanged {
			states, err := client.GetWorkflowStates(issue.Team.ID)
			if err != nil {
				fmt.Printf("Error fetching workflow states: %v\n", err)
				os.Exit(1)
			}

			var stateID string
			for _, state := range states {
				if strings.EqualFold(state.Name, updateState) {
					stateID = state.ID
					break
				}
			}

			if stateID == "" {
				fmt.Printf("Error: state '%s' not found. Available states:\n", updateState)
				for _, state := range states {
					fmt.Printf("  - %s\n", state.Name)
				}
				os.Exit(1)
			}
			input.StateID = stateID
		}

		// Resolve assignee name/email to ID if provided
		if assigneeChanged {
			members, err := client.GetTeamMembers(issue.Team.ID)
			if err != nil {
				fmt.Printf("Error fetching team members: %v\n", err)
				os.Exit(1)
			}

			var assigneeID string
			for _, member := range members {
				if strings.EqualFold(member.Name, updateAssignee) || strings.EqualFold(member.Email, updateAssignee) {
					assigneeID = member.ID
					break
				}
			}

			if assigneeID == "" {
				fmt.Printf("Error: assignee '%s' not found. Available team members:\n", updateAssignee)
				for _, member := range members {
					fmt.Printf("  - %s (%s)\n", member.Name, member.Email)
				}
				os.Exit(1)
			}
			input.AssigneeID = assigneeID
		}

		fmt.Printf("Updating issue %s...\n", issue.Identifier)

		updated, err := client.UpdateIssue(issue.ID, input)
		if err != nil {
			fmt.Printf("Error updating issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nUpdated: %s - %s\n", updated.Identifier, updated.Title)
		fmt.Printf("  Team: %s | State: %s | Priority: %d\n",
			updated.Team.Name, updated.State.Name, updated.Priority)
		if updated.Assignee.Name != "" {
			fmt.Printf("  Assignee: %s\n", updated.Assignee.Name)
		}
	},
}

func init() {
	updateCmd.Flags().StringVarP(&updateIssueID, "issue", "i", "", "Issue identifier or ID (required)")
	updateCmd.Flags().StringVarP(&updateTitle, "title", "t", "", "New issue title")
	updateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "New issue description")
	updateCmd.Flags().IntVarP(&updatePriority, "priority", "p", 0, "Priority (0=None, 1=Urgent, 2=High, 3=Medium, 4=Low)")
	updateCmd.Flags().StringVarP(&updateState, "state", "s", "", "New workflow state name")
	updateCmd.Flags().StringVarP(&updateAssignee, "assignee", "a", "", "New assignee name or email")
}
