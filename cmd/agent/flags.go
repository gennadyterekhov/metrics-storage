package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

type AgentConfig struct {
	Addr           string
	ReportInterval int
	PollInterval   int
}

func getConfig() AgentConfig {
	addressFlag := flag.String(
		"a",
		"localhost:8080",
		"Net address host:port without protocol",
	)

	reportIntervalFlag := flag.Int(
		"r",
		10,
		"[report interval] interval of reporting metrics to server, in seconds",
	)
	pollIntervalFlag := flag.Int(
		"p",
		2,
		"[poll interval] interval of polling metrics from runtime package, in seconds",
	)
	flag.Parse()

	flags := AgentConfig{
		*addressFlag,
		*reportIntervalFlag,
		*pollIntervalFlag,
	}

	overwriteWithEnv(&flags)

	return flags
}

func overwriteWithEnv(flags *AgentConfig) {
	flags.Addr = getAddress(flags.Addr)
	flags.ReportInterval = getReportInterval(flags.ReportInterval)
	flags.PollInterval = getPollInterval(flags.PollInterval)
}

func getAddress(current string) string {
	rawAddress, ok := os.LookupEnv("ADDRESS")
	if ok {
		return rawAddress
	}

	return current
}

func getReportInterval(current int) int {
	rawInterval, ok := os.LookupEnv("REPORT_INTERVAL")
	if ok {
		interval, err := strconv.Atoi(rawInterval)
		if err != nil {
			log.Fatalln("incorrect format of env var REPORT_INTERVAL")
		}
		return interval
	}

	return current
}

func getPollInterval(current int) int {
	rawInterval, ok := os.LookupEnv("POLL_INTERVAL")
	if ok {
		interval, err := strconv.Atoi(rawInterval)
		if err != nil {
			log.Fatalln("incorrect format of env var POLL_INTERVAL")
		}
		return interval
	}

	return current
}
