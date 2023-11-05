package ctxstorage

import (
	"context"
	"fmt"
	"time"
)

type MyKey string

const StorageKey MyKey = "key"

type CtxStorage struct {
	ConnDuration time.Duration
	Retries      int
	Verbose      bool
}

func FromContext(ctx context.Context) (*CtxStorage, error) {
	a := ctx.Value(StorageKey)
	if s, ok := a.(CtxStorage); !ok {
		return nil, fmt.Errorf("wrong storage type")
	} else {
		return &s, nil
	}
}
