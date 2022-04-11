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
    [string]$exporterURL="",
    [string]$exporterHead="",
    [string]$exporterGUID="",
    [string]$upgradeGUID="",
    [string]$licenseGUID="",
    [string]$configGUID="",
    [string]$version=""
)

$projectRootPath = pwd

$win_build = Join-Path -Path $projectRootPath -ChildPath "\scripts\win_exe_build.ps1"
& $win_build -arch $arch -exporterName $exporterName -exporterHead $exporterHead  -exporterURL $exporterURL -dependencyManager "dep"

$win_msi_build = Join-Path -Path $projectRootPath -ChildPath "\scripts\win_msi_build.ps1"
& $win_msi_build -arch $arch -exporterName $exporterName -version $version -exporterGUID $exporterGUID -upgradeGUID $upgradeGUID -licenseGUID $licenseGUID -configGUID $configGUID -definitionGUID $definitionGUID  -pfx_passphrase $pfx_passphrase
