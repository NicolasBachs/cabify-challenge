package logger

type Logger interface {
	Debug(service string, msg string, args ...interface{})
	Info(service string, msg string, args ...interface{})
	Warn(service string, msg string, args ...interface{})
	Error(service string, msg string, args ...interface{})
}
