#!/bin/bash

set -euo pipefail

# ###############################################################
#  Integration variables
root_dir=$1
integration_dir="${root_dir}/exporters/powerdns"
integration_bin_dir="${integration_dir}/target/bin"

# ###############################################################
#  Clone exporter
tmp_dir=$(mktemp -d)
git clone --no-checkout ${EXPORTER_REPO_URL} "${tmp_dir}"
cd "${tmp_dir}"

if [[ -z $EXPORTER_TAG ]]
then
    git -c advice.detachedHead=false checkout ${EXPORTER_COMMIT}
else
    git checkout ${EXPORTER_TAG} -d
fi

# ###############################################################
#  Build exporter
make build

# ###############################################################
#  Move binary to its final path to be copied from the next step
mkdir -p ${integration_bin_dir}

cp "${tmp_dir}/powerdns_exporter" "${integration_bin_dir}/powerdns-exporter"
