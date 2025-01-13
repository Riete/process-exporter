package collector

import (
	"log"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/riete/process-exporter/pkg/storage"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	ctxSwitchVoluntary = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ctxSwitchSubsystem, "voluntary_count"),
		"Process ctx switch voluntary count",
		commonLabels,
		nil,
	)
	ctxSwitchInVoluntary = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, ctxSwitchSubsystem, "involuntary_count"),
		"Process ctx switch involuntary count",
		commonLabels,
		nil,
	)
	fd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, fdSubsystem, "count"),
		"Process fd count",
		commonLabels,
		nil,
	)
	thread = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, threadSubsystem, "count"),
		"Process thread count",
		commonLabels,
		nil,
	)
	processCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "total", "count"),
		"Process total count",
		nil,
		nil,
	)
)

type CommonCollector struct {
	s *storage.Storage
}

func (c CommonCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- ctxSwitchVoluntary
	descs <- ctxSwitchInVoluntary
	descs <- fd
	descs <- thread
	descs <- processCount
}

func (c CommonCollector) Collect(metrics chan<- prometheus.Metric) {

	ch := make(chan *process.Process)
	go c.s.Fetch(ch)
	for p := range ch {
		cmdline := c.s.ProcessCmdline(p.Pid)
		cw, err := p.NumCtxSwitches()
		if err != nil {
			log.Printf("Get [%s] Process CTX Switches Error: %v\n", cmdline, err)
			continue
		}
		fds, err := p.NumFDs()
		if err != nil {
			log.Printf("Get [%s] Process FDs Error: %v\n", cmdline, err)
			continue
		}
		threads, err := p.NumThreads()
		if err != nil {
			log.Printf("Get [%s] Process Threads Error: %v\n", cmdline, err)
			continue
		}
		pid := strconv.Itoa(int(p.Pid))
		metrics <- prometheus.MustNewConstMetric(ctxSwitchVoluntary, prometheus.CounterValue, float64(cw.Voluntary), pid, cmdline)
		metrics <- prometheus.MustNewConstMetric(ctxSwitchInVoluntary, prometheus.CounterValue, float64(cw.Involuntary), pid, cmdline)
		metrics <- prometheus.MustNewConstMetric(fd, prometheus.GaugeValue, float64(fds), pid, cmdline)
		metrics <- prometheus.MustNewConstMetric(thread, prometheus.GaugeValue, float64(threads), pid, cmdline)
	}
	metrics <- prometheus.MustNewConstMetric(processCount, prometheus.GaugeValue, float64(c.s.Total()))
}

func NewCommonCollector(s *storage.Storage) prometheus.Collector {
	return &CommonCollector{s: s}
}
