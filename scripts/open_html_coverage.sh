go test -v -coverpkg=./... -coverprofile coverage/coverage.out ./... > /dev/null
go tool cover -html=coverage/coverage.out