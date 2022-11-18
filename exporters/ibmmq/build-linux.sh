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
HOME=${tmp_dir} MONITORS=mq_prometheus ./buildMonitors.sh

# ###############################################################
#  Move binary to its final path to be copied from the next step
mkdir -p "${integration_bin_dir}"

cp "${tmp_dir}/tmp/mq-metric-samples/bin/mq_prometheus" "${integration_bin_dir}/ibmmq-exporter"
