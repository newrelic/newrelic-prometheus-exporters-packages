#!/bin/bash

set -euo pipefail

# ###############################################################
#  Integration variables
root_dir=$1
integration_dir="${root_dir}/exporters/mongodb"
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

# Since we are using the make build from the repository, the binary doesn't have the correct extension
# 'exe' needed in windows executables, so we add it to the final binary.
cp "${tmp_dir}/mongodb_exporter" "${integration_bin_dir}/mongodb-exporter.exe"
