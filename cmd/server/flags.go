package main

import (
	"flag"
	"github.com/gennadyterekhov/metrics-storage/internal/cliargs"
)

func parseFlags() string {
	netAddressFlag := new(cliargs.NetAddress)
	_ = flag.Value(netAddressFlag)
	flag.Var(netAddressFlag, "a", "Net address host:port")
	flag.Parse()

	return netAddressFlag.String()
}
