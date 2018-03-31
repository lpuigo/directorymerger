package copymanager

import (
	"io"
	"sync"
	"time"
)

const (
	NotStarted = iota
	Started
	Done
)

type Monitor struct {
	mut        sync.Mutex
	c          <-chan uint64

	startTime  time.Time
	updateTime time.Time

	size       uint64
	updateSize uint64

	status     int
}

func NewMonitor(c <-chan uint64) *Monitor {
	res := &Monitor{c: c}
	go res.update()
	return res
}

func (m *Monitor) update() {
	ongoing := true
	for ongoing {
		i, working := <-m.c
		m.mut.Lock()
		if m.startTime.IsZero() {
			m.startTime = time.Now()
			m.updateTime = m.startTime
			m.status = Started
		}
		if !working {
			m.c = nil
			ongoing = false
			m.status = Done
		} else {
			m.size += i
		}
		m.mut.Unlock()
	}
}

func (m *Monitor) GetInfo() (size uint64, sinceStart, absThruput, instThruput float64, inprogress bool) {
	m.mut.Lock()
	defer m.mut.Unlock()

	switch m.status {
	case NotStarted:
		return 0, 0, 0,0, true
	case Started, Done:
		now := time.Now()
		size = m.size
		sinceStart = now.Sub(m.startTime).Seconds()
		sinceUpdate := now.Sub(m.updateTime).Seconds()
		mb := float64(size/1024) / 1024
		mbu := float64((size-m.updateSize)/1024) / 1024

		absThruput = mb / sinceStart
		instThruput = mbu / sinceUpdate
		inprogress = (m.status == Started)

		m.updateTime = time.Now()
		m.updateSize = size
	}
	return
}

type MonitoredReader struct {
	r io.Reader
	c chan uint64
}

func (mr *MonitoredReader) Read(p []byte) (n int, err error) {
	n, err = mr.r.Read(p)
	mr.c <- uint64(n)
	if n == 0 && err == io.EOF {
		close(mr.c)
	}
	return
}

func NewMonitoredReader(r io.Reader) (*MonitoredReader, *Monitor) {
	c := make(chan uint64, 10)
	mr := &MonitoredReader{r: r, c: c}
	m := NewMonitor(c)

	return mr, m
}
