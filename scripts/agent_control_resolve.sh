#!/usr/bin/env bash
# Resolves what is needed to publish a released exporter to the Agent Control catalog from its release tag.
#
# Usage: agent_control_resolve.sh <tag>
#   tag: release tag of the exporter, nri-<exporter>-<version> (e.g. nri-mongodb3-3.1.8)
#
# The exporter is onboarded when exporters/<exporter>/agent-control.yml exists. The outputs are written to
# $GITHUB_OUTPUT, or to the standard output when it is not set:
#   publish           true if the exporter has to be published, false otherwise
#   name, version     exporter name and version, both extracted from the tag
#   agent_type, display_name
#   oci_registry      registry the packages are uploaded to
#   config_directory  folder with the agent type files
#   artifacts         JSON array with the release assets to repackage, see agent_control_repackage.sh
set -euo pipefail

tag=${1:?usage: $0 <tag>}
root_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
output=${GITHUB_OUTPUT:-/dev/stdout}

# Same convention as the pre-release workflow: RELEASE_TAG is nri-${NAME}-${VERSION}.
regex='^nri-(.+)-([0-9]+\.[0-9]+\.[0-9]+)$'
if [[ ! ${tag} =~ ${regex} ]]; then
  echo "::error::Release tag '${tag}' does not match nri-<exporter>-<major>.<minor>.<patch>"
  exit 1
fi
name=${BASH_REMATCH[1]}
version=${BASH_REMATCH[2]}

exporter_file="${root_dir}/exporters/${name}/exporter.yml"
if [ ! -f "${exporter_file}" ]; then
  echo "::error::Exporter '${name}' (from tag '${tag}') not found: ${exporter_file} does not exist"
  exit 1
fi

agent_control_file="${root_dir}/exporters/${name}/agent-control.yml"
if [ ! -f "${agent_control_file}" ]; then
  echo "::notice::Exporter '${name}' is not onboarded to Agent Control (exporters/${name}/agent-control.yml does not exist), nothing to publish"
  echo "publish=false" >> "${output}"
  exit 0
fi

agent_type=$(yq e '.agent_type // ""' "${agent_control_file}")
display_name=$(yq e '.display_name // ""' "${agent_control_file}")
if [ -z "${agent_type}" ] || [ -z "${display_name}" ]; then
  echo "::error::exporters/${name}/agent-control.yml must define agent_type and display_name"
  exit 1
fi

package_linux=$(yq e '.package_linux // false' "${exporter_file}")
goarchs=$(yq e '.package_linux_goarchs // "amd64"' "${exporter_file}")
package_windows=$(yq e '.package_windows // false' "${exporter_file}")

if [ "${package_windows}" = "true" ]; then
  echo "::warning::Windows packages of '${name}' are not published to Agent Control: the zip does not contain the exporter binary"
fi
if [ "${package_linux}" != "true" ]; then
  echo "::notice::Exporter '${name}' has no Linux package, nothing to publish"
  echo "publish=false" >> "${output}"
  exit 0
fi

artifacts='[]'
IFS=',' read -r -a archs <<< "${goarchs}"
for arch in "${archs[@]}"; do
  artifacts=$(jq -c \
    --arg name "${name}" --arg version "${version}" --arg arch "${arch}" \
    '. + [{
      name: ("linux-" + $arch),
      os: "linux",
      arch: $arch,
      format: "tar+gzip",
      source_asset: ("nri-" + $name + "_linux_" + $version + "_" + $arch + ".tar.gz"),
      source_format: "tar.gz",
      keep_files: [
        ("var/db/newrelic-infra/newrelic-integrations/bin/nri-" + $name),
        ("usr/local/prometheus-exporters/bin/" + $name + "-exporter")
      ]
    }]' <<< "${artifacts}")
done

{
  echo "publish=true"
  echo "name=${name}"
  echo "version=${version}"
  echo "agent_type=${agent_type}"
  echo "display_name=${display_name}"
  echo "oci_registry=docker.io/newrelic/nri-${name}-artifacts"
  echo "config_directory=exporters/${name}/.fleetControl"
  echo "artifacts=${artifacts}"
} >> "${output}"
