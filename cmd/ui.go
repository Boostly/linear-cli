package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))

	paginationStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	labelStyle = lipgloss.NewStyle().
			Width(12)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

type screenMode int

const (
	listMode screenMode = iota
	createMode
	selectTeamMode
	selectStateMode
	selectAssigneeMode
)

type createFormData struct {
	titleInput       textinput.Model
	descriptionInput textinput.Model
	teamID           string
	teamName         string
	stateID          string
	stateName        string
	assigneeID       string
	assigneeName     string
	priority         int
}

type model struct {
	// List mode fields
	issues   []Issue
	cursor   int
	selected map[int]struct{}
	loading  bool
	err      error

	// Screen mode
	mode screenMode

	// Create mode fields
	createForm   createFormData
	focusIndex   int
	teams        []Team
	states       []WorkflowState
	members      []User
	selectCursor int

	// Feedback
	successMsg string
	errorMsg   string
}

var priorityLabels = []string{"None", "Urgent", "High", "Medium", "Low"}

func initialModel() model {
	titleInput := textinput.New()
	titleInput.Placeholder = "Issue title"
	titleInput.CharLimit = 256
	titleInput.Width = 50

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Description (optional)"
	descriptionInput.CharLimit = 1000
	descriptionInput.Width = 50

	return model{
		selected: make(map[int]struct{}),
		loading:  true,
		mode:     listMode,
		createForm: createFormData{
			titleInput:       titleInput,
			descriptionInput: descriptionInput,
		},
	}
}

func (m model) Init() tea.Cmd {
	return loadIssues
}

func loadIssues() tea.Msg {
	client, err := getLinearClient()
	if err != nil {
		return errMsg{err}
	}

	issues, err := client.GetIssues()
	if err != nil {
		return errMsg{err}
	}

	return issuesLoadedMsg{issues}
}

func loadTeams() tea.Msg {
	client, err := getLinearClient()
	if err != nil {
		return errMsg{err}
	}

	teams, err := client.GetTeams()
	if err != nil {
		return errMsg{err}
	}

	return teamsLoadedMsg{teams}
}

func loadStates(teamID string) tea.Cmd {
	return func() tea.Msg {
		client, err := getLinearClient()
		if err != nil {
			return errMsg{err}
		}

		states, err := client.GetWorkflowStates(teamID)
		if err != nil {
			return errMsg{err}
		}

		return statesLoadedMsg{states}
	}
}

func loadMembers(teamID string) tea.Cmd {
	return func() tea.Msg {
		client, err := getLinearClient()
		if err != nil {
			return errMsg{err}
		}

		members, err := client.GetTeamMembers(teamID)
		if err != nil {
			return errMsg{err}
		}

		return membersLoadedMsg{members}
	}
}

func createIssue(input IssueCreateInput) tea.Cmd {
	return func() tea.Msg {
		client, err := getLinearClient()
		if err != nil {
			return createErrorMsg{err}
		}

		issue, err := client.CreateIssue(input)
		if err != nil {
			return createErrorMsg{err}
		}

		return issueCreatedMsg{issue}
	}
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type issuesLoadedMsg struct{ issues []Issue }
type teamsLoadedMsg struct{ teams []Team }
type statesLoadedMsg struct{ states []WorkflowState }
type membersLoadedMsg struct{ members []User }
type issueCreatedMsg struct{ issue *CreatedIssue }
type createErrorMsg struct{ err error }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Clear feedback messages on any key press
		m.successMsg = ""
		m.errorMsg = ""

		switch m.mode {
		case listMode:
			return m.updateListMode(msg)
		case createMode:
			return m.updateCreateMode(msg)
		case selectTeamMode:
			return m.updateSelectMode(msg, m.teams, func(i int) {
				m.createForm.teamID = m.teams[i].ID
				m.createForm.teamName = m.teams[i].Name
				// Reset state and assignee when team changes
				m.createForm.stateID = ""
				m.createForm.stateName = ""
				m.createForm.assigneeID = ""
				m.createForm.assigneeName = ""
				m.states = nil
				m.members = nil
			})
		case selectStateMode:
			return m.updateSelectMode(msg, m.states, func(i int) {
				m.createForm.stateID = m.states[i].ID
				m.createForm.stateName = m.states[i].Name
			})
		case selectAssigneeMode:
			return m.updateSelectMode(msg, m.members, func(i int) {
				m.createForm.assigneeID = m.members[i].ID
				m.createForm.assigneeName = m.members[i].Name
			})
		}

	case issuesLoadedMsg:
		m.issues = msg.issues
		m.loading = false

	case teamsLoadedMsg:
		m.teams = msg.teams
		m.loading = false

	case statesLoadedMsg:
		m.states = msg.states
		m.loading = false

	case membersLoadedMsg:
		m.members = msg.members
		m.loading = false

	case issueCreatedMsg:
		m.loading = false
		m.successMsg = fmt.Sprintf("Created: %s - %s", msg.issue.Identifier, msg.issue.Title)
		m.mode = listMode
		m.resetCreateForm()
		return m, loadIssues

	case createErrorMsg:
		m.loading = false
		m.errorMsg = fmt.Sprintf("Error: %v", msg.err)

	case errMsg:
		m.err = msg.err
		m.loading = false
	}

	return m, nil
}

func (m *model) updateListMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.issues)-1 {
			m.cursor++
		}

	case "enter", " ":
		if len(m.issues) > 0 {
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}

	case "c":
		m.mode = createMode
		m.focusIndex = 0
		m.createForm.titleInput.Focus()
		return m, loadTeams
	}

	return m, nil
}

func (m *model) updateCreateMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.mode = listMode
		m.resetCreateForm()
		return m, nil

	case "tab", "down":
		m.focusIndex++
		if m.focusIndex > 6 {
			m.focusIndex = 0
		}
		m.updateFormFocus()

	case "shift+tab", "up":
		m.focusIndex--
		if m.focusIndex < 0 {
			m.focusIndex = 6
		}
		m.updateFormFocus()

	case "left":
		if m.focusIndex == 3 { // Priority field
			m.createForm.priority--
			if m.createForm.priority < 0 {
				m.createForm.priority = 4
			}
		}

	case "right":
		if m.focusIndex == 3 { // Priority field
			m.createForm.priority++
			if m.createForm.priority > 4 {
				m.createForm.priority = 0
			}
		}

	case "enter":
		switch m.focusIndex {
		case 1: // Team field
			if len(m.teams) > 0 {
				m.mode = selectTeamMode
				m.selectCursor = 0
			}
		case 4: // State field
			if m.createForm.teamID != "" {
				if len(m.states) == 0 {
					m.loading = true
					return m, loadStates(m.createForm.teamID)
				}
				m.mode = selectStateMode
				m.selectCursor = 0
			}
		case 5: // Assignee field
			if m.createForm.teamID != "" {
				if len(m.members) == 0 {
					m.loading = true
					return m, loadMembers(m.createForm.teamID)
				}
				m.mode = selectAssigneeMode
				m.selectCursor = 0
			}
		case 6: // Submit button
			return m.submitCreateForm()
		}

	default:
		// Handle text input
		var cmd tea.Cmd
		if m.focusIndex == 0 {
			m.createForm.titleInput, cmd = m.createForm.titleInput.Update(msg)
			return m, cmd
		} else if m.focusIndex == 2 {
			m.createForm.descriptionInput, cmd = m.createForm.descriptionInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *model) updateSelectMode(msg tea.KeyMsg, items interface{}, onSelect func(int)) (tea.Model, tea.Cmd) {
	var length int
	switch v := items.(type) {
	case []Team:
		length = len(v)
	case []WorkflowState:
		length = len(v)
	case []User:
		length = len(v)
	}

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.mode = createMode
		return m, nil

	case "up", "k":
		if m.selectCursor > 0 {
			m.selectCursor--
		}

	case "down", "j":
		if m.selectCursor < length-1 {
			m.selectCursor++
		}

	case "enter":
		if length > 0 {
			onSelect(m.selectCursor)
			m.mode = createMode
		}
	}

	return m, nil
}

func (m *model) updateFormFocus() {
	m.createForm.titleInput.Blur()
	m.createForm.descriptionInput.Blur()

	switch m.focusIndex {
	case 0:
		m.createForm.titleInput.Focus()
	case 2:
		m.createForm.descriptionInput.Focus()
	}
}

func (m *model) resetCreateForm() {
	m.createForm.titleInput.SetValue("")
	m.createForm.descriptionInput.SetValue("")
	m.createForm.teamID = ""
	m.createForm.teamName = ""
	m.createForm.stateID = ""
	m.createForm.stateName = ""
	m.createForm.assigneeID = ""
	m.createForm.assigneeName = ""
	m.createForm.priority = 0
	m.focusIndex = 0
	m.teams = nil
	m.states = nil
	m.members = nil
}

func (m *model) submitCreateForm() (tea.Model, tea.Cmd) {
	// Validate required fields
	if m.createForm.titleInput.Value() == "" {
		m.errorMsg = "Title is required"
		m.focusIndex = 0
		m.createForm.titleInput.Focus()
		return m, nil
	}
	if m.createForm.teamID == "" {
		m.errorMsg = "Team is required"
		m.focusIndex = 1
		return m, nil
	}

	input := IssueCreateInput{
		Title:       m.createForm.titleInput.Value(),
		TeamID:      m.createForm.teamID,
		Description: m.createForm.descriptionInput.Value(),
		Priority:    m.createForm.priority,
		StateID:     m.createForm.stateID,
		AssigneeID:  m.createForm.assigneeID,
	}

	m.loading = true
	return m, createIssue(input)
}

func (m model) View() string {
	switch m.mode {
	case createMode:
		return m.viewCreateForm()
	case selectTeamMode:
		return m.viewTeamSelect()
	case selectStateMode:
		return m.viewStateSelect()
	case selectAssigneeMode:
		return m.viewAssigneeSelect()
	default:
		return m.viewList()
	}
}

func (m model) viewList() string {
	if m.loading {
		return "\n  Loading Linear issues...\n"
	}

	if m.err != nil {
		return fmt.Sprintf("\n  Error: %v\n", m.err)
	}

	s := titleStyle.Render("Linear Issues")
	s += "\n\n"

	if m.successMsg != "" {
		s += "  " + successStyle.Render(m.successMsg) + "\n\n"
	}

	if len(m.issues) == 0 {
		s += "  No issues found.\n"
	} else {
		for i, issue := range m.issues {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checked := " "
			if _, ok := m.selected[i]; ok {
				checked = "x"
			}

			line := fmt.Sprintf("%s [%s] %s - %s (%s)", cursor, checked, issue.Title, issue.State.Name, issue.Team.Name)

			if m.cursor == i {
				s += selectedItemStyle.Render(line)
			} else {
				s += itemStyle.Render(line)
			}
			s += "\n"
		}
	}

	s += "\n" + helpStyle.Render("q: quit | j/k: navigate | space: select | c: create issue")
	return s
}

func (m model) viewCreateForm() string {
	if m.loading {
		return "\n  Loading...\n"
	}

	s := titleStyle.Render("Create New Issue")
	s += "\n\n"

	if m.errorMsg != "" {
		s += "  " + errorStyle.Render(m.errorMsg) + "\n\n"
	}

	// Title field (index 0)
	label := labelStyle.Render("Title:")
	if m.focusIndex == 0 {
		s += fmt.Sprintf("  %s %s\n", focusedStyle.Render(">"), label+m.createForm.titleInput.View())
	} else {
		s += fmt.Sprintf("    %s %s\n", label, m.createForm.titleInput.View())
	}

	// Team field (index 1)
	label = labelStyle.Render("Team:")
	teamDisplay := blurredStyle.Render("[Select team...]")
	if m.createForm.teamName != "" {
		teamDisplay = m.createForm.teamName
	}
	if m.focusIndex == 1 {
		s += fmt.Sprintf("  %s %s %s\n", focusedStyle.Render(">"), label, teamDisplay)
	} else {
		s += fmt.Sprintf("    %s %s\n", label, teamDisplay)
	}

	// Description field (index 2)
	label = labelStyle.Render("Description:")
	if m.focusIndex == 2 {
		s += fmt.Sprintf("  %s %s\n", focusedStyle.Render(">"), label+m.createForm.descriptionInput.View())
	} else {
		s += fmt.Sprintf("    %s %s\n", label, m.createForm.descriptionInput.View())
	}

	// Priority field (index 3)
	label = labelStyle.Render("Priority:")
	priorityDisplay := priorityLabels[m.createForm.priority]
	if m.focusIndex == 3 {
		s += fmt.Sprintf("  %s %s < %s >\n", focusedStyle.Render(">"), label, priorityDisplay)
	} else {
		s += fmt.Sprintf("    %s %s\n", label, priorityDisplay)
	}

	// State field (index 4)
	label = labelStyle.Render("State:")
	stateDisplay := blurredStyle.Render("[Select state...]")
	if m.createForm.stateName != "" {
		stateDisplay = m.createForm.stateName
	} else if m.createForm.teamID == "" {
		stateDisplay = blurredStyle.Render("[Select team first]")
	}
	if m.focusIndex == 4 {
		s += fmt.Sprintf("  %s %s %s\n", focusedStyle.Render(">"), label, stateDisplay)
	} else {
		s += fmt.Sprintf("    %s %s\n", label, stateDisplay)
	}

	// Assignee field (index 5)
	label = labelStyle.Render("Assignee:")
	assigneeDisplay := blurredStyle.Render("[Select assignee...]")
	if m.createForm.assigneeName != "" {
		assigneeDisplay = m.createForm.assigneeName
	} else if m.createForm.teamID == "" {
		assigneeDisplay = blurredStyle.Render("[Select team first]")
	}
	if m.focusIndex == 5 {
		s += fmt.Sprintf("  %s %s %s\n", focusedStyle.Render(">"), label, assigneeDisplay)
	} else {
		s += fmt.Sprintf("    %s %s\n", label, assigneeDisplay)
	}

	// Submit button (index 6)
	s += "\n"
	if m.focusIndex == 6 {
		s += fmt.Sprintf("  %s [ Create Issue ]\n", focusedStyle.Render(">"))
	} else {
		s += "    [ Create Issue ]\n"
	}

	s += "\n" + helpStyle.Render("Tab: next | Shift+Tab: prev | Enter: select/submit | Esc: cancel")
	return s
}

func (m model) viewTeamSelect() string {
	s := titleStyle.Render("Select Team")
	s += "\n\n"

	for i, team := range m.teams {
		cursor := " "
		if m.selectCursor == i {
			cursor = ">"
		}

		line := fmt.Sprintf("%s %s (%s)", cursor, team.Name, team.Key)

		if m.selectCursor == i {
			s += selectedItemStyle.Render(line)
		} else {
			s += itemStyle.Render(line)
		}
		s += "\n"
	}

	s += "\n" + helpStyle.Render("j/k: navigate | Enter: select | Esc: cancel")
	return s
}

func (m model) viewStateSelect() string {
	s := titleStyle.Render("Select State")
	s += "\n\n"

	for i, state := range m.states {
		cursor := " "
		if m.selectCursor == i {
			cursor = ">"
		}

		line := fmt.Sprintf("%s %s", cursor, state.Name)

		if m.selectCursor == i {
			s += selectedItemStyle.Render(line)
		} else {
			s += itemStyle.Render(line)
		}
		s += "\n"
	}

	s += "\n" + helpStyle.Render("j/k: navigate | Enter: select | Esc: cancel")
	return s
}

func (m model) viewAssigneeSelect() string {
	s := titleStyle.Render("Select Assignee")
	s += "\n\n"

	for i, member := range m.members {
		cursor := " "
		if m.selectCursor == i {
			cursor = ">"
		}

		line := fmt.Sprintf("%s %s (%s)", cursor, member.Name, member.Email)

		if m.selectCursor == i {
			s += selectedItemStyle.Render(line)
		} else {
			s += itemStyle.Render(line)
		}
		s += "\n"
	}

	s += "\n" + helpStyle.Render("j/k: navigate | Enter: select | Esc: cancel")
	return s
}

func runTUI() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}
