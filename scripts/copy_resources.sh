#!/bin/bash
root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"

repository_definitions="https://github.com/newrelic/nr-integration-definitions"

target_dir="${integration_dir}/target"
binaries_dir="${target_dir}/bin"
source_dir="${target_dir}/source"

integrations_exec_dir="${source_dir}/var/db/newrelic-infra/newrelic-integrations/bin/"
exporters_exec_dir="${source_dir}/usr/local/prometheus-exporters/bin/"
integrations_config_dir="${source_dir}/etc/newrelic-infra/integrations.d"
exporters_doc_dir="${source_dir}/usr/local/share/doc/prometheus-exporters"
definition_files_dir="${source_dir}/etc/newrelic-infra/definition-files"

integration_sample="${integration_dir}/${integration}-exporter.yml.sample"
integration_license="${integration_dir}/LICENSE"

create_folders_structure() {
  rm -rf "${source_dir}";
	mkdir -p  "${exporters_exec_dir}" \
	  "${integrations_exec_dir}" \
	  "${integrations_config_dir}" \
	  "${exporters_doc_dir}" \
	  "${definition_files_dir}";
}

copy_resources() {
  echo "copying binaries..."
  cp "${binaries_dir}/${integration}-exporter" "${exporters_exec_dir}/${integration}-exporter"
	chmod 755 "${exporters_exec_dir}/${integration}-exporter"
	cp "${binaries_dir}/nri-${integration}" "${integrations_exec_dir}/nri-${integration}"
	chmod 755 "${integrations_exec_dir}/nri-${integration}"
  echo "copying samples..."
  cp "${integration_sample}" "${integrations_config_dir}"
}

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

create_folders_structure
copy_resources
fetch_definition_file
fetch_license