#!/bin/bash
set -eo pipefail

root_dir=$1
integration=$2

integration_dir="${root_dir}/exporters/${integration}"
publish_schema_tmp="${root_dir}/scripts/pkg/s3-publish-schema-tmp.yml"

echo "Loading Variables of the exporter"
source "${root_dir}"/scripts/common_functions.sh
loadVariables "${integration_dir}"/exporter.yml

rm -f ${publish_schema_tmp} || true
if [ "$PACKAGE_LINUX" = "true" ];then
   echo  "Adding linux to ${publish_schema_tmp}"
   cat ${root_dir}/scripts/pkg/s3-publish-schema-linux.yml >> ${publish_schema_tmp}

   yq e -i ".[].arch=[]" ${publish_schema_tmp}
   IFS=',' read -r -a goarchs <<< "$PACKAGE_LINUX_GOARCHS"
   for goarch in "${goarchs[@]}"
   do
     echo  "Adding ${goarch} to linux in ${publish_schema_tmp}"
     yq e -i ".[].arch += [ \"${goarch}\" ]" ${publish_schema_tmp}
   done
fi

if [ "$PACKAGE_WINDOWS" = "true" ];then
  echo  "Adding windows to ${publish_schema_tmp}"
  cat ${root_dir}/scripts/pkg/s3-publish-schema-windows.yml >> ${publish_schema_tmp}
fi
