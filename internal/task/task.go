package task

import (
	"context"
	"fmt"
	"net"
	"time"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
)

const (
	tcpNetwork = "tcp"
	fail       = "fail"
	success    = "success"
	waiting    = "waiting"
	undefined  = "undefined"
)

type ConnChecker interface {
	Check(string, time.Duration) error
}

type checker struct {
	Net string
}

func (t *checker) Check(dst string, timeout time.Duration) error {
	conn, err := net.DialTimeout(t.Net, dst, timeout)
	if err == nil {
		conn.Close()
	}
	return err
}

func makeTcpChecker() *checker {
	return &checker{Net: tcpNetwork}
}

type ScanTask struct {
	Ip      string
	Port    int
	Checker ConnChecker

	runCnt        int
	lastExecution time.Time
	lastDuratuion time.Duration
	err           error
}

func (st *ScanTask) Status() string {
	switch {
	case st.runCnt > 0 && st.err != nil:
		return fail
	case st.runCnt > 0 && st.err == nil:
		return success
	case st.runCnt == 0 && st.err == nil:
		return waiting
	default:
		return undefined
	}

}

func (st ScanTask) String() string {
	return fmt.Sprintf("%v: %v", st.Port, st.Status())
}

func (st *ScanTask) ConnStr() string {
	return fmt.Sprintf("%v:%v", st.Ip, st.Port)
}

func (st *ScanTask) RunCnt() int {
	return st.runCnt
}

func (st *ScanTask) Error() error {
	return st.err
}

func (st *ScanTask) GetPort() int {
	return st.Port
}

func (st *ScanTask) Do(ctx context.Context) {
	var err error
	st.runCnt += 1
	st.lastExecution = time.Now()
	start := time.Now()
	defer func() {
		st.lastDuratuion = time.Since(start)
		st.err = err
	}()

	storage, err := ctxstorage.FromContext(ctx)
	if err != nil {
		return
	}
	err = st.Checker.Check(st.ConnStr(), storage.ConnDuration)
}

func MakeTcpScanTask(ip string, port int) ScanTask {
	return ScanTask{Ip: ip, Port: port, Checker: makeTcpChecker()}
}
