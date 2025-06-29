package cmd

import (
	"fmt"

	"github.com/the-mclain-train/gowo/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectAddCmd = &cobra.Command{
	Use:   "add <projectName> <repo1> [repo2...]",
	Short: "Add a new project",
	Long:  `Adds a new project definition with a name and a list of repositories.`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		repos := args[1:]

		projectsKey := "projects"
		if viper.IsSet(projectsKey + "." + projectName) {
			return fmt.Errorf("project '%s' already exists", projectName)
		}

		newProject := config.Project{Repositories: repos}
		viper.Set(projectsKey+"."+projectName, newProject)

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ Project '%s' added with %d repositories.\n", projectName, len(repos))
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectAddCmd)
}
