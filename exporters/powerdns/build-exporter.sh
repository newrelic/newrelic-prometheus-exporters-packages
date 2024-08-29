#!/bin/bash

set -euo pipefail

# ###############################################################
#  Integration variables
os="$1"
root_dir="$2"
integration_dir="${root_dir}/exporters/powerdns"
integration_bin_dir="${integration_dir}/target/bin"
download_url="https://github.com/newrelic-forks/prometheus-powerdns_exporter/releases/download"

# ###############################################################
#  Validate variables
if ! [ -z "$EXPORTER_COMMIT" ]; then
  echo "This script does not support to build by commit anymore."
  exit 1
fi

if [ -z "$EXPORTER_TAG" ]; then
  echo "This script requires .exporter_tag to download the exporter."
  exit 2
fi

if [ "$os" != "linux" ] && [ "$os" != "windows" ]; then
  echo "os not supported: $os."
  exit 3
fi

# ###############################################################
#  Download the exporter
IFS=',' read -r -a goarchs <<< "$PACKAGE_LINUX_GOARCHS"
for goarch in "${goarchs[@]}"
do
  echo  "Downloading exporter for ${os} ${goarch}"
  mkdir -p "${integration_bin_dir}/${os}_${goarch}"
  curl -qsLo "${integration_bin_dir}/${os}_${goarch}/powerdns-exporter" "${download_url}/${EXPORTER_TAG}/prometheus-powerdns_exporter_${EXPORTER_TAG}_${os}_${goarch}"
done
