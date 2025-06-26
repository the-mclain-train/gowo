package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectRmCmd = &cobra.Command{
	Use:     "rm <projectName>",
	Short:   "Remove a project",
	Long:    `Removes a project definition from the configuration.`,
	Aliases: []string{"remove"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		projectsKey := "projects"
		fullKey := projectsKey + "." + projectName

		if !viper.IsSet(fullKey) {
			return fmt.Errorf("project '%s' not found", projectName)
		}

		// To "remove" a key, we get the whole map, delete the key, and set it back.
		allProjects := viper.GetStringMap(projectsKey)
		delete(allProjects, projectName)
		viper.Set(projectsKey, allProjects)

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save config after removing project: %w", err)
		}

		fmt.Printf("✅ Project '%s' removed.\n", projectName)
		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectRmCmd)
}
