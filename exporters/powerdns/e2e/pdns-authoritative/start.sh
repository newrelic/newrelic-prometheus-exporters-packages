#!/bin/bash

set -euo pipefail

DATABASE_PATH=/data/db.sqlite

mkdir -pv "$(dirname "$DATABASE_PATH")"

if [ -f "$DATABASE_PATH" ]; then
  echo Removing a previously unclean exit...
  rm -v "$DATABASE_PATH"
fi

echo Populating seed database...
sqlite3 "$DATABASE_PATH" < /usr/local/share/pdns/schema.sqlite3.sql

echo "Starting PowerDNS..."

if [ "$#" -gt 0 ]; then
  exec /usr/local/sbin/pdns_server "$@"
else
  exec /usr/local/sbin/pdns_server --daemon=no --guardian=no --control-console --loglevel=9
fi
