package worker

import (
	"context"
	"fmt"
	"sort"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
)

type ResultProcessor interface {
	Success(context.Context, Tasker)
	Fail(context.Context, Tasker)
	Print(context.Context)
}

type SortedResultProcessor struct {
	success []Tasker
	fail    []Tasker
}

func (srp *SortedResultProcessor) init() {
	srp.success = make([]Tasker, 0)
	srp.fail = make([]Tasker, 0)
}

func (srp *SortedResultProcessor) Success(ctx context.Context, t Tasker) {
	srp.success = append(srp.success, t)
}

func (srp *SortedResultProcessor) Fail(ctx context.Context, t Tasker) {
	srp.fail = append(srp.fail, t)
}

func (srp *SortedResultProcessor) Print(ctx context.Context) {
	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		panic(fmt.Errorf("internal context error"))
	}
	sort.SliceStable(srp.success, func(i, j int) bool { return srp.success[i].GetPort() < srp.success[j].GetPort() })
	if len(srp.success) > 0 {
		fmt.Println()
	}
	for _, t := range srp.success {
		fmt.Printf("    %v\n", t)
	}
	if storage.Verbose {
		sort.SliceStable(srp.fail, func(i, j int) bool { return srp.fail[i].GetPort() < srp.fail[j].GetPort() })
		if len(srp.fail) > 0 {
			fmt.Println()
		}
		for _, t := range srp.fail {
			fmt.Printf("    %v\n", t)
		}
	} else {
		if len(srp.success) == 0 {
			fmt.Println("\nNo opened ports")
		}
	}
}

func MakeSortedResultProcessor() *SortedResultProcessor {
	srp := &SortedResultProcessor{}
	srp.init()
	return srp
}

type RealtimeResultProcessor struct {
	success bool
}

func (rrp *RealtimeResultProcessor) Success(_ context.Context, t Tasker) {
	fmt.Printf("    %v\n", t)
	rrp.success = true
}

func (rrp *RealtimeResultProcessor) Fail(ctx context.Context, t Tasker) {
	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		panic(fmt.Errorf("internal context error"))
	}
	if storage.Verbose {
		fmt.Printf("    %v\n", t)
	}
}

func (rrp *RealtimeResultProcessor) Print(ctx context.Context) {
	if !rrp.success {
		fmt.Println("\nNo opened ports")
	}
}

func MakeRealtimeResultProcessor() *RealtimeResultProcessor {
	return &RealtimeResultProcessor{}
}

type JsonResultProcessor struct {
	t []Tasker
}

func (jrp *JsonResultProcessor) init() {
	jrp.t = make([]Tasker, 0)
}

func (jrp *JsonResultProcessor) Success(ctx context.Context, t Tasker) {
	jrp.t = append(jrp.t, t)
}

func (jrp *JsonResultProcessor) Fail(ctx context.Context, t Tasker) {
	jrp.t = append(jrp.t, t)
}

func (jrp *JsonResultProcessor) Print(ctx context.Context) {
	fmt.Println("{")
	l := len(jrp.t) - 1
	var end string
	for i, t := range jrp.t {
		if i != l {
			end = ","
		} else {
			end = ""
		}
		fmt.Printf("  \"%v\":\"%v\"%v\n", t.ConnStr(), t.Status(), end)
	}
	fmt.Println("}")
}

func MakeJsonResultProcessor() *JsonResultProcessor {
	jrp := &JsonResultProcessor{}
	jrp.init()
	return jrp
}
