package collector

import (
	"log"
	"strconv"

	"github.com/shirou/gopsutil/v3/process"

	"github.com/riete/process-exporter/pkg/storage"

	"github.com/prometheus/client_golang/prometheus"
)

var cputime = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, cpuSubsystem, "seconds_total"),
	"Process cpu seconds in each mode.",
	append(commonLabels, "mode"),
	nil,
)

type CpuCollector struct {
	s *storage.Storage
}

func (c CpuCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- cputime
}

func (c CpuCollector) Collect(metrics chan<- prometheus.Metric) {
	ch := make(chan *process.Process)
	go c.s.Fetch(ch)
	for p := range ch {
		cmdline := c.s.ProcessCmdline(p.Pid)
		pid := strconv.Itoa(int(p.Pid))
		cpuTimes, err := p.Times()
		if err != nil {
			log.Printf("Get [%s] Process Cpu Time Error: %v\n", cmdline, err)
			continue
		}
		metrics <- prometheus.MustNewConstMetric(cputime, prometheus.CounterValue, cpuTimes.User, pid, cmdline, "user")
		metrics <- prometheus.MustNewConstMetric(cputime, prometheus.CounterValue, cpuTimes.System, pid, cmdline, "system")
		metrics <- prometheus.MustNewConstMetric(cputime, prometheus.CounterValue, cpuTimes.Iowait, pid, cmdline, "iowait")
	}
}

func NewCpuCollector(s *storage.Storage) prometheus.Collector {
	return &CpuCollector{s: s}
}
