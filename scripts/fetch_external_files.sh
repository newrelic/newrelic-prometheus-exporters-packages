#!/bin/bash
root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"

target_dir="${integration_dir}/target"
source_dir="${target_dir}/source"

exporters_doc_dir="${source_dir}/usr/local/share/doc/prometheus-exporters"
definition_files_dir="${source_dir}/etc/newrelic-infra/definition-files"

integration_license="${integration_dir}/LICENSE"

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


fetch_license