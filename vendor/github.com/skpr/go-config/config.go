package skprconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	// DefaultPath is the default file path used to mount config.
	DefaultPath string = "/etc/skpr/data/config.json"
)

// Config represents the config.
type Config struct {
	path string
	data map[string]interface{}
}

// Load loads Config from file.
func Load(options ...func(config *Config)) (*Config, error) {
	config := &Config{
		path: DefaultPath,
	}
	for _, option := range options {
		option(config)
	}

	if _, err := os.Stat(config.path); os.IsNotExist(err) {
		return config, fmt.Errorf("config file does not exist: %w", err)
	}

	data, err := ioutil.ReadFile(config.path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	var configData map[string]interface{}

	err = json.Unmarshal(data, &configData)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	config.data = configData

	return config, nil
}

// Get returns a string value for the key.
func (c *Config) Get(key string) (string, bool) {
	value, ok := c.getValue(key)
	if value == nil {
		value = ""
	}
	return value.(string), ok
}

func (c *Config) getValue(key string) (interface{}, bool) {
	if _, ok := c.data[key]; !ok {
		return nil, false
	}
	return c.data[key], true
}

// GetWithFallback returns the configured value of a given key, and the fallback
// value if no key does not exist.
func (c *Config) GetWithFallback(key, fallback string) string {
	if _, ok := c.getValue(key); !ok {
		return fallback
	}
	value, _ := c.getValue(key)
	return value.(string)
}
