package main

import (
	"flag"
	"os"
)

func getAddress() string {
	rawAddress, ok := os.LookupEnv("ADDRESS")
	if ok {
		return rawAddress
	}
	addressFlag := flag.String(
		"a",
		"localhost:8080",
		"Net address host:port without protocol",
	)
	flag.Parse()

	return *addressFlag
}
