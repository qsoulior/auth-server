package log

type Logger interface {
	Info(format string, v ...any)
	Error(format string, v ...any)
	Fatal(format string, v ...any)
}
