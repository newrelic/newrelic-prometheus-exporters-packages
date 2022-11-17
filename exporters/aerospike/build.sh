# !/bin/bash
root_dir=$1
integration_dir="${root_dir}/exporters/aerospike"
aerospike_bin_dir="${integration_dir}/target/bin"

tmp_dir=$(mktemp -d)
git clone ${EXPORTER_REPO_URL} "${tmp_dir}"
cd "${tmp_dir}"

if [[ -z $EXPORTER_TAG ]]
then
    git -c advice.detachedHead=false checkout ${EXPORTER_COMMIT}
else
    git checkout ${EXPORTER_TAG} -d
fi

# Prerequisites to start building the binary package
# - Go language

# create a binary file named "aerospike-prometheus-exporter" from the cloned repository
go build -o aerospike-prometheus-exporter

mkdir -p ${aerospike_bin_dir}

# copy binary file
cp "${tmp_dir}/aerospike-prometheus-exporter" "${aerospike_bin_dir}/aerospike-exporter"
# copy config file
cp "${tmp_dir}/ape.toml" "${aerospike_bin_dir}/aerospike-config.yml.sample"
