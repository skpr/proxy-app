package main

import (
	"errors"
	"os"

	skprconfig "github.com/skpr/go-config"
	"github.com/skpr/proxy-app/internal/server"
)

const (
	// EnvSkprConfigKeyAddr used to load the address value from Skpr config.
	EnvSkprConfigKeyAddr = "PROXY_APP_CONFIG_KEY_ADDR"
	// EnvSkprConfigKeyEndpoint used to load the endpoint value from Skpr config.
	EnvSkprConfigKeyEndpoint = "PROXY_APP_CONFIG_KEY_ENDPOINT"
	// EnvSkprConfigKeyUsername used to load the username value from Skpr config.
	EnvSkprConfigKeyUsername = "PROXY_APP_CONFIG_KEY_USERNAME"
	// EnvSkprConfigKeyPassword used to load the password value from Skpr config.
	EnvSkprConfigKeyPassword = "PROXY_APP_CONFIG_KEY_PASSWORD"

	// EnvAddr sets the address for the proxy application.
	EnvAddr = "PROXY_APP_ADDR"
	// EnvEndpoint sets the endpoint for the proxy application.
	EnvEndpoint = "PROXY_APP_ENDPOINT"
	// EnvUsername sets the username for the proxy connection.
	EnvUsername = "PROXY_APP_USERNAME"
	// EnvPassword sets the password for the proxy connection.
	EnvPassword = "PROXY_APP_PASSWORD"
)

func main() {
	skprclient, err := skprconfig.Load()
	if err != nil && !errors.Is(err, skprconfig.ErrNotFound) {
		panic(err)
	}

	params := server.RunParams{
		Addr:     skprclient.GetWithFallback(os.Getenv(EnvSkprConfigKeyAddr), os.Getenv(EnvAddr)),
		Endpoint: skprclient.GetWithFallback(os.Getenv(EnvSkprConfigKeyEndpoint), os.Getenv(EnvEndpoint)),
		Username: skprclient.GetWithFallback(os.Getenv(EnvSkprConfigKeyUsername), os.Getenv(EnvUsername)),
		Password: skprclient.GetWithFallback(os.Getenv(EnvSkprConfigKeyPassword), os.Getenv(EnvPassword)),
	}

	if err := server.Run(params); err != nil {
		panic(err)
	}
}
