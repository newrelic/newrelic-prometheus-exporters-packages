#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/githubactions"

tmp_dir=$(mktemp -d)
git clone ${EXPORTER_REPO_URL} "${tmp_dir}"
cd "${tmp_dir}"
git checkout ${VERSION}
make
rm -rf "${tmp_dir}"