package service

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	pb "github.com/andrey67895/new_test_go_y_practicum/proto"
)

var log = logger.Log()

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	IStorage storage.IStorageData
}

func (s *MetricsServer) UpdateMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	var response pb.MetricsResponse
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.DataLoss, "failed to get metadata")
	}
	for s2, v := range md {
		log.Info(s2, " ", v)
	}
	if config.TrustedSubnet != "" {
		xrip := md["x-real-ip"]
		if len(xrip) == 0 {
			return nil, status.Error(codes.InvalidArgument, "missing 'X-Real-IP' header")
		} else {
			ip := net.ParseIP(xrip[0])
			ones, _ := ip.DefaultMask().Size()
			_, i, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip.To4(), ones))
			mask := i.String()
			if mask != config.TrustedSubnet {
				return nil, status.Error(codes.PermissionDenied, "deny")
			}
		}
	}
	iStorage := s.IStorage
	typeMet := req.GetType()
	nameMet := req.GetId()
	switch typeMet {
	case "gauge":
		valueMet := req.GetValue()
		tErr := iStorage.RetrySaveGauge(ctx, nameMet, valueMet)
		if tErr != nil {
			log.Error(tErr.Error())
			return nil, status.Error(codes.Unavailable, tErr.Error())
		}
	case "counter":
		valueMet := req.GetDelta()
		localCounter, tErr := iStorage.GetCounter(ctx, nameMet)
		if tErr != nil {
			ttErr := iStorage.RetrySaveCounter(ctx, nameMet, valueMet)
			if ttErr != nil {
				log.Error(ttErr.Error())
				return nil, status.Error(codes.Unavailable, ttErr.Error())
			}
		} else {
			ttErr := iStorage.RetrySaveCounter(ctx, nameMet, localCounter+valueMet)
			if ttErr != nil {
				log.Error(ttErr.Error())
				return nil, status.Error(codes.Unavailable, ttErr.Error())
			}
		}
	default:
		err := fmt.Errorf("неверный тип метрики! Допустимые значения: gauge, counter")
		log.Error(err.Error())
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	return &response, nil
}
