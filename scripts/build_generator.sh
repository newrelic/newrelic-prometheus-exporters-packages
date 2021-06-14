#!/bin/bash

root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"


binary_dir="${root_dir}/exporters/${integration}/target/bin"
template_name=${integration}.json.tmpl
template_path="${integration_dir}/${template_name}"


cp "${template_path}" "${root_dir}/nri-config-generator/templates/${template_name}"
cd nri-config-generator && \
  GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.integration=${integration}" \
    -o "${binary_dir}/nri-${integration}" github.com/newrelic/nri-config-generator

echo "executable file was created ${binary_dir}/nri-${integration}"