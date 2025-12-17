package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [issue-identifier]",
	Short: "Open an issue in the browser",
	Long:  `Open a Linear issue in your default web browser by its identifier (e.g., ENG-123).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]

		client, err := getLinearClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		issue, err := client.GetIssue(identifier)
		if err != nil {
			fmt.Printf("Error fetching issue: %v\n", err)
			os.Exit(1)
		}

		if issue.URL == "" {
			fmt.Printf("Error: no URL found for issue %s\n", identifier)
			os.Exit(1)
		}

		fmt.Printf("Opening %s in browser...\n", issue.Identifier)
		if err := openBrowser(issue.URL); err != nil {
			fmt.Printf("Error opening browser: %v\n", err)
			os.Exit(1)
		}
	},
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Start()
}

func init() {
	rootCmd.AddCommand(openCmd)
}
