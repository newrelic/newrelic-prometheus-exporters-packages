#!/bin/bash
set -e

$(/usr/local/prometheus-exporters/bin/nri-redis -short_running -config_path=/usr/local/prometheus-exporters/bin/redis-config.yml)