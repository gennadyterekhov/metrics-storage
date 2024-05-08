package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/gennadyterekhov/metrics-storage/internal/agent"
	"github.com/gennadyterekhov/metrics-storage/internal/logger"
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
	payloadSignatureKeyFlag := flag.String(
		"k",
		"",
		"[key] used to sign requests' bodies so that the server can check authenticity",
	)
	simultaneousRequestsLimitFlag := flag.Int(
		"l",
		5,
		"[limit] used to limit the number of simultaneous requests sent to server",
	)
	flag.Parse()

	flags := agent.AgentConfig{
		Addr:                      *addressFlag,
		IsGzip:                    *gzipFlag,
		ReportInterval:            *reportIntervalFlag,
		PollInterval:              *pollIntervalFlag,
		PayloadSignatureKey:       *payloadSignatureKeyFlag,
		SimultaneousRequestsLimit: *simultaneousRequestsLimitFlag,
		IsBatch:                   true,
	}

	overwriteWithEnv(&flags)

	if flags.SimultaneousRequestsLimit < 1 {
		logger.ZapSugarLogger.Infoln("limit flag < 1, setting to 1")
		flags.SimultaneousRequestsLimit = 1
	}

	return &flags
}

func overwriteWithEnv(flags *agent.AgentConfig) {
	flags.Addr = getAddress(flags.Addr)
	flags.IsGzip = isGzip(flags.IsGzip)
	flags.ReportInterval = getReportInterval(flags.ReportInterval)
	flags.PollInterval = getPollInterval(flags.PollInterval)
	flags.PayloadSignatureKey = getKey(flags.PayloadSignatureKey)
	flags.SimultaneousRequestsLimit = getSimultaneousRequestsLimit(flags.SimultaneousRequestsLimit)
}

func getSimultaneousRequestsLimit(current int) int {
	raw, ok := os.LookupEnv("RATE_LIMIT")
	if ok {
		val, err := strconv.Atoi(raw)
		if err != nil {
			log.Fatalln("incorrect format of env var RATE_LIMIT")
		}
		return val
	}

	return current
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

func getKey(current string) string {
	raw, ok := os.LookupEnv("KEY")
	if ok {
		return raw
	}

	return current
}
