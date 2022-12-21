#!/bin/bash
set -eo pipefail

goreleaser_bin=${GORELEASER_BIN:-goreleaser}

root_dir=$1
integration=$2
goos=$3

integration_dir="${root_dir}/exporters/${integration}"
exporters_doc_dir="${integration_dir}/target"
goreleaser_file="${root_dir}/scripts/pkg/.goreleaser.yml"
goreleaser_file_template="${goreleaser_file}.template"

echo "Loading Variables of of exporter"
source "${root_dir}"/scripts/common_functions.sh
loadVariables "${integration_dir}"/exporter.yml

echo "Package configurator"

echo "Downloading the license and placing it under ${exporters_doc_dir}/${integration}-LICENSE"
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

echo "Packaging"
if [ "$goos" == "windows" ]; then
  GORELEASER_CURRENT_TAG=${VERSION} ${goreleaser_bin} build --config ${goreleaser_file} --snapshot --rm-dist --id ${goos}
  powershell.exe -file "${root_dir}/scripts/win_msi_build.ps1" -arch amd64 -exporterName ${NAME} -version ${VERSION} -exporterGUID ${EXPORTER_GUID} -upgradeGUID ${UPGRADE_GUID} -licenseGUID ${LICENSE_GUID}
else
   cp "${goreleaser_file_template}" "${goreleaser_file}"
   IFS=',' read -r -a goarchs <<< "$PACKAGE_LINUX_GOARCHS"
   for goarch in "${goarchs[@]}"
   do
     echo  "Adding ${goarch} to goreleaser"
     yq e -i ".builds[0].goarch += [ \"${goarch}\" ]" "${goreleaser_file}"
   done

   GORELEASER_CURRENT_TAG=${VERSION} ${goreleaser_bin} release --config "${goreleaser_file}" --rm-dist --snapshot --id ${goos}
   echo "Signing the packages"
   bash ${root_dir}/scripts/sign.sh "${root_dir}" "${integration}"
fi
