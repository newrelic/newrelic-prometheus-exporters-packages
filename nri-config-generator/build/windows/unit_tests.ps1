echo "--- Running tests"

go mod vendor
go test ./...
if (-not $?)
{
    echo "Failed running tests"
    exit -1
}