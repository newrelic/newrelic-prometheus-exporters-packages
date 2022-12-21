<#
    .SYNOPSIS
        This script creates the win .MSI
#>
param (
    # Target architecture: amd64 (default) or 386
    [ValidateSet("amd64", "386")]
    [string]$arch="amd64",
    [string]$exporterName="",
    [string]$exporterGUID="",
    [string]$upgradeGUID="",
    [string]$licenseGUID="",
    [string]$configGUID="",
    [string]$version=""
)

$executable = "$exporterName-exporter.exe"

$pfxPassphrase = $env:OHAI_PFX_PASSPHRASE

# verifying version number format
$v = $version.Split(".")

if ($v.Length -ne 3) {
    echo "-version must follow a numeric major.minor.patch semantic versioning schema (received: $version)"
    exit -1
}

$wrong = $v | ? { (-Not [System.Int32]::TryParse($_, [ref]0)) -or ( $_.Length -eq 0) -or ([int]$_ -lt 0)} | % { 1 }
if ($wrong.Length  -ne 0) {
    echo "-version major, minor and patch must be valid positive integers (received: $version)"
    exit -1
}

echo "===> Import .pfx certificate from GH Secrets"
Import-PfxCertificate -FilePath mycert.pfx -Password (ConvertTo-SecureString -String $pfxPassphrase -AsPlainText -Force) -CertStoreLocation Cert:\CurrentUser\My

echo "===> Show certificate installed"
Get-ChildItem -Path cert:\CurrentUser\My\

echo "===> Configuring version $version for artifacts in $exporterName"

$projectRootPath = pwd

echo "===> Checking MSBuild.exe..."
$msBuild = (Get-ItemProperty hklm:\software\Microsoft\MSBuild\ToolsVersions\4.0).MSBuildToolsPath
if ($msBuild.Length -eq 0) {
    echo "Can't find MSBuild tool. .NET Framework 4.0.x must be installed"
    exit -1
}
echo $msBuild

$env:GOOS="windows"
$env:GOARCH=$arch

echo "===> Building Installer"
Push-Location -Path "scripts\pkg\windows\nri-$arch-installer"

$env:exporterName = $exporterName
$env:IntegrationVersion = $version
$env:UpgradeCode = $upgradeGUID
. $msBuild/MSBuild.exe nri-installer.wixproj

if (-not $?)
{
    echo "Failed building installer"
    Pop-Location
    exit -1
}

echo "===> Making versioned installed copy moving it to $projectRootPath\exporters\$exporterName\target\packages\nri-$exporterName-$arch.$version.msi"
New-item -type directory -path "$projectRootPath\exporters\$exporterName\target\packages" -Force
cp ".\bin\Release\nri-$exporterName-$arch.msi" "$projectRootPath\exporters\$exporterName\target\packages\nri-$exporterName-$arch.$version.msi"

Pop-Location
