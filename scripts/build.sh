#!/bin/bash
set -eo pipefail

goreleaser_bin=${GORELEASER_BIN:-goreleaser}

root_dir=$1
integration=$2
goos=$3

template_name=${integration}.json.tmpl
config_name=${integration}.prometheus.json.tmpl
destination_dir="${root_dir}/nri-config-generator/templates"
integration_dir="${root_dir}/exporters/${integration}"
template_path="${integration_dir}/${integration}.json.tmpl"
config_path="${integration_dir}/${integration}.prometheus.json.tmpl"
goreleaser_file="${root_dir}/scripts/pkg/.goreleaser.yml"
goreleaser_file_template="${goreleaser_file}.template"


echo "Loading Variables of of exporter"
source "${root_dir}"/scripts/common_functions.sh
loadVariables "${integration_dir}"/exporter.yml

if [[ "$PACKAGE_LINUX" != "true" ]] && [[ "$PACKAGE_WINDOWS" != "true" ]]; then
    echo "ERROR: the exporter would not be packaged for any supported OS (at least one of package_linux or package_windows should be set to true)"
    exit 1
fi

echo "Building exporter"
bash ${root_dir}/exporters/${integration}/build-exporter-${goos}.sh ${root_dir}

echo "Finding prometheus template and placing it under nri-config-generator/templates folder"
cp "${template_path}" "${destination_dir}/${template_name}"
if [ -f "$config_path" ]; then
  echo "Using prometheus template $config_path"
  cp "${config_path}" "${destination_dir}/${config_name}"
else
  echo "Using Default prometheus template"
  cp "${destination_dir}/default/config.json.tmpl" "${destination_dir}/${config_name}"
fi

echo "Finding exporter-config-files and placing them under nri-config-generator/templates/exporter-config-files folder"
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

cp "${goreleaser_file_template}" "${goreleaser_file}"
IFS=',' read -r -a goarchs <<< "$PACKAGE_LINUX_GOARCHS"
for goarch in "${goarchs[@]}"
do
  echo  "Adding ${goarch} to goreleaser"
  yq e -i ".builds[0].goarch += [ \"${goarch}\" ]" ${goreleaser_file}
done

GORELEASER_CURRENT_TAG=${VERSION} ${goreleaser_bin} build --config ${goreleaser_file} --snapshot --rm-dist --id ${goos}
