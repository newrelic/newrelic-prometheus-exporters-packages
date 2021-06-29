#!/bin/bash
root_dir=$1
integration=$2

repository_definitions="https://github.com/newrelic-experimental/entity-synthesis-definitions"

tmp_dir=$(mktemp -d)
tmp_dir2=$(mktemp -d)

git clone "${repository_definitions}" "${tmp_dir}"
echo "${DEFINITION_NAMES}"
for definition in ${DEFINITION_NAMES}; do
  OLD_IFS="$IFS"
  IFS=':' read -r definition_name definition_version <<< "${definition}"
  cd ${tmp_dir} ; git checkout ${definition_version} -b latest
  tmp_definition_file="${tmp_dir}/definitions/${definition_name}/definition.yml"
  if [ ! -f "${tmp_definition_file}" ]; then
    echo "Cannot find a definition '${definition_name}' in the entity synthesis repo";
    exit 1;
  fi
  cp "${tmp_definition_file}" "${tmp_dir2}/${definition_name}.yml"
  IFS="$OLD_IFS"
done
yq eval-all "" ${tmp_dir2}/*.yml > "${root_dir}/nri-config-generator/definitions/definitions.yml"
rm -rf "${tmp_dir} ${tmp_dir2}"