package worker

import (
	"context"
	"fmt"
	"sync"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
	"github.com/HoskeOwl/portscan/internal/task"
)

func resultProcessor(ctx context.Context, wg *sync.WaitGroup, rq <-chan task.ScanTask, rh ResultHook) {
	defer wg.Done()
	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		panic(fmt.Errorf("internal context error"))
	}
	for task := range rq {
		if task.Error() != nil && task.RunCnt() <= storage.MaxRunCnt {
			if rh != nil {
				rh.Retry(ctx, task)
			}
		} else {
			if task.Error() == nil {
				if rh != nil {
					rh.Success(ctx, task)
				}
			} else {
				if rh != nil {
					rh.Fail(ctx, task)
				}
			}
		}

	}
}

func RunResultProcessor(ctx context.Context, rq <-chan task.ScanTask, rh ResultHook) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go resultProcessor(ctx, &wg, rq, rh)
	return &wg
}
