Reverse Proxy App
=======================

A lightweight, re-usable & configurable proxy app compatible with Skpr config.

## What problem are we solving?

Skpr currently has no way to proxy one end-point to another using Skpr config.

## Usage

The proxy should not need to be configured, but you can map a couple of 
configurable values to the Skpr configuration or using environment variables.

The environment variable will instruct the proxy which Skpr config key/value
pair to use, for example `PROXY_APP_ENDPOINT` could be set to
`elasticsearch.default.endpoint`.

|          | Environment Variable |
|----------|----------------------|
| Username | `PROXY_APP_USERNAME` |
| Password | `PROXY_APP_PASSWORD` |
| Endpoint | `PROXY_APP_ENDPOINT` |
| Port     | `PROXY_APP_PORT`     |

### How to

#### From Skpr config

In this example, our arguments will come from the Skpr config.

The command should be run from a project compatible with Skpr.

```shell
$ go build -o ${DESINATION_DIRECTORY}/proxy-app . && \
  chmod +x ./proxy-app && \
  cd ${DESINATION_DIRECTORY} && \
  ./proxy-app
```

#### From Docker

You can build the dockerfile using the following, as well as using the build
arguments ass listed above. This example maps to port 7000 - the default 
is 8080. The application will serve on port 7000 of the container, and the
run command maps port 7000 on the host to port 7000 on the container.

```shell
$ docker build --build-arg PROXY_APP_PORT=:7000 --build-arg PROXY_APP_ENDPOINT=https://www.google.com -t skpr/proxy-app .
$ docker run -p 7000:7000 skpr/proxy-app
```

#### How do I know it's running?

Note that if you did not configure the port to run on, it will run on `8080`.
Following the above build/run commands, you can curl `http://localhost:7000` 
and retrieve a 404 which will use the host header as `http://localhost:7000`
but will proxy through to `https://www.google.com`. This endpoint does not
have content for the hostname `http://localhost:7000`.

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