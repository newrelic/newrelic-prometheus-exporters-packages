<#
    .SYNOPSIS
        This script creates the win .MSI
#>
param (
    # Target architecture: amd64 (default) or 386
    [ValidateSet("amd64", "386")]
    [string]$arch="amd64",
    [string]$exporterName="",
    [string]$exporterURL="",
    [string]$exporterHead="",
    [string]$dependencyManager
)

$env:GOPATH = go env GOPATH
$env:GOOS = "windows"
$env:GOARCH = $arch
$env:GO111MODULE = "auto"

$exporterBinaryName = "$exporterName-exporter.exe"
$exporterRepo =  [string]"$exporterURL" -replace 'https?://(www.)?'

$defRepoURL= "https://github.com/newrelic/nr-integration-definitions"
$defFileName = "prometheus_$exporterName.yml"
$defRepoPath = "definitions"
$defFilePath = "$defRepoPath\\definitions\\prometheus_exporters"

$projectRootPath = pwd

echo "--- Cloning definitions files Repo"
$ErrorActionPreference = "SilentlyContinue"
git clone $defRepoURL $defRepoPath
$ErrorActionPreference = "Stop"
$defFileExists = Test-Path $defFilePath\\$defFileName
if ($defFileExists -eq $False)
{
    echo "Cannot find a definition file called $defFileName in the definitions repo"
    exit -1
}

echo "--- Cloning exporter Repo"
Push-Location $env:GOPATH\src
$ErrorActionPreference = "SilentlyContinue"
git clone $exporterURL $exporterRepo
$ErrorActionPreference = "Stop"
Set-Location "$env:GOPATH\src\$exporterRepo"


$ErrorActionPreference = "SilentlyContinue"
git fetch -at
git checkout "$exporterHead"
$ErrorActionPreference = "Stop"

echo "--- Downloading dependencies"
$ErrorActionPreference = "SilentlyContinue"
if ($dependencyManager -eq "modules"){
    go mod download
    echo "using mod"
} elseif ($dependencyManager -eq "dep"){
    echo "using dep"
    dep ensure
}

$ErrorActionPreference = "Stop"


echo "--- Compiling exporter"
go build -v -o $exporterBinaryName
if (-not $?)
{
    echo "Failed building exporter"
    exit -1
}

Pop-Location
New-item -type directory -path .\exporters\$exporterName\target\bin\windows_$arch\ -Force
Copy-Item "$env:GOPATH\src\$exporterRepo\$exporterBinaryName" -Destination ".\exporters\$exporterName\target\bin\windows_$arch\" -Force 
if (-not $?)
{
    echo "Failed building exporter"
    exit -1
}
Copy-Item ".\exporters\$exporterName\target\source\usr\local\share\doc\prometheus-exporters\$exporterName-LICENSE" -Destination ".\exporters\$exporterName\target\bin\windows_$arch\$exporterName-LICENSE" -Force
if (-not $?)
{
    echo "Failed copying license"
    exit -1
}
Copy-Item ".\exporters\$exporterName\$exporterName-exporter.yml.sample" -Destination ".\exporters\$exporterName\target\bin\windows_$arch\$exporterName-exporter.yml.sample" -Force
if (-not $?)
{
    echo "Failed copying config file"
    exit -1
}

Copy-Item ".\exporters\$exporterName\target\source\etc\newrelic-infra\definition-files\$defFileName" -Destination ".\exporters\$exporterName\target\bin\windows_$arch\$defFileName" -Force
if (-not $?)
{
    echo "Failed copying definition file"
    exit -1
}