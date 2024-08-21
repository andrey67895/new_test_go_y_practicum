package helpers

import (
	"math/rand"
	"runtime"
)

// GetMemByStats получение метрики по name
func GetMemByStats(name string) float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	switch name {
	case "Alloc":
		return float64(memStats.Alloc)
	case "BuckHashSys":
		return float64(memStats.BuckHashSys)
	case "Frees":
		return float64(memStats.Frees)
	case "GCCPUFraction":
		return memStats.GCCPUFraction
	case "GCSys":
		return float64(memStats.GCSys)
	case "HeapAlloc":
		return float64(memStats.HeapAlloc)
	case "HeapIdle":
		return float64(memStats.HeapIdle)
	case "HeapInuse":
		return float64(memStats.HeapInuse)
	case "HeapObjects":
		return float64(memStats.HeapObjects)
	case "HeapReleased":
		return float64(memStats.HeapReleased)
	case "HeapSys":
		return float64(memStats.HeapSys)
	case "LastGC":
		return float64(memStats.LastGC)
	case "Lookups":
		return float64(memStats.Lookups)
	case "MCacheInuse":
		return float64(memStats.MCacheInuse)
	case "MCacheSys":
		return float64(memStats.MCacheSys)
	case "MSpanInuse":
		return float64(memStats.MSpanInuse)
	case "MSpanSys":
		return float64(memStats.MSpanSys)
	case "Mallocs":
		return float64(memStats.Mallocs)
	case "NumForcedGC":
		return float64(memStats.NumForcedGC)
	case "NumGC":
		return float64(memStats.NumGC)
	case "OtherSys":
		return float64(memStats.OtherSys)
	case "PauseTotalNs":
		return float64(memStats.PauseTotalNs)
	case "StackInuse":
		return float64(memStats.StackInuse)
	case "Sys":
		return float64(memStats.Sys)
	case "TotalAlloc":
		return float64(memStats.TotalAlloc)
	case "RandomValue":
		return rand.Float64()
	default:
		return 0
	}
}
