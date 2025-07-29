package db

import (
	"dubhe/db/util/log"
	"fmt"
	"sync"
)

func (r *Repo[T]) checkErr(e error) {
	r.log.Info(fmt.Sprintf("repo [%s] error: %s", r.key, e.Error()), log.KV("error", e))

	// 先调用全局注册的错误处理函数

	for _, gh := range globalHandle {
		h := newErrHandler[T](r, e)
		gh(h)
		if h.isPanic {
			panic(h.err)
		}
		if h.isContinue {
			continue
		}
		if h.isCancel {
			return
		}
	}

	if r.onErrFunc != nil {
		for _, f := range r.onErrFunc {
			h := newErrHandler[T](r, e)
			f(h)
			if h.isPanic {
				panic(h.err)
			}
			if h.isContinue {
				continue
			}
			if h.isCancel {
				return
			}
		}
	}
}

var (
	mu           sync.Mutex
	globalHandle []GlobalErrHandle
)

// RegisterGlobalHandle 注册新的全局错误处理器
func RegisterGlobalHandle(handle GlobalErrHandle) {
	mu.Lock()
	defer mu.Unlock()
	globalHandle = append(globalHandle, handle)
}

// ClearGlobalHandles 清空所有全局错误处理器
func ClearGlobalHandles() {
	mu.Lock()
	defer mu.Unlock()
	globalHandle = nil
}

// ListGlobalHandles 返回当前所有已注册的错误处理器（返回副本，避免外部修改）
func ListGlobalHandles() []GlobalErrHandle {
	mu.Lock()
	defer mu.Unlock()
	handlesCopy := make([]GlobalErrHandle, len(globalHandle))
	copy(handlesCopy, globalHandle)
	return handlesCopy
}
