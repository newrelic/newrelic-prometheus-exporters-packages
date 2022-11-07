#!/bin/bash

WORKING_DIR=$1
INTEGRATION_NAME=$2
INTEGRATION_TARGET_DIR="$WORKING_DIR/exporters/$INTEGRATION_NAME/target"

cd $WORKING_DIR
source scripts/common_functions.sh
export EXPORTER_PATH="exporters/$INTEGRATION_NAME/exporter.yml"
loadVariables

TMP_DIR=$(mktemp -d)

FILENAME="${EXPORTER_LICENSE_PATH}"
git clone "${EXPORTER_REPO_URL}" "${TMP_DIR}"
TMP_LICENSE_FILE="${TMP_DIR}/${FILENAME}"

if [ ! -f "${TMP_LICENSE_FILE}" ]; then
    echo "Cannot find a LICENSE file called ${FILENAME} in the exporter repo (${TMP_LICENSE_FILE})";
    exit 1;
fi

cp "${TMP_LICENSE_FILE}" "${INTEGRATION_TARGET_DIR}/${INTEGRATION_NAME}-LICENSE"
rm -rf "${TMP_DIR}"
