go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./... > /dev/null
cat coverage/coverage.out | grep -v ".pb.go" > coverage/coverage.nopb.out
go tool cover -func coverage/coverage.nopb.out