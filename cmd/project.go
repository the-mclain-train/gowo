package cmd

import (
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project definitions",
	Long: `A project is a named collection of git repositories.
These commands allow you to add, remove, list, and inspect projects.`,
	Aliases: []string{"p"},
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
