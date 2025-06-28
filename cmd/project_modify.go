package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/the-mclain-train/gorkspace/internal/config"
)

var (
	addFlag     bool
	removeFlag  bool
	modProjFlag string
	repoFlag    string
)

// modifyCmd represents the modify command
var modifyCmd = &cobra.Command{
	Use:   "modify -p <project> [-a | -r] -R <repo>",
	Short: "Modifies a project",
	Long:  `Modifies a project by adding or removing repositories`,
	RunE: func(cmd *cobra.Command, args []string) error {

		fullProjectKey := "projects." + modProjFlag

		var project config.Project
		if err := config.GetProject(&project, modProjFlag); err != nil {
			return fmt.Errorf("error getting project %s", modProjFlag)
		}

		if addFlag {
			_, ok := config.FindRepo(&project, repoFlag)
			if ok {
				return fmt.Errorf("%s is already in %s", repoFlag, modProjFlag)
			}

			project.Repositories = append(project.Repositories, repoFlag)

			viper.Set(fullProjectKey, project)
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("error writing configuration file")
			}

			fmt.Printf("✅ %s successfully added to %s!\n", repoFlag, modProjFlag)
		}

		if removeFlag {
			rindex, ok := config.FindRepo(&project, repoFlag)

			if !ok {
				return fmt.Errorf("%s not found in %s", repoFlag, modProjFlag)
			}

			// Remove repo from the project and save updated config
			project.Repositories = append(project.Repositories[:rindex], project.Repositories[rindex+1:]...)

			viper.Set(fullProjectKey, project)

			if err := viper.WriteConfig(); err != nil {
				log.Fatalf("Error writing configuration file: %s", err)
			}

			fmt.Printf("✅ %s successfully removed from %s!\n", repoFlag, modProjFlag)
		}

		return nil
	},
}

func init() {
	projectCmd.AddCommand(modifyCmd)

	modifyCmd.Flags().BoolVarP(&addFlag, "add", "a", false, "adds repository to a project")
	modifyCmd.Flags().StringVarP(&modProjFlag, "project", "p", "", "project to modify")
	modifyCmd.Flags().BoolVarP(&removeFlag, "remove", "r", false, "removes repository from a project")
	modifyCmd.Flags().StringVarP(&repoFlag, "repo", "R", "", "repository to update project with")

	modifyCmd.MarkFlagsOneRequired("add", "remove")
	modifyCmd.MarkFlagsMutuallyExclusive("add", "remove")
	modifyCmd.MarkFlagRequired("project")
	modifyCmd.MarkFlagRequired("repo")
}
