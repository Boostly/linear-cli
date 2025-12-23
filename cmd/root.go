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
	Use:   "linear-cli",
	Short: "A CLI for Linear.app",
	Long: `Linear CLI is a command-line interface for interacting with Linear.app.
It allows you to fetch issues, create and update them, search, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
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
	rootCmd.AddCommand(updateCmd)
}
