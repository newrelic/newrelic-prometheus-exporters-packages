#!/bin/bash
root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"

source ${root_dir}/scripts/common_functions.sh
EXPORTER_PATH="${integration_dir}/exporter.yml"
loadVariables

target_dir="${integration_dir}/target"
source_dir="${target_dir}/source"
packages_dir="${target_dir}/packages"
rpm_dir="${packages_dir}/rpm"
deb_dir="${packages_dir}/deb"
tarball_dir="${packages_dir}/tarball"

PROJECT_NAME="${integration}-exporter"
LICENSE="https://newrelic.com/terms (also see LICENSE.txt installed with this package)"
VENDOR="New Relic, Inc."
PACKAGER="New Relic, Inc."
PACKAGE_URL="https://github.com/newrelic/newrelic-prometheus-exporters-packages"
RELEASE=1
DESCRIPTION="Prometheus exporters help exporting existing metrics from third-party systems as Prometheus metrics."
SUMMARY="Prometheus exporter for ${integration} ${EXPORTER_REPO_URL}"
GOARCH=amd64




function create_deb  {
  echo "creating DEB package..."
  mkdir -p "${deb_dir}"
  fpm --verbose -C ${source_dir} -s dir -n ${PROJECT_NAME} -v ${VERSION} --iteration ${RELEASE} --prefix "" --license "${LICENSE}" --vendor "${VENDOR}" -m "${PACKAGER}" --url "${PACKAGE_URL}" --config-files /etc/newrelic-infra/ --description "${DESCRIPTION}" -t deb -p "${deb_dir}/" .
}

function create_rpm {
  echo "creating RPM package..."
  mkdir -p "${rpm_dir}"
  fpm --verbose -C ${source_dir} -s dir -n ${PROJECT_NAME} -v ${VERSION} --iteration ${RELEASE} --prefix "" --license "${LICENSE}" --vendor "${VENDOR}" -m "${PACKAGER}" --url "${PACKAGE_URL}" --config-files /etc/newrelic-infra/ --description "${DESCRIPTION}" -t rpm -p "${rpm_dir}/" --epoch 0 --rpm-summary "${SUMMARY}" .
  echo "fpm --verbose -C ${source_dir} -s dir -n ${PROJECT_NAME} -v ${VERSION} --iteration ${RELEASE} --prefix \"\" --license \"${LICENSE}\" --vendor \"${VENDOR}\" -m \"${PACKAGER}\" --url \"${PACKAGE_URL}\" --config-files /etc/newrelic-infra/ --description \"${DESCRIPTION}\" -t rpm -p \"${rpm_dir}/\" --epoch 0 --rpm-summary \"${SUMMARY}\" ."

}

function create_tarball() {
  echo "creating tarball..."
  tarball_filename="${PROJECT_NAME}_linux_${VERSION}_${GOARCH}.tar.gz"
  mkdir -p "${tarball_dir}"
  mkdir -p "${tarball_dir}"
	tar -czf "${tarball_dir}/${tarball_filename}" -C "${source_dir}" ./
}




create_deb
create_rpm
create_tarball