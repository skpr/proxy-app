Reverse Proxy App
=======================

A lightweight, re-usable & configurable proxy app compatible with Skpr config.

## What problem are we solving?

Skpr otherwise has no way to proxy one end-point to another using Skpr config.

## Usage

The proxy does not need to be configured, but you can map a couple of 
configurable values to the Skpr configuration or using environment variables.

The environment variable will instruct the proxy which Skpr config key/value
pair to use, for example `PROXY_APP_ENDPOINT` could be set to
`elasticsearch.default.endpoint`.

|          | Environment Variable | Environment Variable<br />(Skpr Config Use) | Explaination |
|----------|----------------------|---------------------------------------------|----------------------------------------------------------------------------------------------------------------------|
| Username | `PROXY_APP_USERNAME` | `PROXY_APP_CONFIG_KEY_USERNAME`             | HTTP Basic Authentication username for target endpoint.                                                              |
| Password | `PROXY_APP_PASSWORD` | `PROXY_APP_CONFIG_KEY_PASSWORD`             | HTTP Basic Authentication password for target endpoint.                                                              |
| Endpoint | `PROXY_APP_ENDPOINT` | `PROXY_APP_CONFIG_KEY_ENDPOINT`             | The single endpoint you wish to proxy to which includes the schema and port. Example: `https://54.206.202.192:4040`  |
| Address  | `PROXY_APP_ADDR`     | `PROXY_APP_CONFIG_KEY_ADDR`                 | The local address you wish to access the endpoint from, expressed as a port with an appended colon. Example: `:8080` | 

### How to

#### Build from source

The command should be run from a project compatible with Skpr, so that your
Skpr configuration can get loaded in at runtime - elsewise it will depend
on environment variables.

```shell
$ go build -o ${DESINATION_DIRECTORY}/proxy-app . && \
  chmod +x ./proxy-app && \
  cd ${DESINATION_DIRECTORY} && \
  ./proxy-app
```

#### Build using Docker

You can build the dockerfile using the following, as well as using the build
arguments ass listed above. This example maps to port 7000 - the default 
is 8080. The application will serve on port 7000 of the container, and the
run command maps port 7000 on the host to port 7000 on the container.

```shell
$ docker build --build-arg PROXY_APP_ADDR=:7000 --build-arg PROXY_APP_ENDPOINT=https://www.skpr.com.au -t skpr/proxy-app .
$ docker run -p 7000:7000 skpr/proxy-app
# Or...
 docker run --rm -p 7000:7000 -e PROXY_APP_ADDR=:7000 -e PROXY_APP_ENDPOINT=https://www.skpr.com.au --name skpr-proxy skpr/proxy-app
```

#### Pull from DockerHub

Prepackaged images allow you to configure the container, and an image can
be found on [Docker Hub](https://hub.docker.com/r/skpr/proxy-app)

```shell
$ docker pull skpr/proxy-app:latest
$  docker run --rm -p 7700:7700 -e PROXY_APP_ADDR=:7700 -e PROXY_APP_ENDPOINT=https://www.skpr.com.au --name skpr-proxy skpr/proxy-app
```

#### How do I know it's running?

An endpoint may not know how to handle a request with a different Hostname
header. To test for a 200 status code, we'll include our Host header:

**Example**:
```shell
$ curl --HEAD -H "Host: www.skpr.com.au" localhost:7700
HTTP/1.1 200 OK
Age: 1
Cache-Control: public, max-age=0, must-revalidate
Content-Type: text/html; charset=UTF-8
Date: Wed, 04 Aug 2021 04:42:11 GMT
Etag: "4394817e42b2622c92718ec394e7a606-ssl"
Server: Netlify
Strict-Transport-Security: max-age=31536000
X-Nf-Request-Id: 01FC7QG2RWGWCKJRJKXH55TKQC
```

## Licence

This project is licenced under GPLv3