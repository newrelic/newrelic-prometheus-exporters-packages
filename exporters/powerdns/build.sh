#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/powerdns"
powerdns_bin_dir="${integration_dir}/target/bin"

tmp_dir=$(mktemp -d)
git clone ${EXPORTER_REPO_URL} "${tmp_dir}"
cd "${tmp_dir}"

if [[ -z $EXPORTER_TAG ]]
then
    git -c advice.detachedHead=false checkout ${EXPORTER_COMMIT}
else
    git checkout ${EXPORTER_TAG} -d
fi

goos=$2
goarch=$3

GOOS="$goos" GOARCH="$goarch" make build

mkdir -p ${powerdns_bin_dir}

cp "${tmp_dir}/powerdns_exporter" "${powerdns_bin_dir}/powerdns-exporter"
