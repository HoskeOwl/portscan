package helper

import (
	"github.com/HoskeOwl/portscan/internal/task"
)

func CreateScanTasks(ip string, prs []PortRange) []task.ScanTask {
	t := make([]task.ScanTask, 0)
	for _, pr := range prs {
		for _, p := range pr.Ports() {
			st := task.MakeTcpScanTask(ip, p)
			t = append(t, st)
		}
	}
	return t
}
