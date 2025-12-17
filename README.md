# Linear CLI

A command-line interface for Linear.app. Manage issues, search, comment, and streamline your development workflow directly from the terminal.

## Features

- **Issue Management**: List, create, update, and search issues
- **Quick Actions**: Open issues in browser, add comments
- **Start Command**: One command to assign yourself, set "In Progress", and create a git branch
- **Status Filtering**: Filter issues by any Linear state
- **User Queries**: Get issues assigned to you

## Installation

### Prerequisites
- Go 1.21 or later
- Linear.app account with API access

### Install

```bash
go install github.com/Boostly/linear-cli@latest
```

### Setup

1. Get your API key from [Linear Settings > API](https://linear.app/settings/api)

2. Create a `.env` file:
   ```bash
   # Option 1: Project directory
   echo "LINEAR_API_KEY=your_api_key" > .env

   # Option 2: Global config
   mkdir -p ~/.config/linear
   echo "LINEAR_API_KEY=your_api_key" > ~/.config/linear/.env
   ```

## Commands

### List Issues

```bash
# Your assigned issues
linear-cli me
linear-cli me -s "In Progress"
linear-cli me -d                      # Include descriptions

# All workspace issues
linear-cli list
linear-cli list -s "Backlog" -d
```

### Search

```bash
linear-cli search "authentication bug"
linear-cli search "login" -d          # With descriptions
```

### Open in Browser

```bash
linear-cli open ENG-123
```

### Add Comment

```bash
linear-cli comment ENG-123 -m "Fixed in latest commit"
linear-cli comment ENG-123            # Interactive prompt
```

### Create Issue

```bash
linear-cli create -t "Fix login bug" -T ENG
linear-cli create -t "New feature" -T ENG -d "Description here" -p 2
```

Flags:
- `-t, --title` (required) - Issue title
- `-T, --team` (required) - Team key (e.g., ENG)
- `-d, --description` - Issue description
- `-p, --priority` - Priority (1-4)
- `-s, --state` - Initial state
- `-a, --assignee` - Assignee

### Update Issue

```bash
linear-cli update -i ENG-123 -s "Done"
linear-cli update -i ENG-123 -a "john@example.com" -p 1
```

Flags:
- `-i, --issue` (required) - Issue identifier
- `-t, --title` - New title
- `-d, --description` - New description
- `-s, --state` - New state
- `-a, --assignee` - New assignee
- `-p, --priority` - New priority

### Start Working (Power Command)

The `start` command streamlines beginning work on an issue:

```bash
linear-cli start ENG-123
```

This will:
1. Assign the issue to you (if not already)
2. Set status to "In Progress"
3. Create and checkout a git branch (using Linear's branch naming)

Options:
```bash
linear-cli start ENG-123 --no-branch     # Skip git branch creation
linear-cli start ENG-123 -b "my-branch"  # Custom branch name
```

## Examples

### Daily Workflow

```bash
# What am I working on?
linear-cli me -s "In Progress"

# Pick up new work
linear-cli me -s "Backlog"
linear-cli start ENG-456

# Quick search
linear-cli search "payment"

# Done with a task
linear-cli update -i ENG-456 -s "Done"
linear-cli comment ENG-456 -m "Completed - see PR #123"
```

### With Claude Code

```bash
# Get context for current work
linear-cli me -s "In Progress" -d

# Find issues to work on
linear-cli search "needs review" -d
```

## Building from Source

```bash
git clone https://github.com/Boostly/linear-cli.git
cd linear-cli
go build -o linear-cli .
```

## License

MIT
