package collector

import (
	"log"
	"strconv"

	"github.com/prometheus/procfs"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/riete/process-exporter/pkg/storage"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	networkConnection = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, networkSubsystem, "connections"),
		"Process network connections",
		commonLabels,
		nil,
	)
	networkReceiveBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, networkSubsystem, "receive_bytes_total"),
		"Supervisor Process Network Receive Bytes",
		commonLabels,
		nil,
	)
	networkTransmitBytes = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, networkSubsystem, "transmit_bytes_total"),
		"Supervisor Process Network Transmit Bytes",
		commonLabels,
		nil,
	)
)

type NetworkCollector struct {
	s *storage.Storage
}

func (n NetworkCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- networkConnection
	descs <- networkReceiveBytes
	descs <- networkTransmitBytes
}

func (n NetworkCollector) Collect(metrics chan<- prometheus.Metric) {
	ch := make(chan *process.Process)
	go n.s.Fetch(ch)
	for p := range ch {
		cmdline := n.s.ProcessCmdline(p.Pid)
		conn, err := p.Connections()
		pid := strconv.Itoa(int(p.Pid))
		if err != nil {
			log.Printf("Get [%s] Process Network Connections Error: %v\n", cmdline, err)
		} else {
			metrics <- prometheus.MustNewConstMetric(networkConnection, prometheus.GaugeValue, float64(len(conn)), pid, cmdline)
		}
		pn, err := procfs.NewProc(int(p.Pid))
		if err != nil {
			log.Printf("Get [%s] Process Network Traffic Error: %v\n", cmdline, err)
			continue
		}
		netstat, err := pn.Netstat()
		if err != nil {
			log.Printf("Get [%s] Process Network Traffic Error: %v\n", cmdline, err)
			continue
		}
		if netstat.IpExt.InOctets != nil {
			metrics <- prometheus.MustNewConstMetric(networkReceiveBytes, prometheus.CounterValue, *netstat.IpExt.InOctets, pid, cmdline)
		}
		if netstat.IpExt.OutOctets != nil {
			metrics <- prometheus.MustNewConstMetric(networkTransmitBytes, prometheus.CounterValue, *netstat.IpExt.OutOctets, pid, cmdline)
		}
	}
}

func NewNetworkCollector(s *storage.Storage) prometheus.Collector {
	return &NetworkCollector{s: s}
}
