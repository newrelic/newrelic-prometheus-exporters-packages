#!/bin/bash
root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"

repository_definitions="https://github.com/newrelic/nr-integration-definitions"

target_dir="${integration_dir}/target"
source_dir="${target_dir}/source"

exporters_doc_dir="${source_dir}/usr/local/share/doc/prometheus-exporters"
definition_files_dir="${source_dir}/etc/newrelic-infra/definition-files"

integration_license="${integration_dir}/LICENSE"



fetch_definition_file() {
  tmp_dir=$(mktemp -d)


  filename="prometheus_${integration}.yml"
  git clone "${repository_definitions}" "${tmp_dir}"
  tmp_definition_file="${tmp_dir}/definitions/prometheus_exporters/${filename}"

  if [ ! -f "${tmp_definition_file}" ]; then
    echo "Cannot find a definition file called ${filename} in the definitions repo";
    exit 1;
  fi

  cp "${tmp_definition_file}" "${definition_files_dir}/${filename}"
  rm -rf "${tmp_dir}"
}

fetch_license() {
  tmp_dir=$(mktemp -d)

  filename="${EXPORTER_LICENSE_PATH}"
  git clone "${EXPORTER_REPO_URL}" "${tmp_dir}"
  tmp_license_file="${tmp_dir}/${filename}"

   if [ ! -f "${tmp_license_file}" ]; then
    echo "Cannot find a LICENSE file called ${filename} in the exporter repo";
    exit 1;
  fi

  cp "${tmp_license_file}" "${exporters_doc_dir}/${integration}-LICENSE"
  rm -rf "${tmp_dir}"

}

fetch_definition_file
fetch_license