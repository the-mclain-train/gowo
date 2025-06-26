package cmd

import (
	"fmt"

	"github.com/the-mclain-train/gorkspace/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectShowCmd = &cobra.Command{
	Use:   "show <projectName>",
	Short: "Show details for a specific project",
	Long:  `Displays the list of repositories associated with a given project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		key := "projects." + projectName

		if !viper.IsSet(key) {
			return fmt.Errorf("project '%s' not found", projectName)
		}

		var project config.Project
		if err := viper.UnmarshalKey(key, &project); err != nil {
			return fmt.Errorf("failed to decode project data: %w", err)
		}

		fmt.Printf("Project: %s\n", projectName)
		fmt.Println("Repositories:")
		for _, repo := range project.Repositories {
			fmt.Printf("  - %s\n", repo)
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(projectShowCmd)
}
