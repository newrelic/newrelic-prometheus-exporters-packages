FROM paritytech/gnupg:latest AS KEY_DOWNLOADER
RUN curl https://repo.powerdns.com/FD380FBB-pub.asc -o /tmp/FD380FBB-pub.asc
RUN gpg --dearmor < /tmp/FD380FBB-pub.asc > /tmp/FD380FBB-pub.gpg

FROM debian:bookworm

ENV POWERDNS_API_KEY="apisecret"

COPY pdns.list /etc/apt/sources.list.d/pdns.list
COPY pdns.preference /etc/apt/preferences.d/pdns
COPY --from=KEY_DOWNLOADER /tmp/FD380FBB-pub.gpg /etc/apt/keyrings/powerdns.gpg

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -q -y pdns-server pdns-backend-sqlite3 sqlite3 && \
    rm /etc/powerdns/pdns.d/*.conf && rm /etc/powerdns/*.conf && \
    mkdir /var/run/pdns && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN sqlite3 /tmp/db.sqlite < /usr/share/pdns-backend-sqlite3/schema/schema.sqlite3.sql

COPY pdns.conf /etc/powerdns/pdns.conf

EXPOSE 53/udp 53/tcp

CMD /usr/sbin/pdns_server --daemon=no --guardian=no --control-console --loglevel=9 --api-key=${POWERDNS_API_KEY}
