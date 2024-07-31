package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/gennadyterekhov/metrics-storage/internal/agent"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

func getConfig() *agent.Config {
	var publicKeyFlag *string

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
	if flag.Lookup("-crypto-key") == nil {
		publicKeyFlag = flag.String(
			"-crypto-key",
			"",
			"path to public key file used to encrypt request",
		)
	}
	flag.Parse()

	flags := agent.Config{
		Addr:                      *addressFlag,
		IsGzip:                    *gzipFlag,
		ReportInterval:            *reportIntervalFlag,
		PollInterval:              *pollIntervalFlag,
		PayloadSignatureKey:       *payloadSignatureKeyFlag,
		SimultaneousRequestsLimit: *simultaneousRequestsLimitFlag,
		IsBatch:                   true,
		PublicKeyFilePath:         *publicKeyFlag,
	}

	overwriteWithEnv(&flags)

	if flags.SimultaneousRequestsLimit < 1 {
		logger.Custom.Infoln("limit flag < 1, setting to 1")
		flags.SimultaneousRequestsLimit = 1
	}

	return &flags
}

func overwriteWithEnv(flags *agent.Config) {
	flags.IsGzip = isGzip(flags.IsGzip)
	flags.ReportInterval = getReportInterval(flags.ReportInterval)
	flags.PollInterval = getPollInterval(flags.PollInterval)
	flags.SimultaneousRequestsLimit = getSimultaneousRequestsLimit(flags.SimultaneousRequestsLimit)

	flags.PayloadSignatureKey = getStringFromEnvOrFallback("KEY", flags.PayloadSignatureKey)
	flags.Addr = getStringFromEnvOrFallback("ADDRESS", flags.Addr)
	flags.PublicKeyFilePath = getStringFromEnvOrFallback("CRYPTO_KEY", flags.PublicKeyFilePath)
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

func getStringFromEnvOrFallback(envKey string, fallback string) string {
	fromEnv, ok := os.LookupEnv(envKey)
	if ok {
		return fromEnv
	}

	return fallback
}
