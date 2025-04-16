package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Config struct {
	Workspaces map[string]*WorkspaceConfig `toml:"workspace"`
	Apps       map[string]*AppConfig       `toml:"app"`
}

type WorkspaceConfig struct {
	Name           string `toml:"-"`
	Representation string `toml:"representation"`
}

type AppConfig struct {
	Name        string   `toml:"-"`
	ID          string   `toml:"id"`
	Command     string   `toml:"cmd,omitempty"`
	Size        int      `toml:"size,omitempty"`
	PostCommand []string `toml:"post_cmd,omitempty"`
}

func (c *Config) Validate() error {
	for _, workspace := range c.Workspaces {
		if err := workspace.validate(c); err != nil {
			return fmt.Errorf("workspace '%s' invalid: %w", workspace.Name, err)
		}
	}

	for _, app := range c.Apps {
		if err := app.validate(); err != nil {
			return fmt.Errorf("app '%s' invalid: %w", app.Name, err)
		}
	}

	return nil
}

func (w *WorkspaceConfig) validate(c *Config) error {
	if w.Representation == "" {
		return fmt.Errorf("representation is required")
	}

	if err := w.validateLayout(); err != nil {
		return fmt.Errorf("representation is invalid: %w", err)
	}

	if err := w.validateAppsName(c); err != nil {
		return fmt.Errorf("representation is invalid: %w", err)
	}

	return nil
}

func (w *WorkspaceConfig) validateLayout() error {
	validPrefix := [4]string{"H[", "T[", "V[", "S["}
	for _, prefix := range validPrefix {
		if !strings.HasPrefix(w.Representation, prefix) {
			return fmt.Errorf("must start with 'H[', 'T[', 'V[' or 'S['")
		}
	}

	bracketCount := 0
	for _, char := range w.Representation {
		switch char {
		case '[':
			bracketCount++
		case ']':
			bracketCount--
		default:
		}
	}
	if bracketCount != 0 {
		return fmt.Errorf("brackets are unbalanced")
	}

	return nil
}

func (w *WorkspaceConfig) validateAppsName(c *Config) error {
	representation := w.Representation

	layoutMarkers := [6]string{"H[", "T[", "V[", "S[", "[", "]"}
	for _, marker := range layoutMarkers {
		representation = strings.ReplaceAll(representation, marker, " ")
	}

	apps := strings.Fields(representation)
	for _, app := range apps {
		reg := `^[a-zA-Z0-9_-]+$`
		if !regexp.MustCompile(reg).MatchString(app) {
			return fmt.Errorf("name contains invalid characters - only letters, numbers, underscore, and hyphen are allowed")
		}

		if _, exists := c.Apps[app]; !exists {
			return fmt.Errorf("not defined in the configuration")
		}
	}

	return nil
}

func (a *AppConfig) validate() error {
	if a.ID == "" {
		return fmt.Errorf("id is required")
	}

	if a.Command == "" && !a.isDesktopEntry() {
		return fmt.Errorf("no match with desktop entries. Provide a command to launch it")
	}

	return nil
}

func (a *AppConfig) isDesktopEntry() bool {
	homeDir, _ := os.UserHomeDir()
	desktopEntryDirs := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		filepath.Join(homeDir, ".local/share/applications"),
	}

	for _, dir := range desktopEntryDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		desktopFile := filepath.Join(dir, a.ID+".desktop")
		if _, err := os.Stat(desktopFile); err == nil {
			return true
		}
	}

	return false
}
