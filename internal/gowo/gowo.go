package gowo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// CheckDependencies ensures that git and go are available in the PATH.
func CheckDependencies() error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("'git' command not found: please install git and ensure it's in your PATH")
	}
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("'go' command not found: please install Go and ensure it's in your PATH")
	}
	return nil
}

// RunCommand executes a shell command and streams its output.
func RunCommand(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("==> Running in %s: %s %s\n", dir, name, strings.Join(args, " "))
	return cmd.Run()
}

// FindGoModules recursively finds all go.mod files in a given directory.
func FindGoModules(root string) ([]string, error) {
	var modules []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "go.mod" {
			// We need the directory containing go.mod, not the file itself
			modulePath, err := filepath.Rel(root, filepath.Dir(path))
			if err != nil {
				return err
			}
			modules = append(modules, "./"+modulePath)
		}
		return nil
	})
	return modules, err
}

// GetModuleName reads a go.mod file and returns the module name.
func GetModuleName(modFilePath string) (string, error) {
	data, err := os.ReadFile(modFilePath)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(data), nil
}
