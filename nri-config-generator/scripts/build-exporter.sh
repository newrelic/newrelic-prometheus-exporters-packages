#!/bin/bash
WORKING_DIR=$1
INTEGRATION_NAME=$2
OS=$3
ARCH=$4

cd $WORKING_DIR
source scripts/common_functions.sh
export EXPORTER_PATH="exporters/$INTEGRATION_NAME/exporter.yml"
loadVariables
bash exporters/$INTEGRATION_NAME/build.sh $WORKING_DIR $OS $ARCH
