#!/bin/bash
root_dir=$1
integration=$2

repository_definitions="https://github.com/newrelic-experimental/entity-synthesis-definitions"

tmp_dir=$(mktemp -d)
tmp_dir2=$(mktemp -d)
git clone "${repository_definitions}" "${tmp_dir}"
for definition in ${DEFINITION_NAMES}; do
  definition_name=${definition%%:*}
  definition_version=${definition#*:}
  if git --git-dir ${tmp_dir}/.git show-ref --tags --quiet --verify -- "refs/tags/$definition_version" >/dev/null 2>&1; then
      cd ${tmp_dir}; git checkout ${definition_version} -d
  else
      repo_url=${definition_version%%#*}
      repo_commit=${definition#*#}
      cd ${tmp_dir}; \
        git remote add forked ${repo_url}; \
        git fetch forked;  \
        git switch ${repo_commit} -d;
  fi
  tmp_definition_file="${tmp_dir}/definitions/${definition_name}/definition.yml"
  if [ ! -f "${tmp_definition_file}" ]; then
    echo "Cannot find a definition '${definition_name}' in the entity synthesis repo";
    exit 1;
  fi
  synthesis_attr=$(eval "yq eval '.synthesis' ${tmp_definition_file}")
  if [[ "${synthesis_attr}" = "null" ]] ; then
    echo "missing required attribute 'synthesis' for definitions";
    exit 1;
  fi
  cp "${tmp_definition_file}" "${tmp_dir2}/${definition_name}.yml"
done
mkdir -p ${root_dir}/nri-config-generator/definitions
yq eval-all "" ${tmp_dir2}/*.yml > "${root_dir}/nri-config-generator/definitions/definitions.yml"
rm -rf "${tmp_dir} ${tmp_dir2}"