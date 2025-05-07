package storage

import (
	"strings"
	"sync"

	"github.com/shirou/gopsutil/v3/process"
)

type Storage struct {
	p        []*process.Process
	cmdlines sync.Map
	l        sync.RWMutex
	include  []string
	exclude  []string
	total    int
}

func (s *Storage) SyncOrDie() {
	s.l.Lock()
	defer s.l.Unlock()

	var err error
	var p []*process.Process
	p, err = process.Processes()
	if err != nil {
		panic(err)
	}
	s.total = len(p)

	var ppid int32
	var cmdline string
	s.p = nil

	for _, i := range p {
		ppid, err = i.Ppid()
		if err != nil {
			continue
		}
		// do not track systemd(1) and kernel thread
		if ppid == 0 || ppid == 2 {
			continue
		}
		cmdline, err = i.Cmdline()
		if err != nil {
			continue
		}
		if s.shouldFetchMetric(cmdline) {
			s.p = append(s.p, i)
			s.cmdlines.Store(i.Pid, cmdline)
		}
	}
}

// shouldFetchMetric exclude first then include
func (s *Storage) shouldFetchMetric(cmdline string) bool {
	if len(s.exclude) > 0 {
		for _, i := range s.exclude {
			if i != "" && strings.Contains(cmdline, i) {
				return false
			}
		}
	}

	if len(s.include) > 0 {
		for _, i := range s.include {
			if i == "" || strings.Contains(cmdline, i) {
				return true
			}
		}
		return false
	}
	return true
}

func (s *Storage) ProcessCmdline(pid int32) string {
	cmdline, _ := s.cmdlines.Load(pid)
	return cmdline.(string)
}

func (s *Storage) Fetch(ch chan *process.Process) {
	s.l.RLock()
	defer s.l.RUnlock()
	defer close(ch)
	for _, i := range s.p {
		ch <- i
	}
}

func (s *Storage) Total() int {
	return s.total
}

func New(include, exclude []string) *Storage {
	return &Storage{include: include, exclude: exclude}
}
