package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

// ServerConfig is used to tune server behaviour
type ServerConfig struct {
	// StoreInterval interval in seconds of saving metrics to disk. use 0 to write immediately
	// FileStorage absolute path to json file for db to be saved into. on omission, don't write
	// Restore on true, loads db from file on start
	// PayloadSignatureKey used to check authenticity and to sign response hashes
	Addr                string
	DBDsn               string
	StoreInterval       int
	FileStorage         string
	Restore             bool
	PayloadSignatureKey string
	PrivateKeyFilePath  string
}

func New() ServerConfig {
	return *getConfig()
}

func getConfig() *ServerConfig {
	if strings.HasSuffix(os.Args[0], ".test") {
		return &ServerConfig{}
	}
	var addressFlag *string
	var storeIntervalFlag *int
	var fileStorageFlag *string
	var restoreFlag *bool
	var DBDsnFlag *string
	var payloadSignatureKeyFlag *string
	var privateKeyFlag *string

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
	if flag.Lookup("-crypto-key") == nil {
		privateKeyFlag = flag.String(
			"-crypto-key",
			"",
			"path to private key file used to decrypt response",
		)
	}

	flag.Parse()

	flags := ServerConfig{
		Addr:                *addressFlag,
		StoreInterval:       *storeIntervalFlag,
		FileStorage:         *fileStorageFlag,
		Restore:             *restoreFlag,
		DBDsn:               *DBDsnFlag,
		PayloadSignatureKey: *payloadSignatureKeyFlag,
		PrivateKeyFilePath:  *privateKeyFlag,
	}

	overwriteWithEnv(&flags)

	return &flags
}

func overwriteWithEnv(flags *ServerConfig) {
	flags.StoreInterval = getStoreInterval(flags.StoreInterval)
	flags.Restore = getRestore(flags.Restore)

	flags.Addr = getStringFromEnvOrFallback("ADDRESS", flags.Addr)
	flags.FileStorage = getStringFromEnvOrFallback("FILE_STORAGE_PATH", flags.FileStorage)
	flags.DBDsn = getStringFromEnvOrFallback("DATABASE_DSN", flags.DBDsn)
	flags.PayloadSignatureKey = getStringFromEnvOrFallback("KEY", flags.PayloadSignatureKey)
	flags.PrivateKeyFilePath = getStringFromEnvOrFallback("CRYPTO_KEY", flags.PrivateKeyFilePath)
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
