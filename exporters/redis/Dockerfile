FROM alpine:3.16

RUN apk update && apk add bash

COPY target/bin /usr/local/prometheus-exporters/bin/
COPY entrypoint.sh /

RUN chmod a+x /entrypoint.sh && chmod a+x /usr/local/prometheus-exporters/bin/*

ENTRYPOINT ["/entrypoint.sh"]
