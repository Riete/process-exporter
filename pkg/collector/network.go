package collector

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/riete/process-exporter/pkg/storage"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	networkConnection = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, networkSubsystem, "connections"),
		"Process network connections",
		[]string{"pid", "cmdline"},
		nil,
	)
)

type NetworkCollector struct {
	s *storage.Storage
}

func (n NetworkCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- networkConnection
}

func (n NetworkCollector) Collect(metrics chan<- prometheus.Metric) {
	ch := make(chan *process.Process)
	go n.s.Fetch(ch)
	for p := range ch {
		cmdline := n.s.ProcessCmdline(p.Pid)
		conn, err := p.Connections()
		if err != nil {
			log.Printf("Get [%s] Process Network Connections Error: %v\n", cmdline, err)
			continue
		}
		pid := strconv.Itoa(int(p.Pid))
		metrics <- prometheus.MustNewConstMetric(networkConnection, prometheus.GaugeValue, float64(len(conn)), pid, cmdline)
	}
}

func NewNetworkCollector(s *storage.Storage) prometheus.Collector {
	return &NetworkCollector{s: s}
}
