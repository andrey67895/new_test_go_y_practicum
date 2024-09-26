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

func (s *MetricsServer) GetDataByTypeAndName(ctx context.Context, req *pb.GetDataByTypeAndNameRequest) (*pb.GetDataByTypeAndNameResponse, error) {
	if err := checkIP(ctx); err != nil {
		return nil, err
	}
	result := pb.GetDataByTypeAndNameResponse{}
	result.Id = req.GetId()
	result.Type = req.GetType()
	switch req.GetType() {
	case "gauge":
		localGauge, err := s.IStorage.GetGauge(ctx, req.GetId())
		if err != nil {
			return nil, status.Error(codes.NotFound, "Название метрики не найдено")
		}
		result.Value = localGauge
	case "counter":
		localCounter, err := s.IStorage.GetCounter(ctx, req.GetId())
		if err != nil {
			return nil, status.Error(codes.NotFound, "Название метрики не найдено")
		}
		result.Delta = localCounter
	default:
		return nil, status.Error(codes.NotFound, "Неверный тип метрики! Допустимые значения: gauge, counter")
	}
	return &result, nil
}

func (s *MetricsServer) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	if err := checkIP(ctx); err != nil {
		return nil, err
	}
	data, err := s.IStorage.GetData(ctx)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "failed to get ping")
	}
	return &pb.GetDataResponse{Data: data}, nil
}

func (s *MetricsServer) GetPing(ctx context.Context, _ *pb.Ping) (*pb.Ping, error) {
	if err := checkIP(ctx); err != nil {
		return nil, err
	}
	err := s.IStorage.Ping()
	if err != nil {
		return nil, status.Error(codes.Unavailable, "failed to get ping")
	}
	return &pb.Ping{}, nil
}

func (s *MetricsServer) UpdateMetrics(ctx context.Context, req *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var response pb.UpdateMetricsResponse
	if err := checkIP(ctx); err != nil {
		return nil, err
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

func checkIP(ctx context.Context) error {
	if config.TrustedSubnet != "" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Error(codes.DataLoss, "failed to get metadata")
		}
		xrip := md["x-real-ip"]
		if len(xrip) == 0 {
			return status.Error(codes.InvalidArgument, "missing 'X-Real-IP' header")
		} else {
			ip := net.ParseIP(xrip[0])
			ones, _ := ip.DefaultMask().Size()
			_, i, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", ip.To4(), ones))
			mask := i.String()
			if mask != config.TrustedSubnet {
				return status.Error(codes.PermissionDenied, "deny")
			}
		}
	}
	return nil
}
