#!/bin/bash
set -eo pipefail

root_dir=$1
integration=$2
goos=$3
integration_dir="${root_dir}/exporters/${integration}"
destination_dir="${root_dir}/nri-config-generator/templates"


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

rm -rf "${destination_dir}/exporter-config-files/"*

IFS=',' read -r -a config_files <<< "$EXPORTER_CONFIG_FILES"
for file in "${config_files[@]}"
do
  exporter_config_filename="$(basename ${file}).tmpl"
  exporter_config_path="${integration_dir}/$exporter_config_filename"
  if [ -f "$exporter_config_path" ]; then
    echo "Copying exporter config file $exporter_config_filename"
    cp "${exporter_config_path}" "${destination_dir}/exporter-config-files/${exporter_config_filename}"
  else
    echo "ERROR: Missing exporter configuration file $exporter_config_filename"
    exit 1
  fi
done

integrationName="nri-${integration}"
if [ "$goos" == "windows" ]; then
  integrationName="nri-${integration}.exe"
fi

cd nri-config-generator && \
  BIN_PATH=${binary_dir}/${integrationName} \
  make compile 

echo "executable file was created ${binary_dir}/nri-${integration}"
