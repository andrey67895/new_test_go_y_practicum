package interceptors

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
)

var log = logger.Log()

func RealIPInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	if err := checkIP(ctx); err != nil {
		return nil, err
	}
	m, err := handler(ctx, req)
	if err != nil {
		log.Error("RPC failed with error: %v", err)
	}
	return m, err
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
