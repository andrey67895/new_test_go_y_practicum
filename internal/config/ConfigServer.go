package config

import (
	"flag"
	"fmt"
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
var DatabaseDsn string

func InitServerConfig() {

	flag.StringVar(&HostServer, "a", "localhost:8080", "HostServer for server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		HostServer = envRunAddr
	}
	PortServer = strings.Split(HostServer, ":")[1]
	flag.StringVar(&DatabaseDsn, "d", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", `localhost`, 5434, `docker`, `docker`, `postgres`), "DataBase dsn for server")
	//flag.StringVar(&DatabaseDsn, "d", "", "DataBase dsn for server")

	flag.IntVar(&StoreIntervalServer, "i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&FileStoragePathServer, "f", "tmp/metrics-db.json", "полное имя файла, куда сохраняются текущие значения ")
	flag.BoolVar(&RestoreServer, "r", true, "загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.Parse()
	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		DatabaseDsn = envDatabaseDsn
	}
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
