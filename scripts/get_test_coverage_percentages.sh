#!/bin/bash
go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./... > /dev/null
go tool cover -func coverage.out