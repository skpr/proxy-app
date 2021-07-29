Reverse Proxy App
=======================

A lightweight, re-usable & configurable proxy app.

## What problem are we solving?

Skpr currently has no way to proxy one end-point to another.

## Usage

The proxy should not need to be configured, but you can map a couple of 
configurable values to the Skpr configuration using environment variables or
command line arguments. This will override any Skpr config values - and allow
the proxy to be run without Skpr config being utilized.

The environment variable will instruct the proxy which Skpr config key/value
pair to use, for example `PROXY_APP_CONFIG_KEY_ENDPOINT` could be set to
`elasticsearch.default.endpoint`.

**Note**: username and password are not being utilised by the application yet, 
but configuration and value setting is set up and ready to go.

|          | Environment Variable            | Command line argument |
|----------|---------------------------------|-----------------------|
| Username | `PROXY_APP_CONFIG_KEY_USERNAME` | `--username`          |
| Password | `PROXY_APP_CONFIG_KEY_PASSWORD` | `--password`          |
| Endpoint | `PROXY_APP_CONFIG_KEY_ENDPOINT` | `--url`               |
| Port     | `PROXY_APP_PORT`                | `--port`              |

### How to

#### From Command line

You can run this like any Go file/binary with the specified command line
arguments or environment variables.

**Method 1**: Compile/Run
```shell
$ go run main.go --port 7000 --url https://www.google.com
$ # go build -o ./proxy-app . && chmod +x ./proxy-app
```

**Method 2**: Run
```shell
$ ./proxy-app --port 7000 --url https://www.google.com
```

#### From Docker

You can build the dockerfile using the following, as well as using the build
arguments ass listed above.
```shell
$ docker build --build-arg PROXY_APP_PORT=7000 --build-arg PROXY_APP_CONFIG_KEY_ENDPOINT=https://www.google.com -t skpr/proxy-app .
$ docker run -p 7000:7000 skpr/proxy-app
```

#### How do I know it's running?

Well, following the above arguments or command line flags, you can curl 
`http://localhost:7000` and retrieve a 404 which will use the host header
as `http://localhost:7000` but will proxy through to `https://www.google.com`
which does not have content for the full url `http://localhost:7000`.

**Example**:
```shell
$ curl --HEAD -L http://localhost:7000
HTTP/1.1 404 Not Found
Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000,h3-T051=":443"; ma=2592000,h3-Q050=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000,quic=":443"; ma=2592000; v="46,43"
Content-Length: 1561
Content-Type: text/html; charset=UTF-8
Date: Thu, 29 Jul 2021 00:23:58 GMT
Referrer-Policy: no-referrer
```

## Licence

This project is licenced under GPLv3