package config

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/generics"
	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/iohelpler"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

// ServerConfig is used to tune server behaviour
type ServerConfig struct {
	// IsGrpc on true, use grpc and don't use HTTP
	// Restore on true, loads db from file on start
	// StoreInterval interval in seconds of saving metrics to disk. use 0 to write immediately
	// FileStorage absolute path to json file for db to be saved into. on omission, don't write
	// PayloadSignatureKey used to check authenticity and to sign response hashes
	// TrustedSubnet is a single ip subnet that allowed to make requests in CIDR format
	IsGrpc              bool
	Restore             bool
	StoreInterval       int
	Addr                string
	DBDsn               string
	FileStorage         string
	PayloadSignatureKey string
	PrivateKeyFilePath  string
	TrustedSubnet       *net.IPNet
}

type cliOrJSONConfig struct {
	IsGrpc              bool   `json:"grpc"`
	Restore             bool   `json:"restore"`
	StoreInterval       int    `json:"store_interval"`
	Addr                string `json:"address"`
	DBDsn               string `json:"database_dsn"`
	FileStorage         string `json:"store_file"`
	PrivateKeyFilePath  string `json:"crypto_key"`
	TrustedSubnet       string `json:"trusted_subnet"`
	ConfigFilePath      string `json:"-"`
	PayloadSignatureKey string
}

// New gets config from these places, each overwriting the previous one
// - config file (path taken from CONFIG env var or -config flag)
// - cli flags
// - env vars
func New() *ServerConfig {
	return getConfig()
}

func getConfig() *ServerConfig {
	CLIFlags := declareCLIFlags()

	changingConfig := getConfigFromFile(CLIFlags.ConfigFilePath)

	changingConfig = overwriteWithFlags(changingConfig, CLIFlags)

	changingConfig = overwriteWithEnv(changingConfig)

	resultConfig := parseComplexTypes(changingConfig)

	return resultConfig
}

func parseComplexTypes(config *cliOrJSONConfig) *ServerConfig {
	_, subnet, err := net.ParseCIDR(config.TrustedSubnet)
	if err != nil {
		logger.Custom.Debugln("could not parse subnet.TrustedSubnet from string", err.Error())
		subnet = nil
	}
	resultConfig := &ServerConfig{
		IsGrpc:              config.IsGrpc,
		Restore:             config.Restore,
		StoreInterval:       config.StoreInterval,
		Addr:                config.Addr,
		DBDsn:               config.DBDsn,
		FileStorage:         config.FileStorage,
		PayloadSignatureKey: config.PayloadSignatureKey,
		PrivateKeyFilePath:  config.PrivateKeyFilePath,
		TrustedSubnet:       subnet,
	}
	return resultConfig
}

func declareCLIFlags() *cliOrJSONConfig {
	var isGrpcFlag bool
	var addressFlag string
	var storeIntervalFlag int
	var fileStorageFlag string
	var restoreFlag bool
	var DBDsnFlag string
	var payloadSignatureKeyFlag string
	var privateKeyFlag string
	var configFilePathFlag string
	var trustedSubnetFlag string

	if flag.Lookup("grpc") == nil {
		flag.BoolVar(&isGrpcFlag,
			"grpc",
			false,
			"true to use grpc, false to use HTTP",
		)
	}
	if flag.Lookup("a") == nil {
		flag.StringVar(&addressFlag,
			"a",
			"localhost:8080",
			"[address] Net address host:port without protocol",
		)
	}
	if flag.Lookup("i") == nil {
		flag.IntVar(&storeIntervalFlag,
			"i",
			300,
			"[store interval] interval in seconds of saving metrics to disk. use 0 to write immediately",
		)
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&fileStorageFlag,
			"f",
			"/tmp/metrics-db.json",
			"[file storage] absolute path to json file for db to be saved into. on omission, don't write",
		)
	}
	if flag.Lookup("r") == nil {
		flag.BoolVar(&restoreFlag,
			"r",
			true,
			"[restore] on true, loads db from file on start",
		)
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&DBDsnFlag,
			"d",
			"",
			"[db dsn] format: `host=%s user=%s password=%s dbname=%s sslmode=%s` if empty, ram is used",
		)
	}
	if flag.Lookup("k") == nil {
		flag.StringVar(&payloadSignatureKeyFlag,
			"k",
			"",
			"[key] used to check authenticity (bad request on failure) and to sign response hashes",
		)
	}
	if flag.Lookup("crypto-key") == nil {
		flag.StringVar(&privateKeyFlag,
			"crypto-key",
			"",
			"path to private key file used to decrypt response",
		)
	}
	if flag.Lookup("t") == nil {
		flag.StringVar(&trustedSubnetFlag,
			"t",
			"",
			"a single ip subnet that allowed to make requests",
		)
	}
	if flag.Lookup("c") == nil && flag.Lookup("config") == nil {
		flag.StringVar(&configFilePathFlag, "c", "", "path to config file")
		flag.StringVar(&configFilePathFlag, "config", "", "path to config file")
	}

	flag.Parse()

	flags := &cliOrJSONConfig{
		IsGrpc:              isGrpcFlag,
		Addr:                addressFlag,
		StoreInterval:       storeIntervalFlag,
		FileStorage:         fileStorageFlag,
		Restore:             restoreFlag,
		DBDsn:               DBDsnFlag,
		PayloadSignatureKey: payloadSignatureKeyFlag,
		PrivateKeyFilePath:  privateKeyFlag,
		ConfigFilePath:      configFilePathFlag,
		TrustedSubnet:       trustedSubnetFlag,
	}

	return flags
}

func getConfigFromFile(configFilePathFlag string) *cliOrJSONConfig {
	config := &cliOrJSONConfig{}
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

func overwriteWithFlags(resultConfig *cliOrJSONConfig, CLIFlags *cliOrJSONConfig) *cliOrJSONConfig {
	resultConfig.StoreInterval = generics.Overwrite(resultConfig.StoreInterval, CLIFlags.StoreInterval)
	resultConfig.Restore = generics.Overwrite(resultConfig.Restore, CLIFlags.Restore)

	resultConfig.IsGrpc = generics.Overwrite(resultConfig.IsGrpc, CLIFlags.IsGrpc)
	resultConfig.Addr = generics.Overwrite(resultConfig.Addr, CLIFlags.Addr)
	resultConfig.FileStorage = generics.Overwrite(resultConfig.FileStorage, CLIFlags.FileStorage)
	resultConfig.DBDsn = generics.Overwrite(resultConfig.DBDsn, CLIFlags.DBDsn)
	resultConfig.PayloadSignatureKey = generics.Overwrite(resultConfig.PayloadSignatureKey, CLIFlags.PayloadSignatureKey)
	resultConfig.PrivateKeyFilePath = generics.Overwrite(resultConfig.PrivateKeyFilePath, CLIFlags.PrivateKeyFilePath)
	resultConfig.TrustedSubnet = generics.Overwrite(resultConfig.PrivateKeyFilePath, CLIFlags.TrustedSubnet)

	return resultConfig
}

func overwriteWithEnv(resultConfig *cliOrJSONConfig) *cliOrJSONConfig {
	resultConfig.StoreInterval = getStoreInterval(resultConfig.StoreInterval)
	resultConfig.Restore = getRestore(resultConfig.Restore)

	resultConfig.IsGrpc = getIsGrpc(resultConfig.IsGrpc)
	resultConfig.Addr = getStringFromEnvOrFallback("ADDRESS", resultConfig.Addr)
	resultConfig.FileStorage = getStringFromEnvOrFallback("FILE_STORAGE_PATH", resultConfig.FileStorage)
	resultConfig.DBDsn = getStringFromEnvOrFallback("DATABASE_DSN", resultConfig.DBDsn)
	resultConfig.PayloadSignatureKey = getStringFromEnvOrFallback("KEY", resultConfig.PayloadSignatureKey)
	resultConfig.PrivateKeyFilePath = getStringFromEnvOrFallback("CRYPTO_KEY", resultConfig.PrivateKeyFilePath)
	resultConfig.TrustedSubnet = getStringFromEnvOrFallback("TRUSTED_SUBNET", resultConfig.TrustedSubnet)

	return resultConfig
}

func getIsGrpc(current bool) bool {
	fromEnv, ok := os.LookupEnv("USE_GRPC")
	if ok && (fromEnv == "true" || fromEnv == "TRUE" || fromEnv == "True" || fromEnv == "1") {
		return true
	}
	if ok {
		return false
	}

	return current
}

func getStoreInterval(current int) int {
	fromEnv, ok := os.LookupEnv("STORE_INTERVAL")
	if ok {
		interval, err := strconv.Atoi(fromEnv)
		if err != nil {
			log.Fatalln("incorrect format of env var STORE_INTERVAL")
		}
		return interval
	}

	return current
}

func getRestore(current bool) bool {
	fromEnv, ok := os.LookupEnv("RESTORE")
	if ok && (fromEnv == "true" || fromEnv == "TRUE" || fromEnv == "True" || fromEnv == "1") {
		return true
	}
	if ok {
		return false
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
