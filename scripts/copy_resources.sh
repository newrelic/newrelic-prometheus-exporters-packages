#!/bin/bash
root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"

repository_definitions="https://github.com/newrelic/nr-integration-definitions"

target_dir="${integration_dir}/target"
binaries_dir="${target_dir}/bin"
source_dir="${target_dir}/source"

integrations_exec_dir="${source_dir}/var/db/newrelic-infra/newrelic-integrations/bin"
exporters_exec_dir="${source_dir}/usr/local/prometheus-exporters/bin"
integrations_config_dir="${source_dir}/etc/newrelic-infra/integrations.d"

integration_sample="${integration_dir}/${integration}-config.yml.sample"

copy_resources() {
  echo "copying binaries..."
  cp "${binaries_dir}/${integration}-exporter" "${exporters_exec_dir}/${integration}-exporter"
	chmod 755 "${exporters_exec_dir}/${integration}-exporter"
	cp "${binaries_dir}/nri-${integration}" "${integrations_exec_dir}/nri-${integration}"
	chmod 755 "${integrations_exec_dir}/nri-${integration}"
  echo "copying samples..."
  cp "${integration_sample}" "${integrations_config_dir}"
}

copy_resources
