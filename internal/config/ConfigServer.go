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
var DatabaseDsn string
var HashKeyServer string

func InitServerConfig() {

	flag.StringVar(&HostServer, "a", "localhost:8080", "HostServer for server")
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		HostServer = envRunAddr
	}

	flag.StringVar(&HashKeyServer, "k", "", "Key for hash")
	//flag.StringVar(&DatabaseDsn, "d", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", `localhost`, 6543, `admin`, `admin`, `hr_netology`), "DataBase dsn for server")
	flag.StringVar(&DatabaseDsn, "d", "", "DataBase dsn for server")

	flag.IntVar(&StoreIntervalServer, "i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&FileStoragePathServer, "f", "tmp/metrics-db.json", "полное имя файла, куда сохраняются текущие значения ")
	flag.BoolVar(&RestoreServer, "r", true, "загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.Parse()
	PortServer = strings.Split(HostServer, ":")[1]
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
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		HashKeyServer = envHashKey
	}
}

func getBool(env string) bool {

	boolValue, err := strconv.ParseBool(env)
	if err != nil {
		log.Fatal(err)
	}
	return boolValue
}
