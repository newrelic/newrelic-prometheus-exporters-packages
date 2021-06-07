#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/githubactions"

repository_path="https://github.com/Spendesk/github-actions-exporter"
exporter_version="1.2.6"

tmp_dir=$(mktemp -d)
git clone ${repository_path} "${tmp_dir}"
cd "${tmp_dir}"
git checkout ${exporter_version}
GOOS=linux GOARCH=amd64 go build -v -o "${integration_dir}/target/bin/githubactions-exporter"
rm -rf "${tmp_dir}"