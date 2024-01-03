package worker

import (
	"context"
	"fmt"
	"sync"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
	"github.com/HoskeOwl/portscan/internal/task"
)

func MakeWorkQueue(tasks []task.ScanTask) <-chan task.ScanTask {
	ch := make(chan task.ScanTask)
	go func() {
		defer close(ch)
		for _, t := range tasks {
			ch <- t
		}
	}()
	return ch
}

func worker(ctx context.Context, wg *sync.WaitGroup, tq <-chan task.ScanTask, rq chan<- task.ScanTask) {
	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		panic(fmt.Errorf("internal context error"))
	}
	defer wg.Done()
	for t := range tq {
	RETRY:
		for t.RunCnt() < storage.MaxRunCnt {
			t.Do(ctx)
			if t.Error() == nil {
				break RETRY
			}
		}
		rq <- t
	}
}

func MakeWorkers(ctx context.Context, n int, tq <-chan task.ScanTask, rq chan<- task.ScanTask) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go worker(ctx, &wg, tq, rq)
	}
	return &wg
}
