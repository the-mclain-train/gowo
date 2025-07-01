package git

import (
	"fmt"
	"os"
	"path/filepath"
)

// findProjectRoot searches current and all parent directories for the .gowo-meta file
func findProjectRoot() (string, error) {
	const target = ".gowo-meta"

	startDir, err := os.Getwd()
	if err != nil {
		return "oops", fmt.Errorf("unable to determine current directory, got error: %v", err)
	}

	// get absolute path, relative bad
	currentDir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for starting directory: %w", err)
	}

	// loop indefinitely until we find it or we give up
	for {
		filePath := filepath.Join(currentDir, target)

		if _, err := os.Stat(filePath); err == nil {
			// file found, return its path
			return currentDir, nil
		} else if !os.IsNotExist(err) {
			// an error other than "not found", maybe perms, etc
			return "", fmt.Errorf("error checking for file at %s: %w", filePath, err)
		}

		// move up to the parent directory
		parentDir := filepath.Dir(currentDir)

		// if the parent directory is the same as the current one, we're at the root.
		if parentDir == currentDir {
			return "", fmt.Errorf("%s not found. can't find project root", target)
		}

		// next round, maybe here, maybe not
		currentDir = parentDir
	}
}

// Get projects finds all repos in a workspace, when ran from within the workspace
func GetProjects() ([]string, error) {
	var ps []string

	root, err := findProjectRoot()
	if err != nil {
		return ps, fmt.Errorf("unable to determine repos in this workspace: %v", err)
	}

	all, err := os.ReadDir(root)
	if err != nil {
		return ps, err
	}

	// add the dirs to ps, these should all be repos
	for _, item := range all {
		if item.IsDir() {
			ps = append(ps, item.Name())
		}
	}

	return ps, nil
}
