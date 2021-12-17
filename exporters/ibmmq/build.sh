#!/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/ibmmq"
ibmmq_bin_dir="${integration_dir}/target/bin"

tmp_dir=$(mktemp -d)
git clone ${EXPORTER_REPO_URL} "${tmp_dir}"
cd "${tmp_dir}"

if [[ -z $EXPORTER_TAG ]]
then
    git -c advice.detachedHead=false checkout ${EXPORTER_COMMIT}
else
    git checkout ${EXPORTER_TAG} -d
fi

cd scripts
HOME=${tmp_dir} MONITORS=mq_prometheus ./BuildMonitors.sh

mkdir -p ${ibmmq_bin_dir}

cp "${tmp_dir}/tmp/mq-metric-samples/bin/mq_prometheus" "${ibmmq_bin_dir}/ibmmq-exporter"
