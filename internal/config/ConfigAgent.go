// Package config работа с конфиами Server и Agent
package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
)

var (
	HostAgent           string
	ReportIntervalAgent int
	PollIntervalAgent   int
	HashKeyAgent        string
	RateLimit           int
	CryptoKeyAgent      string
	agentConfig         string
)

// InitAgentConfig Инициализация Конфигурации для агента
func InitAgentConfig() {
	flag.StringVar(&HostAgent, "a", "localhost:8080", "HostServer for server")
	flag.IntVar(&ReportIntervalAgent, "r", 10, "reportInterval for send metrics to server")
	flag.IntVar(&PollIntervalAgent, "p", 2, "pollInterval for update metrics")
	flag.IntVar(&RateLimit, "l", 9, "RateLimit for update metrics")
	flag.StringVar(&HashKeyAgent, "k", "", "Key for hash")
	flag.StringVar(&CryptoKeyAgent, "crypto-key", "", "Key for asymmetric encryption")
	flag.StringVar(&agentConfig, "config", "", "Config json for agent")
	flag.StringVar(&agentConfig, "c", "", "Config json for agent (shorthand)")
	flag.Parse()
	if envConfigServer := os.Getenv(""); envConfigServer != "" {
		agentConfig = envConfigServer
	}
	configIt := agentConfig != ""
	var config model.ConfigAgentModel
	if configIt {
		file, err := os.ReadFile(agentConfig)
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
		HostAgent = envRunAddr
	} else if configIt {
		HostAgent = config.Address
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		ReportIntervalAgent = getValueInEnv(envReportInterval)
	} else if configIt {
		ReportIntervalAgent = getValueInEnv(config.ReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		PollIntervalAgent = getValueInEnv(envPollInterval)
	} else if configIt {
		PollIntervalAgent = getValueInEnv(config.PollInterval)
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		RateLimit = getValueInEnv(envRateLimit)
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		HashKeyAgent = envHashKey
	}
	if envCryptoKeyAgent := os.Getenv("CRYPTO_KEY"); envCryptoKeyAgent != "" {
		CryptoKeyAgent = envCryptoKeyAgent
	} else if configIt {
		CryptoKeyAgent = config.CryptoKey
	}
}

func getValueInEnv(env string) int {
	envInt, err := strconv.Atoi(env)
	if err != nil {
		log.Fatal(err)
	}
	return envInt
}
