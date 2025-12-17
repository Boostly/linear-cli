package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const linearAPIURL = "https://api.linear.app/graphql"

type LinearClient struct {
	apiKey string
	client *http.Client
}

func NewLinearClient(apiKey string) *LinearClient {
	return &LinearClient{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (lc *LinearClient) GetCurrentUser() (*User, error) {
	query := `
		query {
			viewer {
				id
				name
				email
			}
		}
	`

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var viewerResp ViewerResponse
	if err := json.NewDecoder(resp.Body).Decode(&viewerResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &viewerResp.Data.Viewer, nil
}

func (lc *LinearClient) GetIssues() ([]Issue, error) {
	query := `
		query {
			issues(first: 20) {
				nodes {
					id
					identifier
					title
					description
					state {
						id
						name
					}
					team {
						id
						name
					}
					assignee {
						id
						name
					}
					priority
				}
			}
		}
	`

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var issuesResp IssuesResponse
	if err := json.NewDecoder(resp.Body).Decode(&issuesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return issuesResp.Data.Issues.Nodes, nil
}

func (lc *LinearClient) GetMyIssues() ([]Issue, error) {
	// First get current user info
	user, err := lc.GetCurrentUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	// Query for issues assigned to the current user
	query := fmt.Sprintf(`
		query {
			issues(first: 50, filter: { assignee: { id: { eq: "%s" } } }) {
				nodes {
					id
					identifier
					title
					description
					state {
						id
						name
					}
					team {
						id
						name
					}
					assignee {
						id
						name
					}
					priority
				}
			}
		}
	`, user.ID)

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var issuesResp IssuesResponse
	if err := json.NewDecoder(resp.Body).Decode(&issuesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return issuesResp.Data.Issues.Nodes, nil
}

func (lc *LinearClient) GetIssuesByStatus(status string) ([]Issue, error) {
	query := fmt.Sprintf(`
		query {
			issues(first: 50, filter: { state: { name: { eq: "%s" } } }) {
				nodes {
					id
					identifier
					title
					description
					state {
						id
						name
					}
					team {
						id
						name
					}
					assignee {
						id
						name
					}
					priority
				}
			}
		}
	`, status)

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var issuesResp IssuesResponse
	if err := json.NewDecoder(resp.Body).Decode(&issuesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return issuesResp.Data.Issues.Nodes, nil
}

func (lc *LinearClient) GetMyIssuesByStatus(status string) ([]Issue, error) {
	// First get current user info
	user, err := lc.GetCurrentUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	// Query for issues assigned to the current user with specific status
	query := fmt.Sprintf(`
		query {
			issues(first: 50, filter: {
				assignee: { id: { eq: "%s" } },
				state: { name: { eq: "%s" } }
			}) {
				nodes {
					id
					identifier
					title
					description
					state {
						id
						name
					}
					team {
						id
						name
					}
					assignee {
						id
						name
					}
					priority
				}
			}
		}
	`, user.ID, status)

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var issuesResp IssuesResponse
	if err := json.NewDecoder(resp.Body).Decode(&issuesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return issuesResp.Data.Issues.Nodes, nil
}

func getLinearClient() (*LinearClient, error) {
	apiKey := os.Getenv("LINEAR_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("LINEAR_API_KEY environment variable not set")
	}
	return NewLinearClient(apiKey), nil
}

func (lc *LinearClient) GetTeams() ([]Team, error) {
	query := `
		query {
			teams {
				nodes {
					id
					name
					key
				}
			}
		}
	`

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var teamsResp TeamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return teamsResp.Data.Teams.Nodes, nil
}

func (lc *LinearClient) GetWorkflowStates(teamID string) ([]WorkflowState, error) {
	query := fmt.Sprintf(`
		query {
			workflowStates(filter: { team: { id: { eq: "%s" } } }) {
				nodes {
					id
					name
					type
				}
			}
		}
	`, teamID)

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var statesResp WorkflowStatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&statesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return statesResp.Data.WorkflowStates.Nodes, nil
}

func (lc *LinearClient) GetTeamMembers(teamID string) ([]User, error) {
	query := fmt.Sprintf(`
		query {
			team(id: "%s") {
				members {
					nodes {
						id
						name
						email
					}
				}
			}
		}
	`, teamID)

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var membersResp TeamMembersResponse
	if err := json.NewDecoder(resp.Body).Decode(&membersResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return membersResp.Data.Team.Members.Nodes, nil
}

func (lc *LinearClient) CreateIssue(input IssueCreateInput) (*CreatedIssue, error) {
	mutation := `
		mutation IssueCreate($input: IssueCreateInput!) {
			issueCreate(input: $input) {
				success
				issue {
					id
					title
					identifier
					url
					state {
						name
					}
					team {
						name
					}
				}
			}
		}
	`

	// Build input map
	inputMap := map[string]interface{}{
		"title":  input.Title,
		"teamId": input.TeamID,
	}

	if input.Description != "" {
		inputMap["description"] = input.Description
	}
	if input.Priority > 0 {
		inputMap["priority"] = input.Priority
	}
	if input.StateID != "" {
		inputMap["stateId"] = input.StateID
	}
	if input.AssigneeID != "" {
		inputMap["assigneeId"] = input.AssigneeID
	}

	variables := map[string]interface{}{
		"input": inputMap,
	}

	reqBody := GraphQLRequestWithVariables{
		Query:     mutation,
		Variables: variables,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var createResp IssueCreateFullResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(createResp.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", createResp.Errors[0].Message)
	}

	if !createResp.Data.IssueCreate.Success {
		return nil, fmt.Errorf("failed to create issue")
	}

	return &createResp.Data.IssueCreate.Issue, nil
}

func (lc *LinearClient) GetIssue(identifier string) (*Issue, error) {
	query := fmt.Sprintf(`
		query {
			issue(id: "%s") {
				id
				identifier
				title
				description
				state {
					id
					name
				}
				team {
					id
					name
				}
				assignee {
					id
					name
				}
				priority
			}
		}
	`, identifier)

	reqBody := GraphQLRequest{Query: query}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var issueResp SingleIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issueResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(issueResp.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", issueResp.Errors[0].Message)
	}

	if issueResp.Data.Issue.ID == "" {
		return nil, fmt.Errorf("issue not found: %s", identifier)
	}

	return &issueResp.Data.Issue, nil
}

func (lc *LinearClient) UpdateIssue(issueID string, input IssueUpdateInput) (*Issue, error) {
	mutation := `
		mutation IssueUpdate($id: String!, $input: IssueUpdateInput!) {
			issueUpdate(id: $id, input: $input) {
				success
				issue {
					id
					identifier
					title
					description
					state {
						id
						name
					}
					team {
						id
						name
					}
					assignee {
						id
						name
					}
					priority
				}
			}
		}
	`

	// Build input map - only include fields that are being changed
	inputMap := map[string]interface{}{}

	if input.Title != "" {
		inputMap["title"] = input.Title
	}
	if input.Description != "" {
		inputMap["description"] = input.Description
	}
	if input.Priority != nil {
		inputMap["priority"] = *input.Priority
	}
	if input.StateID != "" {
		inputMap["stateId"] = input.StateID
	}
	if input.AssigneeID != "" {
		inputMap["assigneeId"] = input.AssigneeID
	}

	variables := map[string]interface{}{
		"id":    issueID,
		"input": inputMap,
	}

	reqBody := GraphQLRequestWithVariables{
		Query:     mutation,
		Variables: variables,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", linearAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", lc.apiKey)

	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var updateResp IssueUpdateResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(updateResp.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", updateResp.Errors[0].Message)
	}

	if !updateResp.Data.IssueUpdate.Success {
		return nil, fmt.Errorf("failed to update issue")
	}

	return &updateResp.Data.IssueUpdate.Issue, nil
}
