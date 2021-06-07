#!/bin/bash

root_dir=$1
integration=$2

definition_file="prometheus_${integration}.yml"
tmp_dir=$(mktemp -d)
git clone "https://github.com/newrelic/nr-integration-definitions" "${tmp_dir}"

tmp_definition_file="${tmp_dir}/definitions/prometheus_exporters/${definition_file}"

if [ ! -f "${tmp_definition_file}" ]; then
  echo "Cannot find a definition file called ${definition_file} in the definitions repo";
  exit 1;
fi
mkdir -p
cp "${tmp_definition_file}" "${definition_file}"
rm -rf "${tmp_dir}"