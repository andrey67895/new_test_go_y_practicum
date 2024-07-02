package storage

import "context"

type IStorageData interface {
	RetrySaveGauge(ctx context.Context, id string, delta float64) error
	RetrySaveCounter(ctx context.Context, id string, value int64) error
	GetCounter(ctx context.Context, id string) (int64, error)
	GetGauge(ctx context.Context, id string) (float64, error)
	GetData(ctx context.Context) (string, error)
	Ping() error
}
