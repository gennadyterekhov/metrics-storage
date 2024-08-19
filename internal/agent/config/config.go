package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/generics"

	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/iohelpler"

	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

type Config struct {
	Addr                      string `json:"address"`
	IsGzip                    bool
	ReportInterval            int `json:"report_interval"`
	PollInterval              int `json:"poll_interval"`
	IsBatch                   bool
	PayloadSignatureKey       string
	SimultaneousRequestsLimit int
	PublicKeyFilePath         string `json:"crypto_key"`
}

type cliFlags struct {
	Addr                      string
	IsGzip                    bool
	ReportInterval            int
	PollInterval              int
	IsBatch                   bool
	PayloadSignatureKey       string
	SimultaneousRequestsLimit int
	PublicKeyFilePath         string
	ConfigFilePath            string
}

// Init gets config from these places, each overwriting the previous one
// - config file (path taken from CONFIG env var or -config flag)
// - cli flags
// - env vars
func Init() *Config {
	CLIFlags := declareCLIFlags()

	resultConfig := getConfigFromFile(CLIFlags.ConfigFilePath)

	resultConfig = overwriteWithFlags(resultConfig, CLIFlags)

	resultConfig = overwriteWithEnv(resultConfig)

	overwriteWithEnv(resultConfig)

	if resultConfig.SimultaneousRequestsLimit < 1 {
		logger.Custom.Infoln("limit flag < 1, setting to 1")
		resultConfig.SimultaneousRequestsLimit = 1
	}

	return resultConfig
}

func declareCLIFlags() *cliFlags {
	var publicKeyFlag string
	var configFilePathFlag string
	var addressFlag string
	var gzipFlag bool
	var reportIntervalFlag int
	var pollIntervalFlag int
	var payloadSignatureKeyFlag string
	var simultaneousRequestsLimitFlag int

	if flag.Lookup("a") == nil {
		flag.StringVar(
			&addressFlag,
			"a",
			"localhost:8080",
			"[address] Net address host:port without protocol",
		)
	}
	if flag.Lookup("g") == nil {
		flag.BoolVar(
			&gzipFlag,
			"g",
			true,
			"[gzip] use gzip",
		)
	}
	if flag.Lookup("r") == nil {
		flag.IntVar(
			&reportIntervalFlag,
			"r",
			10,
			"[report interval] interval of reporting metrics to server, in seconds",
		)
	}
	if flag.Lookup("p") == nil {
		flag.IntVar(
			&pollIntervalFlag,
			"p",
			2,
			"[poll interval] interval of polling metrics from runtime package, in seconds",
		)
	}
	if flag.Lookup("k") == nil {
		flag.StringVar(
			&payloadSignatureKeyFlag,
			"k",
			"",
			"[key] used to sign requests' bodies so that the server can check authenticity",
		)
	}
	if flag.Lookup("l") == nil {
		flag.IntVar(
			&simultaneousRequestsLimitFlag,
			"l",
			5,
			"[limit] used to limit the number of simultaneous requests sent to server",
		)
	}

	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(
			&publicKeyFlag,
			"crypto-key",
			"",
			"path to public key file used to encrypt request",
		)
	}
	if flag.Lookup("c") == nil && flag.Lookup("config") == nil {
		flag.StringVar(&configFilePathFlag, "c", "", "path to config file")
		flag.StringVar(&configFilePathFlag, "config", "", "path to config file")
	}

	flag.Parse()

	flags := &cliFlags{
		Addr:                      addressFlag,
		IsGzip:                    gzipFlag,
		ReportInterval:            reportIntervalFlag,
		PollInterval:              pollIntervalFlag,
		PayloadSignatureKey:       payloadSignatureKeyFlag,
		SimultaneousRequestsLimit: simultaneousRequestsLimitFlag,
		IsBatch:                   true,
		PublicKeyFilePath:         publicKeyFlag,
		ConfigFilePath:            configFilePathFlag,
	}

	return flags
}

func getConfigFromFile(configFilePathFlag string) *Config {
	config := &Config{}
	configFilePath := getStringFromEnvOrFallback("CONFIG", configFilePathFlag)

	if configFilePath == "" {
		return config
	}

	configBytes, err := iohelpler.GetFileContents(configFilePath)
	if err != nil {
		logger.Custom.Panicln(errors.Wrap(err, "config file is supplied but could not be read").Error())
	}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		logger.Custom.Panicln(errors.Wrap(err, "config file is supplied but could not be decoded").Error())
	}

	return config
}

func overwriteWithFlags(resultConfig *Config, CLIFlags *cliFlags) *Config {
	resultConfig.IsGzip = generics.Overwrite(resultConfig.IsGzip, CLIFlags.IsGzip)
	resultConfig.ReportInterval = generics.Overwrite(resultConfig.ReportInterval, CLIFlags.ReportInterval)
	resultConfig.PollInterval = generics.Overwrite(resultConfig.PollInterval, CLIFlags.PollInterval)
	resultConfig.SimultaneousRequestsLimit = generics.Overwrite(
		resultConfig.SimultaneousRequestsLimit,
		CLIFlags.SimultaneousRequestsLimit,
	)

	resultConfig.PayloadSignatureKey = generics.Overwrite(
		resultConfig.PayloadSignatureKey,
		CLIFlags.PayloadSignatureKey,
	)
	resultConfig.Addr = generics.Overwrite(resultConfig.Addr, CLIFlags.Addr)
	resultConfig.PublicKeyFilePath = generics.Overwrite(resultConfig.PublicKeyFilePath, CLIFlags.PublicKeyFilePath)

	return resultConfig
}

func overwriteWithEnv(resultConfig *Config) *Config {
	resultConfig.IsGzip = isGzip(resultConfig.IsGzip)
	resultConfig.ReportInterval = getReportInterval(resultConfig.ReportInterval)
	resultConfig.PollInterval = getPollInterval(resultConfig.PollInterval)
	resultConfig.SimultaneousRequestsLimit = getSimultaneousRequestsLimit(resultConfig.SimultaneousRequestsLimit)

	resultConfig.PayloadSignatureKey = getStringFromEnvOrFallback("KEY", resultConfig.PayloadSignatureKey)
	resultConfig.Addr = getStringFromEnvOrFallback("ADDRESS", resultConfig.Addr)
	resultConfig.PublicKeyFilePath = getStringFromEnvOrFallback("CRYPTO_KEY", resultConfig.PublicKeyFilePath)

	return resultConfig
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
