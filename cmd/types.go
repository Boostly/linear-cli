package cmd

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Issue struct {
	ID          string `json:"id"`
	Identifier  string `json:"identifier"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	BranchName  string `json:"branchName"`
	State       struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"state"`
	Team struct {
		ID   string `json:"id"`
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

// IssueUpdateInput represents the input for updating an issue
type IssueUpdateInput struct {
	Title       string
	Description string
	Priority    *int // pointer to distinguish 0 from unset
	StateID     string
	AssigneeID  string
}

// IssueUpdateResponse for the mutation response
type IssueUpdateResponse struct {
	Data struct {
		IssueUpdate struct {
			Success bool  `json:"success"`
			Issue   Issue `json:"issue"`
		} `json:"issueUpdate"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// SingleIssueResponse for fetching a single issue
type SingleIssueResponse struct {
	Data struct {
		Issue Issue `json:"issue"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// Comment represents a Linear comment
type Comment struct {
	ID        string `json:"id"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
	User      struct {
		Name string `json:"name"`
	} `json:"user"`
}

// CommentCreateResponse for the comment mutation response
type CommentCreateResponse struct {
	Data struct {
		CommentCreate struct {
			Success bool    `json:"success"`
			Comment Comment `json:"comment"`
		} `json:"commentCreate"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// IssueSearchResponse for search query
type IssueSearchResponse struct {
	Data struct {
		IssueSearch struct {
			Nodes []Issue `json:"nodes"`
		} `json:"issueSearch"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// CommentsResponse for fetching issue comments
type CommentsResponse struct {
	Data struct {
		Issue struct {
			Comments struct {
				Nodes []Comment `json:"nodes"`
			} `json:"comments"`
		} `json:"issue"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}
