package main

import (
	"flag"
	"github.com/gennadyterekhov/metrics-storage/internal/agent"
	"log"
	"os"
	"strconv"
)

func getConfig() *agent.AgentConfig {
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

	flags := agent.AgentConfig{
		*addressFlag,
		*reportIntervalFlag,
		*pollIntervalFlag,
	}

	overwriteWithEnv(&flags)

	return &flags
}

func overwriteWithEnv(flags *agent.AgentConfig) {
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
