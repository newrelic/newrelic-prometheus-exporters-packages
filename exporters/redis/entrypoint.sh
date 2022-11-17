#!/bin/bash
set -e
count= $(/usr/local/prometheus-exporters/bin/nri-redis -short_running -config_path=/usr/local/prometheus-exporters/bin/redis-config.yml)
echo $count