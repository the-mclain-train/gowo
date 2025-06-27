package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootDirflag string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configurations",
	Long:  `The config command allows you to manage local configurations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configKey := "config"

		viper.Set(configKey+"."+"workspaces_directory", rootDirflag)

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✅ Workspaces root set to %s\n", rootDirflag)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&rootDirflag, "root", "r", "", "Root directory for creating projects")
	configCmd.MarkFlagRequired("root")
}
