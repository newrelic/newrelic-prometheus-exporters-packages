#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/ravendb"

repository_path="https://github.com/marcinbudny/ravendb_exporter"
exporter_version="0.3.0"

if [ ! -d "${GOPATH}" ]; then
  echo "GOPATH is empty" ;
	exit 1 ;
fi


tmp_dir="${GOPATH}/src/github.com/marcinbudny/ravendb_exporter"
echo "cloning repository into ${tmp_dir}..."
git clone ${repository_path} "${tmp_dir}"
cd "${tmp_dir}"
git checkout ${exporter_version}
echo "checkout version ${exporter_version}"
echo "building the exporter..."
GO111MODULE=off GOOS=linux GOARCH=amd64 go build -v -o "${integration_dir}/target/bin/ravendb-exporter"
echo "executable file created ${integration_dir}/target/bin/ravendb-exporter"
rm -rf "${tmp_dir}"