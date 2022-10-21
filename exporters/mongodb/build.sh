#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/mongodb"
mongodb_bin_dir="${integration_dir}/target/bin"

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

mkdir -p ${mongodb_bin_dir}

cp "${tmp_dir}/mongodb_exporter" "${mongodb_bin_dir}/mongodb-exporter"
