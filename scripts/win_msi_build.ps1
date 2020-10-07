<#
    .SYNOPSIS
        This script creates the win .MSI
#>
param (
    # Target architecture: amd64 (default) or 386
    [ValidateSet("amd64", "386")]
    [string]$arch="amd64",
    [string]$pfx_passphrase="none",
    [string]$exporterName="",
    [string]$exporterGUID="",
    [string]$licenseGUID="",   
    [string]$configGUID="",
    [string]$version=""
)

$executable = "$exporterName-exporter.exe"

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
Import-PfxCertificate -FilePath mycert.pfx -Password (ConvertTo-SecureString -String $pfx_passphrase -AsPlainText -Force) -CertStoreLocation Cert:\CurrentUser\My

echo "===> Show certificate installed"
Get-ChildItem -Path cert:\CurrentUser\My\

echo "===> Configuring version $version for artifacts in $exporterName"

$projectRootPath = pwd
$windows_set_version = Join-Path -Path $projectRootPath -ChildPath "\scripts\windows_set_version.ps1"
& $windows_set_version -major $v[0] -minor $v[1] -patch $v[2] -exporterName $exporterName -exporterGUID $exporterGUID -licenseGUID $licenseGUID -configGUID $configGUID

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
. $msBuild/MSBuild.exe nri-installer.wixproj

if (-not $?)
{
    echo "Failed building installer"
    Pop-Location
    exit -1
}

echo "===> Making versioned installed copy moving it to $projectRootPath\exporters\$exporterName\target\packages\msi\$exporterName-$arch.msi"
New-item -type directory -path "$projectRootPath\exporters\$exporterName\target\packages\msi" -Force
cp ".\bin\Release\$exporterName-$arch.msi" "$projectRootPath\exporters\$exporterName\target\packages\msi\$exporterName-$arch.$version.msi"

Pop-Location