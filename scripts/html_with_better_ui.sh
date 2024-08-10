go test -v -coverpkg=./... -coverprofile coverage/coverage.out ./... > /dev/null
cat coverage/coverage.out | grep -v ".pb.go" > coverage/coverage.nopb.out
# must have go 1.23 to use
covreport -i coverage/coverage.nopb.out -o coverage/cover.html -cutlines 70,40
php -S localhost:8080
#http://localhost:8080/coverage/cover.html