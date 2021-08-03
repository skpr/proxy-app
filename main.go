package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

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
	ProxyAppUserName string
	ProxyAppPassword string
)

func main() {
	Run()
}

type roundTripperFunc func(*http.Request) (*http.Response, error)
func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Printf("%v", r.Header.Get("Authorization"))
	return f(r)
}

func NewRoundTripper(original http.RoundTripper) http.RoundTripper {
	if original == nil {
		original = http.DefaultTransport
	}
	return roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		request.Header.Set("Authorization", fmt.Sprintf("Basic %s:%s", ProxyAppUserName, ProxyAppPassword))
		response, err := original.RoundTrip(request)
		return response, err
	})
}


func Run() {

	// Load the config.
	skprclient, _ := skprconfig.Load()
	//if err != nil && !errors.Is(err, skprconfig.ErrNotFound) {
	//	panic(err)
	//}
	ProxyAppEndpoint := skprclient.GetWithFallback(os.Getenv(CONFIG_ENDPOINT), os.Getenv("PROXY_APP_ENDPOINT"))
	ProxyAppUserName = skprclient.GetWithFallback(os.Getenv(CONFIG_USERNAME), os.Getenv("PROXY_APP_USERNAME"))
	ProxyAppPassword = skprclient.GetWithFallback(os.Getenv(CONFIG_PASSWORD), os.Getenv("PROXY_APP_PASSWORD"))

	ProxyAppPort := os.Getenv("PROXY_APP_PORT")

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	var targets []*middleware.ProxyTarget
	for _, t := range strings.Split(ProxyAppEndpoint, ",") {
		// Setup proxy
		url, err := url.Parse(t)
		if err != nil {
			e.Logger.Fatal(err)
		}
		targets = append(targets, &middleware.ProxyTarget{
			URL: url,
		})
	}

	rrb := middleware.NewRoundRobinBalancer(targets)
	var rt http.RoundTripper
	rt = NewRoundTripper(rt)

	ProxyConfig := &middleware.ProxyConfig{
		Skipper:    nil,
		Balancer:   rrb,
		Rewrite:    nil,
		ContextKey: "target",
		Transport: rt,
	}

	// Debug messaging - also add target to LB.
	fmt.Printf("Proxy configured to use port %s\n", ProxyAppPort)
	for _, target := range targets {
		ProxyConfig.Balancer.AddTarget(target)
		fmt.Printf("Starting proxy on http://localhost%s for endpoint %s\n", ProxyAppPort, target.URL)
	}

	// Start serving the proxy as configured.
	e.Use(middleware.ProxyWithConfig(*ProxyConfig))
	e.Logger.Fatal(e.Start(ProxyAppPort))

}
