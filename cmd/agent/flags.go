package main

import (
	"flag"
	"github.com/gennadyterekhov/metrics-storage/internal/cliargs"
)

func parseFlags() (string, int, int) {
	netAddressFlag := new(cliargs.NetAddress)
	_ = flag.Value(netAddressFlag)
	flag.Var(netAddressFlag, "a", "Net address host:port")

	reportIntervalFlag := flag.Int(
		//&reportIntervalFlag,
		"r",
		//1,
		10,
		"[report interval] interval of reporting metrics to server, in seconds",
	)
	pollIntervalFlag := flag.Int(
		//&pollIntervalFlag,
		"p",
		//1,
		2,
		"[poll interval] interval of polling metrics from runtime package, in seconds",
	)

	flag.Parse()

	return netAddressFlag.String(), *reportIntervalFlag, *pollIntervalFlag
}
