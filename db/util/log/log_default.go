package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type KVField struct {
	key string
	val any
}

func (f *KVField) Key() string {
	return f.key
}
func (f *KVField) Value() any {
	return f.val
}
func KV(key string, val any) *KVField {
	return &KVField{key: key, val: val}
}

type defaultLogger struct {
	mu     sync.Mutex
	prefix string  // logger 名字，如模块名
	fields []Field // 附加的上下文字段，WithFields支持
}

func NewDefaultLogger(prefix string) Log {
	return &defaultLogger{prefix: prefix}
}

func (l *defaultLogger) log(level string, msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ts := time.Now().Format("2006-01-02 15:04:05")

	// 打印前缀和级别
	fmt.Fprintf(os.Stdout, "%s [%s] %s: %s", ts, l.prefix, level, msg)

	// 打印上下文字段
	allFields := append(l.fields, fields...)
	if len(allFields) > 0 {
		fmt.Fprint(os.Stdout, " |")
		for _, f := range allFields {
			fmt.Fprintf(os.Stdout, " %s=%v", f.Key(), f.Value())
		}
	}
	fmt.Fprintln(os.Stdout)
}

func (l *defaultLogger) Debug(msg string, fields ...Field) {
	l.log("DEBUG", msg, fields...)
}

func (l *defaultLogger) Info(msg string, fields ...Field) {
	l.log("INFO", msg, fields...)
}

func (l *defaultLogger) Warn(msg string, fields ...Field) {
	l.log("WARN", msg, fields...)
}

func (l *defaultLogger) Error(msg string, fields ...Field) {
	l.log("ERROR", msg, fields...)
}

func (l *defaultLogger) Debugf(format string, args ...any) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Infof(format string, args ...any) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Warnf(format string, args ...any) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *defaultLogger) WithFields(fields ...Field) Log {
	// 复制当前 logger 并追加新字段，支持链式上下文字段
	newFields := append([]Field{}, l.fields...)
	newFields = append(newFields, fields...)

	return &defaultLogger{
		prefix: l.prefix,
		fields: newFields,
	}
}

// Named 返回一个带新命名空间的 logger，支持链式调用
func (l *defaultLogger) Named(name string) Log {
	newPrefix := name
	if l.prefix != "" {
		newPrefix = l.prefix + "." + name
	}
	return &defaultLogger{
		prefix: newPrefix,
		fields: l.fields,
	}
}
