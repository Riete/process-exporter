package collector

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/riete/process-exporter/pkg/storage"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	readCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ioSubsystem, "read_count"),
		"Process io read count",
		commonLabels,
		nil,
	)
	writeCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ioSubsystem, "write_count"),
		"Process io write count",
		commonLabels,
		nil,
	)
	readBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ioSubsystem, "read_bytes"),
		"Process io read bytes",
		commonLabels,
		nil,
	)
	writeBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ioSubsystem, "write_bytes"),
		"Process io write bytes",
		commonLabels,
		nil,
	)
)

type IOCollector struct {
	s *storage.Storage
}

func (i IOCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- readCount
	descs <- writeCount
	descs <- readBytes
	descs <- writeBytes
}

func (i IOCollector) Collect(metrics chan<- prometheus.Metric) {
	ch := make(chan *process.Process)
	go i.s.Fetch(ch)
	for p := range ch {
		cmdline := i.s.ProcessCmdline(p.Pid)
		pid := strconv.Itoa(int(p.Pid))
		iostat, err := p.IOCounters()
		if err != nil {
			log.Printf("Get [%s] Process IO Stats Error: %v\n", cmdline, err)
			continue
		}
		metrics <- prometheus.MustNewConstMetric(readCount, prometheus.CounterValue, float64(iostat.ReadCount), pid, cmdline)
		metrics <- prometheus.MustNewConstMetric(writeCount, prometheus.CounterValue, float64(iostat.WriteCount), pid, cmdline)
		metrics <- prometheus.MustNewConstMetric(readBytes, prometheus.CounterValue, float64(iostat.ReadBytes), pid, cmdline)
		metrics <- prometheus.MustNewConstMetric(writeBytes, prometheus.CounterValue, float64(iostat.WriteBytes), pid, cmdline)
	}
}

func NewIOCollector(s *storage.Storage) prometheus.Collector {
	return &IOCollector{s: s}
}
