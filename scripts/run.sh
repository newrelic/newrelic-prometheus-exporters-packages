#!/bin/bash
root_dir=$1

rm -r dist/*
mkdir dist
mkdir -p dist/etc/newrelic-infra/integrations.d
mkdir -p dist/var/db/newrelic-infra/newrelic-integrations/bin
mkdir -p dist/usr/local/prometheus-exporters/bin

integrations=(powerdns)

for integration in "${integrations[@]}"
do
  echo "build integration ${integration}"
  make "build-${integration}"
  make "package-${integration}"
  cp exporters/${integration}/target/source/etc/newrelic-infra/integrations.d/*  dist/etc/newrelic-infra/integrations.d/
  cp exporters/${integration}/target/source/var/db/newrelic-infra/newrelic-integrations/bin/*  dist/var/db/newrelic-infra/newrelic-integrations/bin/
  cp exporters/${integration}/target/source/usr/local/prometheus-exporters/bin/*  dist/usr/local/prometheus-exporters/bin/
done
