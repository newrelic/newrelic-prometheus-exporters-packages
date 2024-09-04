#!/bin/bash

location=$(dirname "${BASH_SOURCE[0]}")
${location}/build-exporter.sh "linux" "$@"
