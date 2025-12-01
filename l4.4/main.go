package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"
	"runtime"
	"runtime/debug"
	"time"
)

var memAllocGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "memory_Alloc",
		Help: "Heap Alloc",
	})

var memMallocGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "memory_Malloc",
		Help: "The number of live object",
	})
var numGCGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "gc_Num",
		Help: "number of completed GC cycles",
	})

var lastGCGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "gc_Last",
		Help: "last garbage collection finished",
	})

var percentGSGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "gc_percent",
		Help: "Current GC percent",
	},
)

func main() {
	var m runtime.MemStats
	debug.SetGCPercent(70)

	prometheus.MustRegister(memAllocGauge, memMallocGauge, numGCGauge, lastGCGauge, percentGSGauge)

	go func() {
		for {
			runtime.ReadMemStats(&m)
			memAllocGauge.Set(float64(m.Alloc))
			memMallocGauge.Set(float64(m.Mallocs))
			numGCGauge.Set(float64(m.NumGC))
			lastGCGauge.Set(float64(m.LastGC))
			percentGSGauge.Set(float64(-1))
			time.Sleep(3 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/debug/pprof/", pprof.Index)
	http.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	http.HandleFunc("/debug/pprof/profile", pprof.Profile)
	http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	http.HandleFunc("/debug/pprof/trace", pprof.Trace)

	fmt.Println("Listening to port", ":8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
