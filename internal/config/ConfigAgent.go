// Package config работа с конфиами Server и Agent
package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	HostAgent           string
	ReportIntervalAgent int
	PollIntervalAgent   int
	HashKeyAgent        string
	RateLimit           int
	CryptoKeyAgent      string
)

// InitAgentConfig Инициализация Конфигурации для агента
func InitAgentConfig() {
	flag.StringVar(&HostAgent, "a", "localhost:8080", "HostServer for server")
	flag.IntVar(&ReportIntervalAgent, "r", 10, "reportInterval for send metrics to server")
	flag.IntVar(&PollIntervalAgent, "p", 2, "pollInterval for update metrics")
	flag.IntVar(&RateLimit, "l", 9, "RateLimit for update metrics")
	flag.StringVar(&HashKeyAgent, "k", "", "Key for hash")
	flag.StringVar(&CryptoKeyAgent, "crypto-key", "", "Key for asymmetric encryption")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		HostAgent = envRunAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		ReportIntervalAgent = getValueInEnv(envReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		PollIntervalAgent = getValueInEnv(envPollInterval)
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		RateLimit = getValueInEnv(envRateLimit)
	}
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		HashKeyAgent = envHashKey
	}
	if envCryptoKeyAgent := os.Getenv("CRYPTO_KEY"); envCryptoKeyAgent != "" {
		CryptoKeyAgent = envCryptoKeyAgent
	}
}

func getValueInEnv(env string) int {
	envInt, err := strconv.Atoi(env)
	if err != nil {
		log.Fatal(err)
	}
	return envInt
}
