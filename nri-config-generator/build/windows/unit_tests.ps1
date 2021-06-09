echo "--- Running tests"
cd nri-config-generator
go test ./...
if (-not $?)
{
    echo "Failed running tests"
    exit -1
}