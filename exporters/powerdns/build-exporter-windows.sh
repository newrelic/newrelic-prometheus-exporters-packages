#!/bin/bash

# This script is not currently used as we are not packaging
# this integration for Windows. It is kept for testing purposes.

location=$(dirname "${BASH_SOURCE[0]}")
${location}/build-exporter.sh "windows" "$@"
