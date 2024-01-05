package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riete/process-exporter/pkg/collector"
	"github.com/riete/process-exporter/pkg/storage"
)

func main() {
	listenPort := flag.String("listen-port", "10921", "listen port")
	pnc := flag.String("process-name-contains", "", "fuzzy match process, ',' separated")
	flag.Parse()

	s := storage.New(strings.Split(*pnc, ","))
	go func() {
		for {
			s.SyncOrDie()
			time.Sleep(5 * time.Minute)
		}
	}()
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collector.NewCpuCollector(s),
		collector.NewMemoryCollector(s),
		collector.NewNetworkCollector(s),
		collector.NewIOCollector(s),
		collector.NewCommonCollector(s),
	)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	log.Printf("server listen at :%s, open http://127.0.0.1:%s/metrics to view process metrics", *listenPort, *listenPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *listenPort), nil))
}
