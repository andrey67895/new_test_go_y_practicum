package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

var HostAgent string
var ReportIntervalAgent int
var PollIntervalAgent int
var HashKeyAgent string

func InitAgentConfig() {
	flag.StringVar(&HostAgent, "a", "localhost:8080", "HostServer for server")
	flag.IntVar(&ReportIntervalAgent, "r", 10, "reportInterval for send metrics to server")
	flag.IntVar(&PollIntervalAgent, "p", 2, "pollInterval for update metrics")
	flag.StringVar(&HashKeyAgent, "k", "123", "Key for hash")
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
	if envHashKey := os.Getenv("KEY"); envHashKey != "" {
		HashKeyAgent = envHashKey
	}
}

func getValueInEnv(env string) int {
	envInt, err := strconv.Atoi(env)
	if err != nil {
		log.Fatal(err)
	}
	return envInt
}
