FROM golang:latest AS builder

LABEL stage=builder
RUN mkdir -p /go/src/github.com/skpr/proxy
COPY . /go/src/github.com/skpr/proxy
WORKDIR /go/src/github.com/skpr/proxy
RUN GO111MODULE=on go mod vendor
RUN GO111MODULE=on go mod verify
RUN GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o pygmy-go-linux-amd64-static .

FROM skpr/base:1.x

ARG PROXY_APP_CONFIG_KEY_USERNAME
ARG PROXY_APP_CONFIG_KEY_PASSWORD
ARG PROXY_APP_CONFIG_KEY_ENDPOINT
ARG PROXY_APP_PORT

ENV PROXY_APP_CONFIG_KEY_USERNAME ${PROXY_APP_CONFIG_KEY_USERNAME:-}
ENV PROXY_APP_CONFIG_KEY_PASSWORD ${PROXY_APP_CONFIG_KEY_PASSWORD:-}
ENV PROXY_APP_CONFIG_KEY_ENDPOINT ${PROXY_APP_CONFIG_KEY_ENDPOINT:-}
ENV PROXY_APP_PORT ${PROXY_APP_PORT:-}

RUN apk add --no-cache tini
ENTRYPOINT ["/sbin/tini", "--"]

COPY --from=builder /go/src/github.com/skpr/proxy/pygmy-go-linux-amd64-static /bin/proxy
RUN chmod +x /bin/proxy

CMD ["/bin/proxy"]
EXPOSE ${PROXY_APP_PORT}
