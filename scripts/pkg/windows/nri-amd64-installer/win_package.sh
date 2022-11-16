#!/bin/bash

#TODO use snakecase
INTEGRATION_NAME='powerdns'
INTEGRATION_VERSION="0.0.2"
UPGRADE_CODE="2818F0DF-0DB1-4152-B16E-88B51CA4DC9F"
BINARIES_PATH="'..\..\..\..\exporters\powerdns\target\bin\windows_amd64\'"
PLATFORM='x64'
PFX_PASS='fake'
PFX_PATH='fake.pfx'
PFX_DESCRIPTION='fake-cert.com'

MSBuild="'$(find /c/Windows/Microsoft.NET/Framework64/ -path "*/v4*/MSBuild.exe")'"
SIGN_TOOL="'$(find /c/Program*Files*86*/ -path "*/x64/signtool.exe")'"

# TODO remove
# rm -rf bin obj

eval $MSBuild nri-installer.wixproj \
    -p:IntegrationName=$INTEGRATION_NAME \
    -p:Platform=$PLATFORM \
    -p:IntegrationVersion=$INTEGRATION_VERSION \
    -p:BinariesPath=$BINARIES_PATH \
    -p:UpgradeCode=$UPGRADE_CODE

eval $SIGN_TOOL sign -fd SHA256 -f $PFX_PATH -p $PFX_PASS -d $PFX_DESCRIPTION -n $PFX_DESCRIPTION bin/x64/nri-powerdns-amd64.msi

