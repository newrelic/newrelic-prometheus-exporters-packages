#!/bin/bash
root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"

target_dir="${integration_dir}/target"
source_dir="${target_dir}/source"
packages_dir="${target_dir}/packages"

PROJECT_NAME="nri-${integration}"
LICENSE="https://newrelic.com/terms (also see LICENSE.txt installed with this package)"
VENDOR="New Relic, Inc."
PACKAGER="New Relic, Inc."
PACKAGE_URL="https://github.com/newrelic/newrelic-prometheus-exporters-packages"
RELEASE=1
DESCRIPTION="Prometheus exporters help exporting existing metrics from third-party systems as Prometheus metrics."
SUMMARY="Prometheus exporter for ${integration} ${EXPORTER_REPO_URL}"
GOARCH=amd64
MINAGENTVER="newrelic-infra >= 1.25.0"

create_deb()  {
  echo "creating DEB package..."
  mkdir -p "${packages_dir}"
  fpm --verbose -C ${source_dir} -s dir -n ${PROJECT_NAME} -v ${VERSION} --iteration ${RELEASE} --prefix "" --license "${LICENSE}" --vendor "${VENDOR}" -m "${PACKAGER}" --url "${PACKAGE_URL}" --config-files /etc/newrelic-infra/ --description "${DESCRIPTION}" -t deb -p "${packages_dir}/" --depends ${MINAGENTVER} .
}

create_rpm() {
  echo "creating RPM package..."
  fpm --verbose -C ${source_dir} -s dir -n ${PROJECT_NAME} -v ${VERSION} --iteration ${RELEASE} --prefix "" --license "${LICENSE}" --vendor "${VENDOR}" -m "${PACKAGER}" --url "${PACKAGE_URL}" --config-files /etc/newrelic-infra/ --description "${DESCRIPTION}" -t rpm -p "${packages_dir}/" --epoch 0 --rpm-summary "${SUMMARY}" --depends ${MINAGENTVER} .
  echo "fpm --verbose -C ${source_dir} -s dir -n ${PROJECT_NAME} -v ${VERSION} --iteration ${RELEASE} --prefix \"\" --license \"${LICENSE}\" --vendor \"${VENDOR}\" -m \"${PACKAGER}\" --url \"${PACKAGE_URL}\" --config-files /etc/newrelic-infra/ --description \"${DESCRIPTION}\" -t rpm -p \"${packages_dir}/\" --epoch 0 --rpm-summary \"${SUMMARY}\" ."

}

create_tarball() {
  echo "creating tarball..."
  tarball_filename="${PROJECT_NAME}_linux_${VERSION}_${GOARCH}.tar.gz"
	tar -czf "${packages_dir}/${tarball_filename}" -C "${source_dir}" ./
}

sign_rpm() {
  echo "===> Create .rpmmacros to sign rpm's from Goreleaser"
  echo "%_gpg_name ${GPG_MAIL}" >> ~/.rpmmacros
  echo "%_signature gpg" >> ~/.rpmmacros
  echo "%_gpg_path ~/.gnupg" >> ~/.rpmmacros
  echo "%_gpgbin /usr/bin/gpg" >> ~/.rpmmacros
  echo "%__gpg_sign_cmd   %{__gpg} gpg --no-verbose --no-armor --batch --pinentry-mode loopback --passphrase ${GPG_PASSPHRASE} --no-secmem-warning -u "%{_gpg_name}" -sbo %{__signature_filename} %{__plaintext_filename}" >> ~/.rpmmacros

  find ${packages_dir} -regex ".*\.\(rpm\)" | while read rpm_file; do
    echo "===> Signing $rpm_file"
    rpm --addsign "$rpm_file"
    echo "===> Sign verification $rpm_file"
    rpm -v --checksig $rpm_file
  done
}

sign_deb() {
  GNUPGHOME="${HOME}/.gnupg"
  echo "${GPG_PASSPHRASE}" > ${GNUPGHOME}/gpg-passphrase
  echo "passphrase-file ${GNUPGHOME}/gpg-passphrase" >> $GNUPGHOME/gpg.conf
  echo 'allow-loopback-pinentry' >> ${GNUPGHOME}/gpg-agent.conf
  echo 'pinentry-mode loopback' >> ${GNUPGHOME}/gpg.conf
  echo 'use-agent' >> ${GNUPGHOME}/gpg.conf
  echo RELOADAGENT | gpg-connect-agent

  find ${packages_dir} -regex ".*\.\(deb\)" | while read deb_file; do
    echo "===> Signing $deb_file"
    debsigs --sign=origin --verify --check -v -k ${GPG_MAIL} $deb_file
  done
}

mkdir -p "${packages_dir}"
create_deb
create_rpm
create_tarball

if [ -z "$GPG_PRIVATE_KEY_BASE64" ];then
    echo "GPG_PRIVATE_KEY_BASE64 env variable missing package are not signed";
    exit 1;
fi
echo "===> Importing GPG private key from GHA secrets..."
printf %s ${GPG_PRIVATE_KEY_BASE64} | base64 -d | gpg --batch --import -
sign_rpm
sign_deb