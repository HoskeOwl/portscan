package scanworker

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	portRangeSeparator string = ":"
	rangesSeparator    string = ","
	defaultCapacity    int    = 10
)

type PortScanResult struct {
	Port      int
	Connected bool
}

func checkPort(ctx context.Context, ip string, timeoutMillis int, portTask <-chan int, resultQueue chan<- PortScanResult) {
	for {
		select {
		case <-ctx.Done():
			return
		case port, ok := <-portTask:
			if !ok {
				return
			}
			result := PortScanResult{Port: port, Connected: false}
			connection, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), time.Duration(timeoutMillis*int(time.Millisecond)))
			result.Connected = (err == nil)
			if err == nil {
				connection.Close()
			}
			resultQueue <- result
		}
	}
}

type PortRange struct {
	Start int
	End   int //not include
}

func (pr PortRange) String() string {
	if pr.End-pr.Start == 1 {
		return fmt.Sprintf("%v", pr.Start)
	}
	return fmt.Sprintf("%v:%v", pr.Start, pr.End)
}

func parseRange(portsRange string) (*PortRange, error) {
	if strings.Contains(portsRange, portRangeSeparator) {
		p := strings.Split(portsRange, portRangeSeparator)
		if len(p) != 2 {
			return nil, fmt.Errorf("wrong range syntax: '%v'", portsRange)
		}
		start, err := strconv.Atoi(p[0])
		if err != nil {
			return nil, err
		}
		end, err := strconv.Atoi(p[1])
		if err != nil {
			return nil, err
		}
		return &PortRange{Start: start, End: end}, nil
	}
	start, err := strconv.Atoi(portsRange)
	if err != nil {
		return nil, err
	}
	return &PortRange{Start: start, End: (start + 1)}, nil
}

func ParsePortRanges(ranges string) ([]PortRange, error) {
	pRanges := make([]PortRange, 0, defaultCapacity)
	for _, r := range strings.Split(ranges, rangesSeparator) {
		pr, err := parseRange(r)
		if err != nil {
			return nil, err
		}
		pRanges = append(pRanges, *pr)
	}
	return pRanges, nil
}

type PortScanQueue struct {
	ip               string
	routinesCount    int
	portRanges       []PortRange
	connectTimeoutMs int

	started          bool
	ctx              context.Context
	ports            chan int
	queueTaskResults chan PortScanResult
	success          []int
}

func (pq *PortScanQueue) startScan(totalLen int) error {
	if pq.started {
		return fmt.Errorf("already scanning")
	}
	pq.started = true
	cnt := totalLen
	if pq.routinesCount < totalLen {
		cnt = pq.routinesCount
	}
	for i := 0; i < cnt; i++ {
		go checkPort(pq.ctx, pq.ip, pq.connectTimeoutMs, pq.ports, pq.queueTaskResults)
	}
	for _, r := range pq.portRanges {
		for i := r.Start; i < r.End; i++ {
			pq.ports <- i
		}
	}
	return nil
}

func (pq *PortScanQueue) processResults(waitResults int) {
	for {
		select {
		case r := <-pq.queueTaskResults:
			if r.Connected {
				pq.success = append(pq.success, r.Port)
			}
			waitResults -= 1
			if waitResults <= 0 {
				pq.started = false
				return
			}
		case <-pq.ctx.Done():
			pq.started = false
			return
		}
	}
}

func StartPortScan(ctx context.Context, routinesCount int, ip string, portRanges []PortRange, connectTimeoutMs int) ([]int, error) {
	var totalLen int
	for _, r := range portRanges {
		totalLen += (r.End - r.Start)
	}
	if totalLen == 0 {
		return nil, fmt.Errorf("Missing ports")
	}
	queueCtx, queueCancel := context.WithCancel(ctx)
	ps := PortScanQueue{
		ip:               ip,
		routinesCount:    routinesCount,
		portRanges:       portRanges,
		connectTimeoutMs: connectTimeoutMs,
		ctx:              queueCtx,
		ports:            make(chan int, totalLen),
		queueTaskResults: make(chan PortScanResult, totalLen),
		success:          make([]int, 0, defaultCapacity),
	}
	err := ps.startScan(totalLen)
	if err != nil {
		queueCancel()
		return nil, err
	}
	ps.processResults(totalLen)
	queueCancel()
	return ps.success, nil
}
