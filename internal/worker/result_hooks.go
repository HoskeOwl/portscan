package worker

import (
	"context"
	"fmt"
	"sort"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
	"github.com/HoskeOwl/portscan/internal/task"
)

type ResultHook interface {
	Success(context.Context, task.ScanTask)
	Retry(context.Context, task.ScanTask)
	Fail(context.Context, task.ScanTask)
	Print(context.Context)
}

type SortedResultHook struct {
	success []task.ScanTask
	fail    []task.ScanTask
}

func (srh *SortedResultHook) init() {
	srh.success = make([]task.ScanTask, 0)
	srh.fail = make([]task.ScanTask, 0)
}

func (srh *SortedResultHook) Success(ctx context.Context, t task.ScanTask) {
	srh.success = append(srh.success, t)
}

func (srh *SortedResultHook) Fail(ctx context.Context, t task.ScanTask) {
	srh.fail = append(srh.fail, t)
}

func (srh *SortedResultHook) Retry(ctx context.Context, t task.ScanTask) {}

func (srh *SortedResultHook) Print(ctx context.Context) {
	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		panic(fmt.Errorf("internal context error"))
	}
	sort.SliceStable(srh.success, func(i, j int) bool { return srh.success[i].GetPort() < srh.success[j].GetPort() })
	if len(srh.success) > 0 {
		fmt.Println()
	}
	for _, t := range srh.success {
		fmt.Printf("    %v\n", t)
	}
	if storage.Verbose {
		sort.SliceStable(srh.fail, func(i, j int) bool { return srh.fail[i].GetPort() < srh.fail[j].GetPort() })
		if len(srh.fail) > 0 {
			fmt.Println()
		}
		for _, t := range srh.fail {
			fmt.Printf("    %v\n", t)
		}
	} else {
		if len(srh.success) == 0 {
			fmt.Println("\nNo opened ports")
		}
	}
}

func MakeSortedResultHook() *SortedResultHook {
	srh := &SortedResultHook{}
	srh.init()
	return srh
}

type RealtimeResultHook struct {
	success bool
}

func (rrh *RealtimeResultHook) Success(_ context.Context, t task.ScanTask) {
	fmt.Printf("    %v\n", t)
	rrh.success = true
}

func (rrh *RealtimeResultHook) Fail(ctx context.Context, t task.ScanTask) {
	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		panic(fmt.Errorf("internal context error"))
	}
	if storage.Verbose {
		fmt.Printf("    %v\n", t)
	}
}

func (rrh *RealtimeResultHook) Retry(ctx context.Context, t task.ScanTask) {
	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		panic(fmt.Errorf("internal context error"))
	}
	if storage.Verbose {
		fmt.Printf("    %v, will retry\n", t)
	}
}

func (rrh *RealtimeResultHook) Print(ctx context.Context) {
	if !rrh.success {
		fmt.Println("\nNo opened ports")
	}
}

func MakeRealtimeResultHook() *RealtimeResultHook {
	return &RealtimeResultHook{}
}

type JsonResultHook struct {
	t []task.ScanTask
}

func (jrh *JsonResultHook) init() {
	jrh.t = make([]task.ScanTask, 0)
}

func (jrh *JsonResultHook) Success(ctx context.Context, t task.ScanTask) {
	jrh.t = append(jrh.t, t)
}

func (jrh *JsonResultHook) Retry(ctx context.Context, t task.ScanTask) {}

func (jrh *JsonResultHook) Fail(ctx context.Context, t task.ScanTask) {
	jrh.t = append(jrh.t, t)
}

func (jrh *JsonResultHook) Print(ctx context.Context) {
	fmt.Println("{")
	l := len(jrh.t) - 1
	var end string
	for i, t := range jrh.t {
		if i != l {
			end = ","
		} else {
			end = ""
		}
		fmt.Printf("  \"%v\":\"%v\"%v\n", t.ConnStr(), t.Status(), end)
	}
	fmt.Println("}")
}

func MakeJsonResultHook() *JsonResultHook {
	jrh := &JsonResultHook{}
	jrh.init()
	return jrh
}
