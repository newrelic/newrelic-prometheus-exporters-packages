#!/bin/bash

set -euo pipefail -x

# ###############################################################
#  Integration variables
root_dir=$1
integration_dir="${root_dir}/exporters/ibmmq"
integration_bin_dir="${integration_dir}/target/bin"
ibmmq_client_libs_version="9.3.4.0"

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

go mod tidy

curl -qsLO "https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist/${ibmmq_client_libs_version}-IBM-MQC-Redist-Win64.zip"
unzip "${ibmmq_client_libs_version}-IBM-MQC-Redist-Win64.zip" -d "IBM-MQC-Redist-Win64"
export CGO_CFLAGS="-I$(cygpath -aw .)\\IBM-MQC-Redist-Win64\\tools\\c\\include -D_WIN64"
export CGO_LDFLAGS="-L$(cygpath -aw .)\\IBM-MQC-Redist-Win64\\bin64 -lmqm"

go build \
  -ldflags="-X \"main.BuildStamp=$(date +%Y%m%d-%H%M%S)\" -X \"main.BuildPlatform=$(uname -s)/$(uname -i)\" -X \"main.GitCommit=$(git rev-list -1 HEAD --abbrev-commit 2>/dev/null || echo "Unknown")\"" \
  -o ibmmq-exporter.exe \
  .\\cmd\\mq_prometheus

# ###############################################################
#  Move binary to its final path to be copied from the next step
mkdir -p "${integration_bin_dir}"

cp "ibmmq-exporter.exe" "${integration_bin_dir}/ibmmq-exporter.exe"
