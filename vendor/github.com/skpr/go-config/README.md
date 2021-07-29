# skprconfig

[![CircleCI](https://circleci.com/gh/skpr/go-config.svg?style=svg)](https://circleci.com/gh/skpr/go-config)

This is a go package providing an interface to read config values on the skpr.io platform.

## Usage

```go
import "github.com/skpr/skprconfig"

// Load the config.
config, err := skprconfig.Load()
if err != nil {
  panic("failed to load config")
}

// Get a string value.
bar, ok := config.Get("token")
if !ok {
  panic("token config key not found")
}

// Get the configured value for "port", with a default fallback if missing.
listenPort := config.GetWithFallback("port", "8888")

```
