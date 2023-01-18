#!/bin/bash

set -euo pipefail

# ###############################################################
#  Integration variables
root_dir=$1
integration_dir="${root_dir}/exporters/mongodb3"
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
make build

IFS=',' read -r -a goarchs <<< "$PACKAGE_LINUX_GOARCHS"
for goarch in "${goarchs[@]}"
do
  echo  "Build exporter Linux ${goarch}"
  GOARCH=${goarch} make build
  mkdir -p "${integration_bin_dir}/linux_${goarch}"
  cp "${tmp_dir}/mongodb_exporter" "${integration_bin_dir}/linux_${goarch}/mongodb-exporter"
done
