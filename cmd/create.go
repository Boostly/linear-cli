package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	createTitle       string
	createTeam        string
	createDescription string
	createPriority    int
	createState       string
	createAssignee    string
	createParent      string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long: `Create a new issue in your Linear workspace.

Requires at minimum a title and team. Team can be specified by name or key.

Examples:
  linear-cli create -t "Bug fix" -T "Engineering"
  linear-cli create -t "New feature" -T ENG -d "Description here" -p 2
  linear-cli create -t "Task" -T "Product" -s "Backlog" -a "john@example.com"
  linear-cli create -t "Sub-task" -T ENG -P ENG-123  # Create sub-issue`,
	Run: func(cmd *cobra.Command, args []string) {
		if createTitle == "" {
			fmt.Println("Error: title is required (-t, --title)")
			os.Exit(1)
		}
		if createTeam == "" {
			fmt.Println("Error: team is required (-T, --team)")
			os.Exit(1)
		}

		client, err := getLinearClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Resolve team name/key to team ID
		teams, err := client.GetTeams()
		if err != nil {
			fmt.Printf("Error fetching teams: %v\n", err)
			os.Exit(1)
		}

		var teamID string
		for _, team := range teams {
			if strings.EqualFold(team.Name, createTeam) || strings.EqualFold(team.Key, createTeam) {
				teamID = team.ID
				break
			}
		}

		if teamID == "" {
			fmt.Printf("Error: team '%s' not found. Available teams:\n", createTeam)
			for _, team := range teams {
				fmt.Printf("  - %s (%s)\n", team.Name, team.Key)
			}
			os.Exit(1)
		}

		// Build create input
		input := IssueCreateInput{
			Title:       createTitle,
			TeamID:      teamID,
			Description: createDescription,
			Priority:    createPriority,
		}

		// Resolve state if provided
		if createState != "" {
			states, err := client.GetWorkflowStates(teamID)
			if err != nil {
				fmt.Printf("Error fetching workflow states: %v\n", err)
				os.Exit(1)
			}

			var stateID string
			for _, state := range states {
				if strings.EqualFold(state.Name, createState) {
					stateID = state.ID
					break
				}
			}

			if stateID == "" {
				fmt.Printf("Error: state '%s' not found. Available states:\n", createState)
				for _, state := range states {
					fmt.Printf("  - %s\n", state.Name)
				}
				os.Exit(1)
			}
			input.StateID = stateID
		}

		// Resolve assignee if provided
		if createAssignee != "" {
			members, err := client.GetTeamMembers(teamID)
			if err != nil {
				fmt.Printf("Error fetching team members: %v\n", err)
				os.Exit(1)
			}

			var assigneeID string
			for _, member := range members {
				if strings.EqualFold(member.Name, createAssignee) || strings.EqualFold(member.Email, createAssignee) {
					assigneeID = member.ID
					break
				}
			}

			if assigneeID == "" {
				fmt.Printf("Error: assignee '%s' not found. Available team members:\n", createAssignee)
				for _, member := range members {
					fmt.Printf("  - %s (%s)\n", member.Name, member.Email)
				}
				os.Exit(1)
			}
			input.AssigneeID = assigneeID
		}

		// Resolve parent issue if provided
		if createParent != "" {
			parentIssue, err := client.GetIssue(createParent)
			if err != nil {
				fmt.Printf("Error: parent issue '%s' not found: %v\n", createParent, err)
				os.Exit(1)
			}
			input.ParentID = parentIssue.ID
		}

		fmt.Println("Creating issue...")

		issue, err := client.CreateIssue(input)
		if err != nil {
			fmt.Printf("Error creating issue: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nCreated: %s - %s\n", issue.Identifier, issue.Title)
		fmt.Printf("  Team: %s | State: %s\n", issue.Team.Name, issue.State.Name)
		fmt.Printf("  URL: %s\n", issue.URL)
	},
}

func init() {
	createCmd.Flags().StringVarP(&createTitle, "title", "t", "", "Issue title (required)")
	createCmd.Flags().StringVarP(&createTeam, "team", "T", "", "Team name or key (required)")
	createCmd.Flags().StringVarP(&createDescription, "description", "d", "", "Issue description")
	createCmd.Flags().IntVarP(&createPriority, "priority", "p", 0, "Priority (0=None, 1=Urgent, 2=High, 3=Medium, 4=Low)")
	createCmd.Flags().StringVarP(&createState, "state", "s", "", "Workflow state name")
	createCmd.Flags().StringVarP(&createAssignee, "assignee", "a", "", "Assignee name or email")
	createCmd.Flags().StringVarP(&createParent, "parent", "P", "", "Parent issue identifier for sub-issues (e.g., ENG-123)")
}
