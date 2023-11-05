package task_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
	"github.com/HoskeOwl/portscan/internal/task"
	"github.com/stretchr/testify/assert"
)

type TestCheckerOk struct{}

func (t TestCheckerOk) Check(_ string, _ time.Duration) error {
	return nil
}

type TestCheckerErr struct{}

func (t TestCheckerErr) Check(_ string, _ time.Duration) error {
	return fmt.Errorf("test error")
}

func TestTaskOk(t *testing.T) {
	ip := "1.1.1.1"
	port := 33
	storage := ctxstorage.CtxStorage{}
	ctx := context.WithValue(context.Background(), ctxstorage.StorageKey, storage)
	task := task.ScanTask{Ip: ip, Port: port, Checker: &TestCheckerOk{}}
	assert.Equal(t, task.Status(), "waiting")
	task.Do(ctx)
	assert.Equal(t, task.Status(), "success")
	assert.Nil(t, task.Error())
}

func TestTaskFail(t *testing.T) {
	ip := "1.1.1.1"
	port := 33
	storage := ctxstorage.CtxStorage{}
	ctx := context.WithValue(context.Background(), ctxstorage.StorageKey, storage)
	task := task.ScanTask{Ip: ip, Port: port, Checker: &TestCheckerErr{}}
	assert.Equal(t, task.Status(), "waiting")
	task.Do(ctx)
	assert.Equal(t, task.Status(), "fail")
	assert.NotNil(t, task.Error())
}

func TestTaskConnString(t *testing.T) {
	ip := "5.5.5.5"
	port := 555
	task := task.ScanTask{Ip: ip, Port: port, Checker: &TestCheckerOk{}}
	exp := fmt.Sprintf("%v:%v", ip, port)
	assert.Equal(t, task.ConnStr(), exp)
}
