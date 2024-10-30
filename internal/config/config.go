// Package config for the redirect application.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// File for configuring the proxy application behaviours.
type File struct {
	// ResponseHeaders to be added to all responses.
	ResponseHeaders map[string]string `yaml:"responseHeaders"`
}

// Load config file from path.
func Load(path string) (File, error) {
	var file File

	data, err := os.ReadFile(path)
	if err != nil {
		return file, fmt.Errorf("failed to load file: %w", err)
	}

	err = yaml.Unmarshal(data, &file)
	if err != nil {
		return file, fmt.Errorf("failed to marshal data: %w", err)
	}

	return file, nil
}
