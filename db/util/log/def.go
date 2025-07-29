package log

// Field 抽象日志字段，底层实现映射到具体日志框架
type Field interface {
	Key() string
	Value() any
}

// Log 日志接口，支持结构化和格式化日志
type Log interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)

	// WithFields 返回带附加字段的新 Logger（链式调用）
	WithFields(fields ...Field) Log

	// Named 返回命名的 Logger，方便日志分组
	Named(name string) Log
}
