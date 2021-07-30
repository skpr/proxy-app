package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	skprconfig "github.com/skpr/go-config"
)

const (
	CONFIG_ENDPOINT = "PROXY_APP_CONFIG_KEY_ENDPOINT"
	CONFIG_USERNAME = "PROXY_APP_CONFIG_KEY_USERNAME"
	CONFIG_PASSWORD = "PROXY_APP_CONFIG_KEY_PASSWORD"
)

func main() {
	Run()
}

func Run() {

	// Load the config.
	skprclient, _ := skprconfig.Load()
	//if err != nil && !errors.Is(err, skprconfig.ErrNotFound) {
	//	panic(err)
	//}
	ProxyAppEndpoint := skprclient.GetWithFallback(os.Getenv(CONFIG_ENDPOINT), os.Getenv("PROXY_APP_ENDPOINT"))
	// @TODO.
	// ProxyAppUserName := skprclient.GetWithFallback(os.Getenv(CONFIG_USERNAME), os.Getenv("PROXY_APP_USERNAME"))
	// @TODO.
	// ProxyAppPassword := skprclient.GetWithFallback(os.Getenv(CONFIG_PASSWORD), os.Getenv("PROXY_APP_PASSWORD"))

	ProxyAppPort := os.Getenv("PROXY_APP_PORT")
	if ProxyAppPort == "" {
		ProxyAppPort = ":8080"
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Setup proxy
	url, err := url.Parse(ProxyAppEndpoint)
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL:  url,
		},
	}

	// Debug messaging.
	fmt.Printf("Proxy configured to use port %s\n", ProxyAppPort)
	for _, target := range targets {
		fmt.Printf("Starting proxy on http://localhost%s for endpoint %s\n", ProxyAppPort, target.URL)
	}

	// Start serving the proxy as configured.
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	e.Logger.Fatal(e.Start(ProxyAppPort))

}
