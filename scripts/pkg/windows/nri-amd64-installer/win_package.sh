#!/bin/bash

PLATFORM="x64"
PFX_PATH='mycert.pfx'

# Parse arguments
for i in "$@"
do
case $i in
    -exporterName=*)
    INTEGRATION_NAME="${i#*=}"
    ;;
    -version=*)
    INTEGRATION_VERSION="${i#*=}"
    ;;
    -upgradeGUID=*)
    UPGRADE_CODE="${i#*=}"
    ;;
    -pfx_passphrase=*)
    PFX_PASSPHRASE="${i#*=}"
    ;;
esac
done

# verify version number format
if ! [[ $INTEGRATION_VERSION =~ ^[0-9]+((\.[0-9]+)){2}$ ]]; then
  echo "-version must follow a numeric major.minor.patch semantic versioning schema (received: $INTEGRATION_VERSION)"
  exit
fi

# Remove binaries
rm -rf bin obj

MSBuild="'$(find /c/Windows/Microsoft.NET/Framework64/ -path "*/v4*/MSBuild.exe")'"
SIGN_TOOL="'$(find /c/Program*Files*86*/ -path "*/SignTool/signtool.exe")'"
BINARIES_PATH="'..\..\..\..\exporters\\${INTEGRATION_NAME}\target\bin\windows_amd64\'"

eval $MSBuild nri-installer.wixproj \
    -p:IntegrationName=$INTEGRATION_NAME \
    -p:Platform=$PLATFORM \
    -p:IntegrationVersion=$INTEGRATION_VERSION \
    -p:BinariesPath=$BINARIES_PATH \
    -p:UpgradeCode=$UPGRADE_CODE

eval "$SIGN_TOOL" sign -fd SHA256 -f $PFX_PATH -p "$PFX_PASSPHRASE" bin/x64/nri-"$INTEGRATION_NAME"-amd64.msi
