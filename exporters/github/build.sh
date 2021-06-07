#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/github"

source ${root_dir}/scripts/common_functions.sh
EXPORTER_PATH="${integration_dir}/exporter.yml"
loadVariables

tmp_dir=$(mktemp -d)
git clone ${EXPORTER_REPO_URL} "${tmp_dir}"
cd "${tmp_dir}"
git checkout ${VERSION}
GOOS=linux GOARCH=amd64 go build -v -o "${integration_dir}/target/bin/github-exporter"
rm -rf "${tmp_dir}"