package main

import (
	"fmt"
	"runtime"
)

var HOST = "http://localhost:8080"

func main() {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)
	fmt.Printf("Total allocated memory (in bytes): %d\n", memStats.Alloc)

	//http.Post()
	//r.Post("/update/{type}/{name}/{value}", handlers.MetHandler)
	//err := http.ListenAndServe(":8080", r)
	//if err != nil {
	//	return
	//}
}
