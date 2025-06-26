package cmd

import (
	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage Go workspaces",
	Long: `A workspace is a directory on your filesystem that contains the
cloned repositories of a Project, with a go.work file linking them.`,
	Aliases: []string{"ws"},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
}
