package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/generics"
	"github.com/gennadyterekhov/metrics-storage/internal/common/helper/iohelpler"
	"github.com/gennadyterekhov/metrics-storage/internal/common/logger"
)

// ServerConfig is used to tune server behaviour
type ServerConfig struct {
	// StoreInterval interval in seconds of saving metrics to disk. use 0 to write immediately
	// FileStorage absolute path to json file for db to be saved into. on omission, don't write
	// Restore on true, loads db from file on start
	// PayloadSignatureKey used to check authenticity and to sign response hashes
	Addr                string `json:"address"`
	DBDsn               string `json:"database_dsn"`
	StoreInterval       int    `json:"store_interval"`
	FileStorage         string `json:"store_file"`
	Restore             bool   `json:"restore"`
	PayloadSignatureKey string
	PrivateKeyFilePath  string `json:"crypto_key"`
}

type cliFlags struct {
	Addr                string
	DBDsn               string
	StoreInterval       int
	FileStorage         string
	Restore             bool
	PayloadSignatureKey string
	PrivateKeyFilePath  string
	ConfigFilePath      string
}

// New gets config from these places, each overwriting the previous one
// - config file (path taken from CONFIG env var or -config flag)
// - cli flags
// - env vars
func New() *ServerConfig {
	return getConfig()
}

func getConfig() *ServerConfig {
	if strings.HasSuffix(os.Args[0], ".test") {
		return &ServerConfig{}
	}
	CLIFlags := declareCLIFlags()

	resultConfig := getConfigFromFile(CLIFlags.ConfigFilePath)

	resultConfig = overwriteWithFlags(resultConfig, CLIFlags)

	resultConfig = overwriteWithEnv(resultConfig)

	return resultConfig
}

func declareCLIFlags() *cliFlags {
	var addressFlag *string
	var storeIntervalFlag *int
	var fileStorageFlag *string
	var restoreFlag *bool
	var DBDsnFlag *string
	var payloadSignatureKeyFlag *string
	var privateKeyFlag *string
	var configFilePathFlag string

	if flag.Lookup("a") == nil {
		addressFlag = flag.String(
			"a",
			"localhost:8080",
			"[address] Net address host:port without protocol",
		)
	}
	if flag.Lookup("i") == nil {
		storeIntervalFlag = flag.Int(
			"i",
			300,
			"[store interval] interval in seconds of saving metrics to disk. use 0 to write immediately",
		)
	}
	if flag.Lookup("f") == nil {
		fileStorageFlag = flag.String(
			"f",
			"/tmp/metrics-db.json",
			"[file storage] absolute path to json file for db to be saved into. on omission, don't write",
		)
	}
	if flag.Lookup("r") == nil {
		restoreFlag = flag.Bool(
			"r",
			true,
			"[restore] on true, loads db from file on start",
		)
	}
	if flag.Lookup("d") == nil {
		DBDsnFlag = flag.String(
			"d",
			"",
			"[db dsn] format: `host=%s user=%s password=%s dbname=%s sslmode=%s` if empty, ram is used",
		)
	}
	if flag.Lookup("k") == nil {
		payloadSignatureKeyFlag = flag.String(
			"k",
			"",
			"[key] used to check authenticity (bad request on failure) and to sign response hashes",
		)
	}
	if flag.Lookup("crypto-key") == nil {
		privateKeyFlag = flag.String(
			"crypto-key",
			"",
			"path to private key file used to decrypt response",
		)
	}
	if flag.Lookup("c") == nil && flag.Lookup("config") == nil {
		flag.StringVar(&configFilePathFlag, "c", "", "path to config file")
		flag.StringVar(&configFilePathFlag, "config", "", "path to config file")
	}

	flag.Parse()

	flags := &cliFlags{
		Addr:                *addressFlag,
		StoreInterval:       *storeIntervalFlag,
		FileStorage:         *fileStorageFlag,
		Restore:             *restoreFlag,
		DBDsn:               *DBDsnFlag,
		PayloadSignatureKey: *payloadSignatureKeyFlag,
		PrivateKeyFilePath:  *privateKeyFlag,
		ConfigFilePath:      configFilePathFlag,
	}

	return flags
}

func getConfigFromFile(configFilePathFlag string) *ServerConfig {
	config := &ServerConfig{}
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

func overwriteWithFlags(resultConfig *ServerConfig, CLIFlags *cliFlags) *ServerConfig {
	resultConfig.StoreInterval = generics.Overwrite(resultConfig.StoreInterval, CLIFlags.StoreInterval)
	resultConfig.Restore = generics.Overwrite(resultConfig.Restore, CLIFlags.Restore)

	resultConfig.Addr = generics.Overwrite(resultConfig.Addr, CLIFlags.Addr)
	resultConfig.FileStorage = generics.Overwrite(resultConfig.FileStorage, CLIFlags.FileStorage)
	resultConfig.DBDsn = generics.Overwrite(resultConfig.DBDsn, CLIFlags.DBDsn)
	resultConfig.PayloadSignatureKey = generics.Overwrite(resultConfig.PayloadSignatureKey, CLIFlags.PayloadSignatureKey)
	resultConfig.PrivateKeyFilePath = generics.Overwrite(resultConfig.PrivateKeyFilePath, CLIFlags.PrivateKeyFilePath)

	return resultConfig
}

func overwriteWithEnv(resultConfig *ServerConfig) *ServerConfig {
	resultConfig.StoreInterval = getStoreInterval(resultConfig.StoreInterval)
	resultConfig.Restore = getRestore(resultConfig.Restore)

	resultConfig.Addr = getStringFromEnvOrFallback("ADDRESS", resultConfig.Addr)
	resultConfig.FileStorage = getStringFromEnvOrFallback("FILE_STORAGE_PATH", resultConfig.FileStorage)
	resultConfig.DBDsn = getStringFromEnvOrFallback("DATABASE_DSN", resultConfig.DBDsn)
	resultConfig.PayloadSignatureKey = getStringFromEnvOrFallback("KEY", resultConfig.PayloadSignatureKey)
	resultConfig.PrivateKeyFilePath = getStringFromEnvOrFallback("CRYPTO_KEY", resultConfig.PrivateKeyFilePath)

	return resultConfig
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
