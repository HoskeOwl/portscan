package worker

import (
	"context"
	"fmt"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
	"github.com/HoskeOwl/portscan/internal/task"
)

type Tasker interface {
	String() string
	RunCnt() int
	Do(context.Context)
	Error() error
	GetPort() int
	Status() string
	ConnStr() string
}

type Hooker func(context.Context, Tasker)

type Pool struct {
	queueChan     chan Tasker
	resultQueue   chan Tasker
	waitResultCnt int

	WorkersCnt  int
	Started     bool
	Finished    bool
	SuccessHook Hooker
	RetrieHook  Hooker
	FailHook    Hooker
}

func (p *Pool) AddTask(t Tasker) {
	p.queueChan <- t
}

func (p *Pool) init(cnt int) {
	p.queueChan = make(chan Tasker, cnt)
	p.resultQueue = make(chan Tasker, cnt)
	p.waitResultCnt += cnt
}

func (p *Pool) execution(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case t, ok := <-p.queueChan:
			if !ok {
				return
			}
			t.Do(ctx)
			p.resultQueue <- t
		}
	}
}

func (p *Pool) resultChecker(ctx context.Context) int {
	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		panic(fmt.Errorf("internal context error"))
	}
	success := 0
	for {
		select {
		case task := <-p.resultQueue:
			if task.Error() != nil && task.RunCnt() < storage.Retries {
				p.queueChan <- task
				if p.RetrieHook != nil {
					p.RetrieHook(ctx, task)
				}
			} else {
				p.waitResultCnt -= 1
				if task.Error() == nil {
					success += 1
					if p.SuccessHook != nil {
						p.SuccessHook(ctx, task)
					}
				} else {
					if p.FailHook != nil {
						p.FailHook(ctx, task)
					}
				}
			}
			if p.waitResultCnt == 0 {
				return success
			}
		case <-ctx.Done():
			p.Finished = true
			return success
		}
	}
}

func (p *Pool) Execute(ctx context.Context) (successCnt int) {
	if p.Started || p.Finished {
		return
	}
	wctx, cancel := context.WithCancel(ctx)
	for i := 0; i < p.WorkersCnt; i++ {
		go p.execution(wctx)
	}
	successCnt = p.resultChecker(wctx)
	cancel()
	return
}

func (p *Pool) WithSuccessHook(h Hooker) *Pool {
	p.SuccessHook = h
	return p
}

func (p *Pool) WithFailHook(h Hooker) *Pool {
	p.FailHook = h
	return p
}

func MakePool(ctx context.Context, t []task.ScanTask, workersCnt int) *Pool {
	pool := Pool{WorkersCnt: workersCnt}
	pool.init(len(t))
	//  Do not use "range" because range use same address for the variable (make a copy)
	for i := 0; i < len(t); i++ {
		pool.AddTask(&(t[i]))

	}
	return &pool
}
