package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

var HostServer string
var StoreIntervalServer int
var FileStoragePathServer string
var RestoreServer bool
var PortServer string

func InitServerConfig() {
	flag.StringVar(&HostServer, "a", "localhost:8080", "HostServer for server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		HostServer = envRunAddr
	}
	PortServer = strings.Split(HostServer, ":")[1]

	flag.IntVar(&StoreIntervalServer, "i", 10, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&FileStoragePathServer, "f", "tmp/metrics-db.json", "полное имя файла, куда сохраняются текущие значения ")
	flag.BoolVar(&RestoreServer, "r", true, "загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.Parse()
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		StoreIntervalServer = getValueInEnv(envStoreInterval)
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePathServer = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		RestoreServer = getBool(envRestore)
	}
}

func getBool(env string) bool {

	boolValue, err := strconv.ParseBool(env)
	if err != nil {
		log.Fatal(err)
	}
	return boolValue
}
