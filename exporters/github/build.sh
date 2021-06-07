#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/github"

repository_path="https://github.com/infinityworks/github-exporter"
exporter_version="2.0.4"

tmp_dir=$(mktemp -d)
git clone ${repository_path} "${tmp_dir}"
cd "${tmp_dir}"
git checkout ${exporter_version}
GOOS=linux GOARCH=amd64 go build -v -o "${integration_dir}/target/bin/github-exporter"
rm -rf "${tmp_dir}"