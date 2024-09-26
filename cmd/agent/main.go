package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	pb "github.com/andrey67895/new_test_go_y_practicum/proto"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

var log = logger.Log()
var metricsName = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
	"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC",
	"Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc", "RandomValue",
}

var count = model.NewCount("PollCount", 0)
var metrics = model.NewMetrics()

func main() {
	log.Infof("Build version: %s", getValueOrNA(&buildVersion))
	log.Infof("Build date: %s", getValueOrNA(&buildDate))
	log.Infof("Build commit: %s", getValueOrNA(&buildCommit))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	var wg sync.WaitGroup

	conn, err := grpc.NewClient(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewMetricsServiceClient(conn)

	config.InitAgentConfig()
	go updateMetrics(time.Duration(config.PollIntervalAgent))
	go sendMetrics(&wg, time.Duration(config.ReportIntervalAgent), client, ctx)
	server := http.Server{}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Info("got interruption signal")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Info("server shutdown returned an err: %v\n", err)
	}
	log.Info("final")
}

func getValueOrNA(value *string) string {
	if value != nil && *value != "" {
		return *value
	} else {
		return "N/A"
	}
}

func updateMetrics(pollInterval time.Duration) {
	for {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			metrics.SetDataMetricsForMap(metricsName)
		}()
		go func() {
			defer wg.Done()
			getMemByGopsutil()
		}()
		wg.Wait()
		count.UpdateCountPlusOne()
		time.Sleep(pollInterval * time.Second)
	}
}

func workerRequestJSON(client pb.MetricsServiceClient, wg *sync.WaitGroup, inCh <-chan model.JSONMetrics, outCh chan<- error, ctx context.Context) {
	defer wg.Done()
	for tModel := range inCh {
		err := retrySendRequestJSON(client, tModel, ctx)
		outCh <- err
	}
}

func sendMetrics(wg *sync.WaitGroup, pollInterval time.Duration, client pb.MetricsServiceClient, ctx context.Context) {
	for {
		wg.Add(1)
		time.Sleep(pollInterval * time.Second)

		var JSONMetricsList []model.JSONMetrics
		for k, v := range metrics.GetDataMetrics() {
			gauge := v.GetMetrics()
			JSONMetricsList = append(JSONMetricsList, model.JSONMetrics{
				ID:    k,
				MType: "gauge",
				Value: &gauge})

		}
		tCounter := model.JSONMetrics{}
		tCounter.ID = count.GetName()
		tCounter.MType = "counter"
		tCounter.SetDelta(count.GetMetrics())
		JSONMetricsList = append(JSONMetricsList, tCounter)

		inputCh := make(chan model.JSONMetrics)
		outputCh := make(chan error)
		tWG := &sync.WaitGroup{}

		go func() {
			defer close(inputCh)
			for i := range JSONMetricsList {
				inputCh <- JSONMetricsList[i]
			}
		}()
		go func() {
			for i := 0; i < config.RateLimit; i++ {
				tWG.Add(1)
				go workerRequestJSON(client, tWG, inputCh, outputCh, ctx)
			}
			tWG.Wait()
			close(outputCh)
		}()
		for res := range outputCh {
			if res != nil {
				log.Error(res.Error())
				return
			}
		}
		count.ClearCount()
		wg.Done()
	}
}

func importPublicKey() *rsa.PublicKey {
	file, err := os.ReadFile(config.CryptoKeyAgent)
	if err != nil {
		return nil
	}
	block, _ := pem.Decode(file)
	if block == nil {
		return nil
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Error("Error: ", err.Error())
		return nil
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub
	default:
		break
	}
	return nil
}

func identifyIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("failed to identify IP addresses: %w", err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.To4(), nil
		}
	}
	return nil, errors.New("ip address is not found")
}

func encrypt(msg []byte) []byte {
	if config.CryptoKeyAgent != "" {
		publicKey := importPublicKey()
		cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
		if err != nil {
			log.Error("Error:", err.Error())
		}
		return cipherText
	}
	return msg
}

func getFloat64OrNull(t *float64) float64 {
	if t != nil {
		return *t
	}
	return 0
}

func getInt64OrNull(t *int64) int64 {
	if t != nil {
		return *t
	}
	return 0
}

func sendRequestJSON(client pb.MetricsServiceClient, tJSON model.JSONMetrics, ctx context.Context) error {
	req := pb.MetricsRequest{
		Id:    tJSON.ID,
		Delta: getInt64OrNull(tJSON.Delta),
		Value: getFloat64OrNull(tJSON.Value),
		Type:  tJSON.MType,
	}
	var md metadata.MD

	ip, err := identifyIP()
	if err != nil {
		log.Error(err.Error())
	} else {
		md = metadata.New(map[string]string{"x-real-ip": ip.To4().String()})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	_, err = client.UpdateMetrics(ctx, &req, grpc.Header(&md))
	if err != nil {
		return err
	}
	return err
}

func retrySendRequestJSON(client pb.MetricsServiceClient, tJSON model.JSONMetrics, ctx context.Context) error {
	err := sendRequestJSON(client, tJSON, ctx)
	if err != nil {
		log.Error(err.Error())
		for i := 1; i <= 5; i = i + 2 {
			timer := time.NewTimer(time.Duration(i) * time.Second)
			t := <-timer.C
			log.Info(t.Local())
			err = sendRequestJSON(client, tJSON, ctx)
			if err == nil {
				break
			}
		}
	}
	return err
}

func sendHashKey(r *http.Request, data []byte) {
	if config.HashKeyAgent != "" {
		hBody := bytes.Clone(data)
		h := hmac.New(sha256.New, []byte(config.HashKeyAgent))
		h.Write(hBody)
		r.Header.Add("HashSHA256", fmt.Sprintf("%x", h.Sum(nil)))
	}
}

func getMemByGopsutil() {
	v, _ := mem.VirtualMemory()
	metrics.SetDataMetrics("TotalMemory", model.NewGauge("TotalMemory", float64(v.Total)))
	metrics.SetDataMetrics("FreeMemory", model.NewGauge("FreeMemory", float64(v.Free)))
	c, _ := cpu.Percent(0, true)
	for i, percent := range c {
		metrics.SetDataMetrics(fmt.Sprintf("CPUutilization%d", i+1), model.NewGauge(fmt.Sprintf("CPUutilization%d", i+1), percent))
	}
}
