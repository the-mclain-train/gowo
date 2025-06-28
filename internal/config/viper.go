package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
	configDir  = ".gorkspace"
)

// Project defines the structure for a project with its repositories.
type Project struct {
	Repositories []string `mapstructure:"repositories"`
}

// EnsureConfig ensures the configuration file and directory exist.
func EnsureConfig() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}
	configPath := filepath.Join(home, configDir)
	configFilePath := filepath.Join(configPath, fmt.Sprintf("%s.%s", configName, configType))

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("Creating config directory at %s\n", configPath)
		if err := os.MkdirAll(configPath, 0755); err != nil {
			return "", fmt.Errorf("could not create config directory: %w", err)
		}
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Printf("Creating default config file at %s\n", configFilePath)
		// Create a bare-bones config file
		initialConfig := []byte("projects: {}\n")
		if err := os.WriteFile(configFilePath, initialConfig, 0644); err != nil {
			return "", fmt.Errorf("could not create config file: %w", err)
		}
	}
	return configFilePath, nil
}

// InitViper initializes viper to read from the config file.
func InitViper() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(home, configDir)
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; we can guide the user to create one
			// or create one automatically. We'll do the latter.
			if _, err := EnsureConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating default config: %v\n", err)
				os.Exit(1)
			}
			// Reread the config now that it's created
			if err := viper.ReadInConfig(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading created config: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Config file was found but another error was produced
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}
}

// GetProject validates a project exists and sets it to p
func GetProject(p *Project, key string) error {
	fullKey := "projects." + key
	if !viper.IsSet(fullKey) {
		return fmt.Errorf("project '%s' not found. Use 'gws project ls' to see available projects", key)
	}

	if err := viper.UnmarshalKey(fullKey, p); err != nil {
		return fmt.Errorf("could not decode project '%s': %w", key, err)
	}
	return nil

}

// FindRepo validates a repo exists in a project and returns the index
func FindRepo(p *Project, repo string) (int, bool) {
	rindex := -1
	for i, r := range p.Repositories {
		if r == repo {
			return i, true
		}
	}
	return rindex, false
}
