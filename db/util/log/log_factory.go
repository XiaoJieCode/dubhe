package log

import (
	"sync"
)

var (
	LogFactory = *NewLogFactory(NewDefaultLogger("default"))
)

// Factory 管理多个命名的日志实例
type Factory struct {
	mu            sync.RWMutex
	loggers       map[string]Log
	defaultLogger Log
}

// NewLogFactory 创建日志工厂，传入默认日志实例
func NewLogFactory(defaultLogger Log) *Factory {
	return &Factory{
		loggers:       make(map[string]Log),
		defaultLogger: defaultLogger,
	}
}

// Register 注册一个命名日志实例，重复注册会覆盖
func (f *Factory) Register(name string, logger Log) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.loggers[name] = logger
}

// Get 获取命名日志实例，如果不存在返回默认日志
func (f *Factory) Get(name string) Log {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if logger, ok := f.loggers[name]; ok {
		return logger
	}
	return f.defaultLogger
}
