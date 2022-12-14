#!/bin/bash
root_dir=$1
integration=$2
packages_dir="${root_dir}/exporters/${integration}/target/packages"

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

import_keys(){
if [ -z "$GPG_PRIVATE_KEY_BASE64" ];then
    echo "GPG_PRIVATE_KEY_BASE64 env variable missing package are not signed";
    exit 1;
fi
echo "===> Importing GPG private key from GHA secrets..."
printf %s ${GPG_PRIVATE_KEY_BASE64} | base64 -d | gpg --batch --import -
}

import_keys
sign_rpm
sign_deb
