package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var projectLsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List all saved projects",
	Long:    `Lists the names of all projects saved in the configuration file.`,
	Aliases: []string{"list"},
	Run: func(cmd *cobra.Command, args []string) {
		projects := viper.GetStringMap("projects")
		if len(projects) == 0 {
			fmt.Println("No projects found. Use 'gowo project add' to create one.")
			return
		}

		fmt.Println("Available Projects:")
		keys := make([]string, 0, len(projects))
		for k := range projects {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("- %s\n", k)
		}
	},
}

func init() {
	projectCmd.AddCommand(projectLsCmd)
}
