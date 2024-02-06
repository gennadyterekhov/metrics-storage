package main

import (
	"flag"
	"github.com/gennadyterekhov/metrics-storage/internal/cliargs"
	"os"
)

func parseFlags() string {
	rawAddress, ok := os.LookupEnv("ADDRESS")
	if ok {
		return rawAddress
	}
	netAddressFlag := new(cliargs.NetAddress)
	_ = flag.Value(netAddressFlag)
	flag.Var(netAddressFlag, "a", "Net address host:port")
	flag.Parse()

	return netAddressFlag.String()
}
