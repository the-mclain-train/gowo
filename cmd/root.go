package cmd

import (
	"fmt"
	"os"

	"github.com/the-mclain-train/gorkspace/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gorkspace",
	Short: "gorkspace is a CLI tool to manage Go workspaces based on project definitions.",
	Long: `A fast and flexible CLI for Go developers to manage complex Go workspaces.
Define projects as collections of repositories and let gorkspace handle the setup.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your command: '%s'", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitViper)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}
