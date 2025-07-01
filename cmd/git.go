package cmd

import (
	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:     "git",
	Short:   "Performs git operations",
	Long:    `Performs git operations when combined with sub commands`,
	Aliases: []string{"g"},
}

func init() {
	rootCmd.AddCommand(gitCmd)
}
