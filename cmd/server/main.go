package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/interceptors"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/service"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/andrey67895/new_test_go_y_practicum/internal/transport/router"
	pb "github.com/andrey67895/new_test_go_y_practicum/proto"
)

var buildVersion string
var buildDate string
var buildCommit string

var log = logger.Log()

func main() {
	log.Infof("Build version: %s", getValueOrNA(&buildVersion))
	log.Infof("Build date: %s", getValueOrNA(&buildDate))
	log.Infof("Build commit: %s", getValueOrNA(&buildCommit))

	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.RealIpInterceptor))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	var wg sync.WaitGroup
	config.InitServerConfig()
	var st storage.IStorageData
	if config.DatabaseDsn != "" {
		st = storage.InitDB(ctx)
	} else {
		st = storage.InMemStorage{}
		if config.FileStoragePathServer != "" {
			if config.RestoreServer {
				RestoringDataFromFile(config.FileStoragePathServer)
			}
			go SaveDataForInterval(&wg, ctx, config.FileStoragePathServer, config.StoreIntervalServer)
		}
	}
	pb.RegisterMetricsServiceServer(s, &service.MetricsServer{
		IStorage: st,
	})

	server := http.Server{
		Addr:    ":" + config.PortServer,
		Handler: router.GetRoutersForServer(st),
	}
	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatal("listen and serve returned err: %v", err)
		}
	}()
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()
	<-ctx.Done()
	log.Info("got interruption signal")
	s.GracefulStop()
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

func RestoringDataFromFile(fname string) {
	data, err := os.ReadFile(fname)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var tModel []model.JSONMetrics
	if err := json.Unmarshal(data, &tModel); err != nil {
		log.Error(err.Error())
	}
	for i := 0; i < len(tModel); i++ {
		metric := tModel[i]
		SaveData(metric)
	}
}

func SaveData(tModel model.JSONMetrics) {
	typeMet := tModel.MType
	nameMet := tModel.ID

	switch typeMet {
	case "gauge":
		valueMet := tModel.GetValue()
		err := storage.LocalNewMemStorageGauge.SetGauge(nameMet, valueMet)
		if err != nil {
			log.Error(err.Error())
			return
		}
	case "counter":
		valueMet := tModel.GetDelta()
		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			err := storage.LocalNewMemStorageCounter.SetCounter(nameMet, valueMet)
			if err != nil {
				log.Error(err.Error())
				return
			} else {
				tModel.SetDelta(localCounter + valueMet)
				err = storage.LocalNewMemStorageCounter.SetCounter(nameMet, tModel.GetDelta())
				if err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	}

}

func SaveDataForInterval(wg *sync.WaitGroup, ctx context.Context, fname string, storeInterval int) {
	if storeInterval > 0 {
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		for range ticker.C {
			select {
			case <-ctx.Done():
				log.Info("Save data finish")
				ticker.Stop()
				return
			default:
				wg.Add(1)
				storage.SaveDataInFile(fname)
				log.Infoln("Save Data file at: ", time.Now())
				wg.Done()
			}
		}
	}
}
