#!/bin/bash

root_dir=$1
integration=$2
integration_dir="${root_dir}/exporters/${integration}"

rm -r "${integration_dir}/target"
echo "delete directory ${integration_dir}/target"
