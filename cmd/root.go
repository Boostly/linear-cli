package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	statusFilter    string
	showDescription bool
)

var rootCmd = &cobra.Command{
	Use:   "linear",
	Short: "A CLI and TUI for Linear.app integration",
	Long: `Linear TUI is a command-line interface and terminal user interface 
for interacting with Linear.app. It allows you to fetch issues, filter by status,
and view detailed information about your Linear workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runTUI(); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(meCmd)
	rootCmd.AddCommand(createCmd)

	rootCmd.SetHelpTemplate(`{{.Long}}

Usage:
  {{.UseLine}}

Available Commands:
  list        List all issues in the workspace
  me          List issues assigned to you
  create      Create a new issue
  tui         Run the interactive TUI (default)
  help        Help about any command

Flags:
  -h, --help   Help for {{.Name}}

Use "{{.CommandPath}} [command] --help" for more information about a command.

Examples:
  linear                              # Run the interactive TUI
  linear tui                          # Run the interactive TUI explicitly
  linear list                         # List all issues
  linear list -s "Backlog"            # List all backlog issues
  linear me                           # List your assigned issues
  linear me -s "In Progress"          # List your in-progress issues
  linear create -t "Bug fix" -T ENG   # Create a new issue
`)
}
