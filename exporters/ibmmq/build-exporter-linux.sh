#!/bin/bash

set -euo pipefail

# ###############################################################
#  Integration variables
root_dir=$1
integration_dir="${root_dir}/exporters/ibmmq"
integration_bin_dir="${integration_dir}/target/bin"

# ###############################################################
#  Clone exporter
tmp_dir=$(mktemp -d)
git clone --no-checkout "${EXPORTER_REPO_URL}" "${tmp_dir}"
cd "${tmp_dir}"

if [[ -z "${EXPORTER_TAG}" ]]
then
    git -c advice.detachedHead=false checkout "${EXPORTER_COMMIT}"
else
    git checkout "${EXPORTER_TAG}" -d
fi

# ###############################################################
#  Build exporter
cd scripts

## We are not considering GOARCH since the biulding script is not taking it into account


IFS=',' read -r -a goarchs <<< "$PACKAGE_LINUX_GOARCHS"
for goarch in "${goarchs[@]}"
do
  DOCKER_DEFAULT_PLATFORM=linux/${goarch} HOME=${tmp_dir} MONITORS=mq_prometheus ./buildMonitors.sh
  mkdir -p "${integration_bin_dir}/linux_${goarch}"
  cp "${tmp_dir}/tmp/mq-metric-samples/bin/mq_prometheus" "${integration_bin_dir}/linux_${goarch}/ibmmq-exporter"
done
