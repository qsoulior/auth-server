package log

type Logger interface {
	Info(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}
