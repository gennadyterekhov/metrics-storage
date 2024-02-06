package main

import (
	"flag"
	"github.com/gennadyterekhov/metrics-storage/internal/cliargs"
	"log"
	"os"
	"strconv"
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

	return getAddress(netAddressFlag), getReportInterval(reportIntervalFlag), getPollInterval(pollIntervalFlag)
}

func getAddress(netAddressFlag *cliargs.NetAddress) string {
	rawAddress, ok := os.LookupEnv("ADDRESS")
	if ok {
		return rawAddress
	}

	return netAddressFlag.String()
}

func getReportInterval(reportIntervalFlag *int) int {
	rawInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		interval, err := strconv.Atoi(rawInterval)
		if err != nil {
			log.Fatalln("incorrect format of env var REPORT_INTERVAL")
		}
		return interval
	}

	return *reportIntervalFlag
}

func getPollInterval(pollIntervalFlag *int) int {
	rawInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		interval, err := strconv.Atoi(rawInterval)
		if err != nil {
			log.Fatalln("incorrect format of env var POLL_INTERVAL")
		}
		return interval
	}

	return *pollIntervalFlag
}
