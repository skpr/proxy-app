package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	skprconfig "github.com/skpr/go-config"
)

const (
	CONFIG_ENDPOINT = "PROXY_APP_CONFIG_KEY_ENDPOINT"
	CONFIG_USERNAME = "PROXY_APP_CONFIG_KEY_USERNAME"
	CONFIG_PASSWORD = "PROXY_APP_CONFIG_KEY_PASSWORD"
)

var (
	username = kingpin.Flag("username", "Username").String()
	password = kingpin.Flag("password", "Password").String()
	endpoint = kingpin.Flag("url", "URL to proxy service to, which includes the port.").Short('u').String()
	port     = kingpin.Flag("port", "Port to expose the service to the host").Default("8080").Short('p').Int()
)

func main() {

	kingpin.Parse()

	// Load the config.
	config, err := skprconfig.Load()
	if err == nil {
		// If the parameters were supplied to the CLI, they should override the configuration.
		if *endpoint == "" {
			endpointEnvVar := os.Getenv(CONFIG_ENDPOINT)
			*endpoint = config.GetWithFallback(endpointEnvVar, *endpoint)
		}
		if *username == "" {
			usernameEnvVar := os.Getenv(CONFIG_USERNAME)
			*username = config.GetWithFallback(usernameEnvVar, *endpoint)
		}
		if *password == "" {
			passwordEnvVar := os.Getenv(CONFIG_PASSWORD)
			*password = config.GetWithFallback(passwordEnvVar, *endpoint)
		}
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// If the user has changed the endpoint via environment variable, respect this.
	var urlToUse string
	if *endpoint == "" {
		urlToUse = os.Getenv(CONFIG_ENDPOINT)
		if urlToUse == "" {
			urlToUse = *endpoint
		}
	} else {
		urlToUse = *endpoint
	}

	// Setup proxy
	url, err := url.Parse(urlToUse)
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL:  url,
		},
	}

	// If the user has changed the port via environment variable, respect this.
	var portToUse string
	if *port == 8080 {
		portToUse = os.Getenv("PROXY_APP_PORT")
		if portToUse == "" {
			portToUse = fmt.Sprintf("%d", *port)
		}
	} else {
		portToUse = fmt.Sprintf("%d", *port)
	}

	// Debug messaging.
	fmt.Printf("Proxy configured to use port %s\n", portToUse)
	for _, target := range targets {
		fmt.Printf("Starting proxy on http://localhost:%s for endpoint %s\n", portToUse, target.URL)
	}

	// Start serving the proxy as configured.
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	localAddress := fmt.Sprintf(":%s", portToUse)
	e.Logger.Fatal(e.Start(localAddress))

}
