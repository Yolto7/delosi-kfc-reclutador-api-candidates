package logger

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
)

type Logger interface {
	Debug(input any)
	Info(input any)
	Warn(input any)
	Error(input any)
}