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
		"[address] Net address host:port without protocol",
	)
	gzipFlag := flag.Bool(
		"g",
		true,
		"[gzip] use gzip",
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
		Addr:           *addressFlag,
		IsGzip:         *gzipFlag,
		ReportInterval: *reportIntervalFlag,
		PollInterval:   *pollIntervalFlag,
	}

	overwriteWithEnv(&flags)

	return &flags
}

func overwriteWithEnv(flags *agent.AgentConfig) {
	flags.Addr = getAddress(flags.Addr)
	flags.IsGzip = isGzip(flags.IsGzip)
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

func isGzip(gzip bool) bool {
	fromEnv, ok := os.LookupEnv("GZIP")
	if ok && (fromEnv == "true" || fromEnv == "TRUE" || fromEnv == "True" || fromEnv == "1") {
		return true
	}
	if ok {
		return false
	}

	return gzip
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
