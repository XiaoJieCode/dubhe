package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// 彩色终端支持判断
// var isColorTerminal = term.IsTerminal(int(os.Stdout.Fd()))
var isColorTerminal = true

// ANSI 颜色定义
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
	colorCyan   = "\033[36m"
)

// KV 字段实现
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

// 默认实现
type defaultLogger struct {
	mu     sync.Mutex
	prefix string
	fields []Field
}

func NewDefaultLogger(prefix string) Log {
	return &defaultLogger{prefix: prefix}
}

func (l *defaultLogger) log(level string, msg string, fields ...Field) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ts := time.Now().Format("2006-01-02 15:04:05")
	tsStr := ts
	if isColorTerminal {
		tsStr = colorGray + ts + colorReset
	}

	// 彩色 prefix
	prefixStr := l.prefix
	if isColorTerminal {
		prefixStr = colorCyan + l.prefix + colorReset
	}

	// 彩色级别标签
	levelStr := level
	if isColorTerminal {
		switch level {
		case "DEBUG":
			levelStr = colorGray + level + colorReset
		case "INFO":
			levelStr = colorGreen + level + colorReset
		case "WARN":
			levelStr = colorYellow + level + colorReset
		case "ERROR":
			levelStr = colorRed + level + colorReset
		}
	}

	// 彩色消息内容（可自定义增强）
	msgStr := msg
	if isColorTerminal {
		msgStr = colorBlue + msg + colorReset
	}

	// 打印主内容
	fmt.Fprintf(os.Stdout, "%s [%s] %s: %s", tsStr, prefixStr, levelStr, msgStr)

	// 附加字段
	allFields := append(l.fields, fields...)
	if len(allFields) > 0 {
		fmt.Fprint(os.Stdout, " |")
		for _, f := range allFields {
			key := f.Key()
			val := f.Value()
			if isColorTerminal {
				fmt.Fprintf(os.Stdout, " %s=%v", colorCyan+key+colorReset, val)
			} else {
				fmt.Fprintf(os.Stdout, " %s=%v", key, val)
			}
		}
	}
	fmt.Fprintln(os.Stdout)
}

// 各级别方法
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

// 格式化方法
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

// 附加字段
func (l *defaultLogger) WithFields(fields ...Field) Log {
	newFields := append([]Field{}, l.fields...)
	newFields = append(newFields, fields...)
	return &defaultLogger{
		prefix: l.prefix,
		fields: newFields,
	}
}

// 命名子 logger
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
