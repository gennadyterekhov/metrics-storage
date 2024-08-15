go test -v -coverpkg=./... -coverprofile coverage/coverage.out ./... > /dev/null
#cat coverage/coverage.out | grep -v ".pb.go" > coverage/coverage.out
go tool cover -html=coverage/coverage.out