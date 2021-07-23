#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/powerdns"
powerdns_bin_dir="${integration_dir}/target/bin"

pwd "${powerdns_exporter_dir}"
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

ls ${tmp_dir}

cp "${tmp_dir}/powerdns_exporter" "${powerdns_bin_dir}/powerdns-exporter"

pwd "${powerdns_bin_dir}/powerdns-exporter"


