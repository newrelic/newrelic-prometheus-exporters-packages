#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/powerdns"
tmp_dir=$(mktemp -d)
git clone ${EXPORTER_REPO_URL} "${tmp_dir}"
cd "${tmp_dir}"

if [[ -z $EXPORTER_TAG ]]
then
    git -c advice.detachedHead=false checkout ${EXPORTER_COMMIT}
else
    git checkout ${EXPORTER_TAG} -d
fi

make build

mv "${tmp_dir}/powerdns_exporter" "${integration_dir}/target/bin/powerdns-exporter"

cd ../

rm -rf "${tmp_dir}"

