package main

import (
	"crypto/subtle"
	"errors"
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

var (
	PROXY_APP_ENDPOINT string
	PROXY_APP_USERNAME string
	PROXY_APP_PASSWORD string
	PROXY_APP_PORT string
)

func init() {
	// Load the config.
	skprclient, err := skprconfig.Load()
	if err != nil && !errors.Is(err, skprconfig.ErrNotFound) {
		panic(err)
	}
	PROXY_APP_ENDPOINT = skprclient.GetWithFallback(os.Getenv(CONFIG_ENDPOINT), os.Getenv("PROXY_APP_ENDPOINT"))
	PROXY_APP_USERNAME = skprclient.GetWithFallback(os.Getenv(CONFIG_USERNAME), os.Getenv("PROXY_APP_USERNAME"))
	PROXY_APP_PASSWORD = skprclient.GetWithFallback(os.Getenv(CONFIG_PASSWORD), os.Getenv("PROXY_APP_PASSWORD"))

	if PROXY_APP_PORT = os.Getenv("PROXY_APP_PORT"); PROXY_APP_PORT == "" {
		PROXY_APP_PORT = "8080"
	}

	fmt.Println(PROXY_APP_PORT)
	fmt.Println(os.Getenv("PROXY_APP_PORT"))
}

func main() {

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Setup proxy
	url, err := url.Parse(PROXY_APP_ENDPOINT)
	if err != nil {
		e.Logger.Fatal(err)
	}
	targets := []*middleware.ProxyTarget{
		{
			URL:  url,
		},
	}

	// Debug messaging.
	fmt.Printf("Proxy configured to use port %s\n", PROXY_APP_PORT)
	for _, target := range targets {
		fmt.Printf("Starting proxy on http://localhost:%s for endpoint %s\n", PROXY_APP_PORT, target.URL)
	}

	// Handle Authentication
	if PROXY_APP_USERNAME != "" && PROXY_APP_PASSWORD != "" {
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte(PROXY_APP_USERNAME)) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte(PROXY_APP_PASSWORD)) == 1 {
				return true, nil
			}
			return false, nil
		}))
	}

	// Start serving the proxy as configured.
	e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
	localAddress := fmt.Sprintf(":%s", PROXY_APP_PORT)
	e.Logger.Fatal(e.Start(localAddress))

}
