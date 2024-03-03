package config

import (
	"flag"
	"github.com/gennadyterekhov/metrics-storage/internal/testhelper"
	"log"
	"os"
	"strconv"
)

type ServerConfig struct {
	Addr          string
	StoreInterval int
	FileStorage   string
	Restore       bool
}

var Conf *ServerConfig = getConfig()

func getConfig() *ServerConfig {
	if testhelper.IsTest() {
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
		"/tmp/metrics-db.json",
		"[file storage] absolute path to json db. on omission, dont write to db",
	)
	restoreFlag := flag.Bool(
		"r",
		true,
		"[restore] on true, loads db from file on start",
	)
	flag.Parse()

	flags := ServerConfig{
		Addr:          *addressFlag,
		StoreInterval: *storeIntervalFlag,
		FileStorage:   *fileStorageFlag,
		Restore:       *restoreFlag,
	}

	overwriteWithEnv(&flags)

	return &flags
}

func overwriteWithEnv(flags *ServerConfig) {
	flags.Addr = getAddress(flags.Addr)
	flags.StoreInterval = getStoreInterval(flags.StoreInterval)
	flags.FileStorage = getFileStorage(flags.FileStorage)
	flags.Restore = getRestore(flags.Restore)
}

func getAddress(current string) string {
	rawAddress, ok := os.LookupEnv("ADDRESS")
	if ok {
		return rawAddress
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