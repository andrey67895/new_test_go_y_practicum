package storage

type IStorageData interface {
	RetrySaveGauge(id string, delta float64) error
	RetrySaveCounter(id string, value int64) error
	GetCounter(id string) (int64, error)
	GetGauge(id string) (float64, error)
	GetData() (string, error)
	Ping() error
}
