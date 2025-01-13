package collector

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/riete/process-exporter/pkg/storage"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	memoryRss = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, memorySubsystem, "rss_bytes"),
		"Process rss memory bytes",
		commonLabels,
		nil,
	)
)

type MemoryCollector struct {
	s *storage.Storage
}

func (m MemoryCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- memoryRss
}

func (m MemoryCollector) Collect(metrics chan<- prometheus.Metric) {
	ch := make(chan *process.Process)
	go m.s.Fetch(ch)
	for p := range ch {
		cmdline := m.s.ProcessCmdline(p.Pid)
		mem, err := p.MemoryInfo()
		if err != nil {
			log.Printf("Get [%s] Process Memory Error: %v\n", cmdline, err)
			continue
		}
		pid := strconv.Itoa(int(p.Pid))
		metrics <- prometheus.MustNewConstMetric(memoryRss, prometheus.GaugeValue, float64(mem.RSS), pid, cmdline)
	}
}

func NewMemoryCollector(s *storage.Storage) prometheus.Collector {
	return &MemoryCollector{s: s}
}
