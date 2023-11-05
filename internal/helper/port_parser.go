package helper

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	OldPortRangeSeparator string = ":"
	PortRangeSeparator    string = "-"
	RangesSeparator       string = ","
	// minPort                      = 1
	// maxPort                      = 65535
)

type PortRange struct {
	Start int
	End   int //not include
}

func (pr PortRange) String() string {
	if pr.End-pr.Start == 1 {
		return fmt.Sprintf("%v", pr.Start)
	}
	return fmt.Sprintf("%v%v%v", pr.Start, PortRangeSeparator, pr.End-1)
}

func (pr *PortRange) Ports() []int {
	ports := make([]int, 0, pr.End-pr.Start)
	for i := pr.Start; i < pr.End; i++ {
		ports = append(ports, i)
	}
	return ports
}

func ParseRange(portsRange string) (*PortRange, error) {
	var sep string
	switch {
	case strings.Contains(portsRange, PortRangeSeparator):
		sep = PortRangeSeparator
	case strings.Contains(portsRange, OldPortRangeSeparator):
		sep = OldPortRangeSeparator
	default:
		// single port
		start, err := strconv.Atoi(portsRange)
		if err != nil {
			return nil, err
		}
		return &PortRange{Start: start, End: (start + 1)}, nil

	}

	p := strings.Split(portsRange, sep)
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
	if start > end {
		start, end = end, start
	}
	return &PortRange{Start: start, End: end + 1}, nil

}

func ParsePortRanges(ranges string) ([]PortRange, error) {
	sr := strings.Split(ranges, RangesSeparator)
	pRanges := make([]PortRange, 0, len(sr))
	for _, r := range sr {
		pr, err := ParseRange(r)
		if err != nil {
			return nil, err
		}
		pRanges = append(pRanges, *pr)
	}
	return pRanges, nil
}
