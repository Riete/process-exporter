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
	search   []string
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
			panic(err)
		}
		// do not track systemd(1) and kernel thread
		if ppid == 0 || ppid == 2 {
			continue
		}
		cmdline, err = i.Cmdline()
		if err != nil {
			continue
		}
		found := false
		for _, c := range s.search {
			if strings.Contains(cmdline, c) {
				found = true
				break
			}
		}
		if found {
			s.p = append(s.p, i)
			s.cmdlines.Store(i.Pid, cmdline)
		}
	}
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

func New(search []string) *Storage {
	return &Storage{search: search}
}
