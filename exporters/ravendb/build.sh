#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/ravendb"
GOPATH=$(go env GOPATH)

if [ ! -d "${GOPATH}" ]; then
  echo "GOPATH is empty" ;
	exit 1 ;
fi


tmp_dir="${GOPATH}/src/github.com/marcinbudny/ravendb_exporter"
echo "cloning repository into ${tmp_dir}..."
git clone ${EXPORTER_REPO_URL} "${tmp_dir}"
cd "${tmp_dir}"
git checkout ${VERSION}
echo "checkout version ${VERSION}"
echo "building the exporter..."


GO111MODULE=off GOOS=linux GOARCH=amd64 go build -v -o "${integration_dir}/target/bin/ravendb-exporter"
echo "executable file created ${integration_dir}/target/bin/ravendb-exporter"



rm -rf "${tmp_dir}"