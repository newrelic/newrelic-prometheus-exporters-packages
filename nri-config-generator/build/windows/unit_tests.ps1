echo "--- Running tests"
ls
go test ./...
if (-not $?)
{
    echo "Failed running tests"
    exit -1
}