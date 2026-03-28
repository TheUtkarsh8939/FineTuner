package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func cleanPath(configPath string) (string, error) {
	trimmedPath := strings.TrimSpace(configPath)
	if trimmedPath == "" {
		return "", fmt.Errorf("config path is empty")
	}

	// Convert relative paths like ./config.json into an absolute path
	// based on the current working directory.
	absPath, err := filepath.Abs(trimmedPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve config path %q: %w", trimmedPath, err)
	}
	absPath = filepath.Clean(absPath)
	return absPath, nil
}
