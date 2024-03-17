package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type ServerConfig struct {
	Addr          string
	StoreInterval int
	FileStorage   string
	Restore       bool
	DBDsn         string
}

var Conf *ServerConfig = getConfig()

func getConfig() *ServerConfig {
	if strings.HasSuffix(os.Args[0], ".test") {
		return &ServerConfig{}
	}

	addressFlag := flag.String(
		"a",
		"localhost:8080",
		"[address] Net address host:port without protocol",
	)
	storeIntervalFlag := flag.Int(
		"i",
		300,
		"[store interval] interval in seconds of saving metrics to disk. use 0 to write immediately",
	)
	fileStorageFlag := flag.String(
		"f",
		"/tmp/metrics-2_db.json",
		"[file storage] absolute path to json 2_db. on omission, dont write to 2_db",
	)
	restoreFlag := flag.Bool(
		"r",
		true,
		"[restore] on true, loads 2_db from file on start",
	)
	DBDsnFlag := flag.String(
		"d",
		"",
		"[db dsn] format: `host=%s user=%s password=%s dbname=%s sslmode=%s` if empty, ram is used",
	)
	flag.Parse()

	flags := ServerConfig{
		Addr:          *addressFlag,
		StoreInterval: *storeIntervalFlag,
		FileStorage:   *fileStorageFlag,
		Restore:       *restoreFlag,
		DBDsn:         *DBDsnFlag,
	}

	overwriteWithEnv(&flags)

	return &flags
}

func overwriteWithEnv(flags *ServerConfig) {
	flags.Addr = getAddress(flags.Addr)
	flags.StoreInterval = getStoreInterval(flags.StoreInterval)
	flags.FileStorage = getFileStorage(flags.FileStorage)
	flags.Restore = getRestore(flags.Restore)
	flags.DBDsn = getDBDsn(flags.DBDsn)
}

func getAddress(current string) string {
	rawAddress, ok := os.LookupEnv("ADDRESS")
	if ok {
		return rawAddress
	}

	return current
}

func getDBDsn(current string) string {
	raw, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		return raw
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

func getFileStorage(current string) string {
	rawInterval, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		return rawInterval
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
