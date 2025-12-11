package cmd

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Issue struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       struct {
		Name string `json:"name"`
	} `json:"state"`
	Team struct {
		Name string `json:"name"`
	} `json:"team"`
	Assignee struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"assignee"`
	Priority int `json:"priority"`
}

type IssuesResponse struct {
	Data struct {
		Issues struct {
			Nodes []Issue `json:"nodes"`
		} `json:"issues"`
	} `json:"data"`
}

type ViewerResponse struct {
	Data struct {
		Viewer User `json:"viewer"`
	} `json:"data"`
}

type GraphQLRequest struct {
	Query string `json:"query"`
}

// GraphQLRequestWithVariables is used for mutations with variables
type GraphQLRequestWithVariables struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// Team represents a Linear team
type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

// TeamsResponse for GraphQL teams query
type TeamsResponse struct {
	Data struct {
		Teams struct {
			Nodes []Team `json:"nodes"`
		} `json:"teams"`
	} `json:"data"`
}

// WorkflowState represents a workflow state (e.g., "Backlog", "In Progress")
type WorkflowState struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// WorkflowStatesResponse for GraphQL workflowStates query
type WorkflowStatesResponse struct {
	Data struct {
		WorkflowStates struct {
			Nodes []WorkflowState `json:"nodes"`
		} `json:"workflowStates"`
	} `json:"data"`
}

// TeamMembersResponse for fetching team members
type TeamMembersResponse struct {
	Data struct {
		Team struct {
			Members struct {
				Nodes []User `json:"nodes"`
			} `json:"members"`
		} `json:"team"`
	} `json:"data"`
}

// IssueCreateInput represents the input for creating an issue
type IssueCreateInput struct {
	Title       string
	TeamID      string
	Description string
	Priority    int
	StateID     string
	AssigneeID  string
}

// IssueCreateResponse for the mutation response
type IssueCreateResponse struct {
	Data struct {
		IssueCreate struct {
			Success bool  `json:"success"`
			Issue   Issue `json:"issue"`
		} `json:"issueCreate"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// CreatedIssue contains additional fields returned when creating an issue
type CreatedIssue struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
	State      struct {
		Name string `json:"name"`
	} `json:"state"`
	Team struct {
		Name string `json:"name"`
	} `json:"team"`
}

// IssueCreateFullResponse for the mutation response with full issue details
type IssueCreateFullResponse struct {
	Data struct {
		IssueCreate struct {
			Success bool         `json:"success"`
			Issue   CreatedIssue `json:"issue"`
		} `json:"issueCreate"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}
