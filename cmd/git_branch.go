package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/the-mclain-train/gowo/internal/git"
	"github.com/the-mclain-train/gowo/internal/gowo"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch <branch>",
	Args:  cobra.ExactArgs(1),
	Short: "Switch git branch",
	Long: `Switches current project (from within project/workspace) to the specified branch. 
	
If <branch> doesn't exist, it will be created.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO strip whitespace on branch if it becomes an issue, can use strings.Fields
		branch := args[0]

		repos, err := git.GetProjects()
		if err != nil {
			return err
		}

		for _, repo := range repos {
			gowo.RunCommand(repo, "git", "switch", "-c", branch)
			if err != nil {
				// maybe branch exists, try without -c flag
				gowo.RunCommand(repo, "git", "switch", branch)
				if err != nil {
					// k something is wrong cry uncle
					return fmt.Errorf("unable to switch to switch to %s in %s: %v", branch, repo, err)
				}
			}
		}

		fmt.Printf("\n✅ Success! All project on branch %s!\n", branch)
		return nil
	},
}

func init() {
	gitCmd.AddCommand(branchCmd)
}
