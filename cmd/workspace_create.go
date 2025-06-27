package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/the-mclain-train/gorkspace/internal/config"
	gws "github.com/the-mclain-train/gorkspace/internal/gorkspace"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	projectFlag   string
	pathFlag      string
	workspaceName string
)

var workspaceCreateCmd = &cobra.Command{
	Use:   "create [flags]",
	Short: "Create a new Go workspace from a project",
	Long: `Creates a directory, clones all project repositories into it,
initializes a Go workspace, and adds all discovered modules.`,
	// Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 1. Pre-flight checks
		if err := gws.CheckDependencies(); err != nil {
			return err
		}

		// 2. Validate Project
		projectKey := "projects." + projectFlag
		if !viper.IsSet(projectKey) {
			return fmt.Errorf("project '%s' not found. Use 'gws project ls' to see available projects", projectFlag)
		}

		var project config.Project
		if err := viper.UnmarshalKey(projectKey, &project); err != nil {
			return fmt.Errorf("could not decode project '%s': %w", projectFlag, err)
		}

		// 3. Determine Workspace Path
		parentDir := pathFlag
		if parentDir == "" {
			parentDir = viper.GetString("config.workspaces_directory")
			fmt.Println("")
		}

		if parentDir == "" {
			// If still empty, default to current directory
			parentDir = "."
		}

		// Expand tilde
		if strings.HasPrefix(parentDir, "~") {
			home, _ := os.UserHomeDir()
			parentDir = filepath.Join(home, parentDir[1:])
		}

		workspacePath, err := filepath.Abs(filepath.Join(parentDir, workspaceName))
		if err != nil {
			return fmt.Errorf("could not determine absolute workspace path: %w", err)
		}

		if _, err := os.Stat(workspacePath); !os.IsNotExist(err) {
			return fmt.Errorf("workspace directory '%s' already exists", workspacePath)
		}

		// 4. Create Workspace Directory
		fmt.Printf("🚀 Creating workspace '%s' at %s\n", workspaceName, workspacePath)
		if err := os.MkdirAll(workspacePath, 0755); err != nil {
			return fmt.Errorf("could not create workspace directory: %w", err)
		}

		// 5. Clone Repos
		fmt.Println("\nCloning repositories...")
		for _, repoURL := range project.Repositories {
			// Construct a proper git clone URL
			fullRepoURL := "https://" + repoURL
			err := gws.RunCommand(workspacePath, "git", "clone", fullRepoURL)
			if err != nil {
				// Clean up created directory on failure
				os.RemoveAll(workspacePath)
				return fmt.Errorf("failed to clone repository %s: %w", repoURL, err)
			}
		}

		// 6. Initialize Go Workspace
		fmt.Println("\nInitializing Go workspace...")
		if err := gws.RunCommand(workspacePath, "go", "work", "init"); err != nil {
			os.RemoveAll(workspacePath)
			return fmt.Errorf("failed to initialize go workspace: %w", err)
		}

		// 7. Find and Use Modules
		fmt.Println("\nScanning for Go modules...")
		modules, err := gws.FindGoModules(workspacePath)
		if err != nil {
			return fmt.Errorf("failed to find go modules: %w", err)
		}

		if len(modules) == 0 {
			fmt.Println("⚠️ No Go modules found in the cloned repositories.")
		} else {
			fmt.Println("Found modules, adding to workspace...")
			// Create the 'go work use' command with all modules at once
			args := append([]string{"work", "use"}, modules...)
			if err := gws.RunCommand(workspacePath, "go", args...); err != nil {
				os.RemoveAll(workspacePath)
				return fmt.Errorf("failed to add modules to workspace: %w", err)
			}
		}

		// 8. Create Metadata File
		metaFilePath := filepath.Join(workspacePath, ".gorkspace-meta")
		metaContent := []byte(fmt.Sprintf("project: %s\n", projectFlag))
		if err := os.WriteFile(metaFilePath, metaContent, 0644); err != nil {
			fmt.Printf("Warning: could not write metadata file: %v\n", err)
		}

		fmt.Printf("\n✅ Success! Workspace '%s' is ready.\n", workspaceName)
		fmt.Printf("To get started, run: cd %s\n", workspacePath)
		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceCreateCmd)

	// Required flag for the project name
	workspaceCreateCmd.Flags().StringVarP(&projectFlag, "project", "p", "", "The name of the project to use for the workspace (required)")
	workspaceCreateCmd.Flags().StringVarP(&workspaceName, "name", "n", "", "The name of the workspace to create")
	workspaceCreateCmd.MarkFlagRequired("project")
	workspaceCreateCmd.MarkFlagRequired("name")

	// Optional flag for the parent directory path
	workspaceCreateCmd.Flags().StringVarP(&pathFlag, "directory", "d", "", "The parent directory to create the workspace in (defaults to config or current dir)")
}
