package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
)

var (
	HostServer            string
	StoreIntervalServer   int
	FileStoragePathServer string
	RestoreServer         bool
	PortServer            string
	DatabaseDsn           string
	HashKeyServer         string
	CryptoKeyServer       string
	serverConfig          string
)

var log = logger.Log()

// InitServerConfig Инициализация Конфигурации для сервера
func InitServerConfig() {
	flag.StringVar(&HostServer, "a", "localhost:8080", "HostServer for server")

	flag.StringVar(&HashKeyServer, "k", "", "Key for hash")
	//flag.StringVar(&DatabaseDsn, "d", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", `localhost`, 6543, `admin`, `admin`, `hr_netology`), "DataBase dsn for server")
	flag.StringVar(&DatabaseDsn, "d", "", "DataBase dsn for server")

	flag.IntVar(&StoreIntervalServer, "i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&FileStoragePathServer, "f", "tmp/metrics-db.json", "полное имя файла, куда сохраняются текущие значения ")
	flag.BoolVar(&RestoreServer, "r", true, "загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&CryptoKeyServer, "crypto-key", "", "Key for asymmetric encryption")
	flag.StringVar(&serverConfig, "config", "", "Config json for server")
	flag.StringVar(&serverConfig, "c", "", "Config json for server (shorthand)")

	flag.Parse()
	if envConfigServer := os.Getenv(""); envConfigServer != "" {
		serverConfig = envConfigServer
	}
	configIt := serverConfig != ""
	var config model.ConfigServerModel
	if configIt {
		file, err := os.ReadFile(serverConfig)
		if err != nil {
			log.Error("Ошибка чтения конфига json: ", err)
			configIt = false
		}
		err = json.Unmarshal(file, &config)
		if err != nil {
			log.Error("Ошибка формирования конфига json: ", err)
			configIt = false
		}
	}
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		HostServer = envRunAddr
	} else if configIt {
		HostServer = config.Address
	}
	PortServer = strings.Split(HostServer, ":")[1]

	if envDatabaseDsn := os.Getenv("DATABASE_DSN"); envDatabaseDsn != "" {
		DatabaseDsn = envDatabaseDsn
	} else if configIt {
		DatabaseDsn = config.DatabaseDsn
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		StoreIntervalServer = getValueInEnv(envStoreInterval)
	} else if configIt {
		StoreIntervalServer = getValueInEnv(config.StoreInterval)
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		FileStoragePathServer = envFileStoragePath
	} else if configIt {
		FileStoragePathServer = config.StoreFile
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		RestoreServer = getBool(envRestore)
	} else if configIt {
		RestoreServer = config.Restore
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		HashKeyServer = envHashKey
	}
	if envCryptoKeyServer := os.Getenv("CRYPTO_KEY"); envCryptoKeyServer != "" {
		CryptoKeyServer = envCryptoKeyServer
	} else if configIt {
		CryptoKeyServer = config.CryptoKey
	}
}

func getBool(env string) bool {

	boolValue, err := strconv.ParseBool(env)
	if err != nil {
		log.Fatal(err)
	}
	return boolValue
}
