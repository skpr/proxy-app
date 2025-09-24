FROM alpine:latest

ENV PROXY_APP_ADDR=":8080"
ENV PROXY_APP_CONFIG_FILE_PATH="/etc/skpr/proxy-app/config.yaml"

ADD example/config.yaml /etc/skpr/proxy-app/config.yaml

RUN cat $PROXY_APP_CONFIG_FILE_PATH

COPY proxy-app /usr/local/bin/proxy-app
RUN chmod +x /usr/local/bin/proxy-app

CMD ["/usr/local/bin/proxy-app"]
