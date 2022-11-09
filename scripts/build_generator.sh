#!/bin/bash

root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"


binary_dir="${root_dir}/exporters/${integration}/target/bin"
template_name=nri-${integration}.json.tmpl
template_path="${integration_dir}/${integration}.json.tmpl"
config_name=nri-${integration}.prometheus.json.tmpl
config_path="${integration_dir}/${integration}.prometheus.json.tmpl"

cp "${template_path}" "${root_dir}/nri-config-generator/templates/${template_name}"

if [ -f "$config_path" ]; then
  echo "Using prometheus template $config_path"
  cp "${config_path}" "${root_dir}/nri-config-generator/templates/${config_name}"
else
  echo "Using Default prometheus template"
  cp "${root_dir}/nri-config-generator/templates/default/config.json.tmpl" "${root_dir}/nri-config-generator/templates/${config_name}"
fi

cd nri-config-generator && \
  BIN_PATH=${binary_dir}/nri-${integration} \
  make compile 

echo "executable file was created ${binary_dir}/nri-${integration}"
