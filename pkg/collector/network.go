package collector

import (
	"log"
	"strconv"

	"github.com/shirou/gopsutil/v3/net"

	"github.com/prometheus/procfs"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/riete/process-exporter/pkg/storage"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	networkTCPConnection = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, networkSubsystem, "tcp_connections"),
		"Process network tcp connections",
		commonLabels,
		nil,
	)
	networkTCPConnectionStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, networkSubsystem, "tcp_connections_status"),
		"Process network tcp connections status",
		append(commonLabels, "status"),
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
	descs <- networkTCPConnection
	descs <- networkTCPConnectionStatus
	descs <- networkReceiveBytes
	descs <- networkTransmitBytes
}

func (n NetworkCollector) Collect(metrics chan<- prometheus.Metric) {
	ch := make(chan *process.Process)
	go n.s.Fetch(ch)
	for p := range ch {
		cmdline := n.s.ProcessCmdline(p.Pid)
		pid := strconv.Itoa(int(p.Pid))
		tcpConns, err := net.ConnectionsPid("tcp", p.Pid)
		if err != nil {
			log.Printf("Get [%s] Process Network TCP Connections Error: %v\n", cmdline, err)
		} else {
			connStatus := make(map[string]float64)
			for _, c := range tcpConns {
				connStatus[c.Status] += 1
			}
			for s, c := range connStatus {
				metrics <- prometheus.MustNewConstMetric(networkTCPConnectionStatus, prometheus.GaugeValue, c, pid, cmdline, s)
			}
			metrics <- prometheus.MustNewConstMetric(networkTCPConnection, prometheus.GaugeValue, float64(len(tcpConns)), pid, cmdline)
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
