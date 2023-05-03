package log

import (
	"log"
	"os"
)

type ConsoleLogger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	FatalLog *log.Logger
}

func (logger *ConsoleLogger) Info(format string, v ...any) {
	logger.InfoLog.Printf(format, v...)
}

func (logger *ConsoleLogger) Error(format string, v ...any) {
	logger.ErrorLog.Printf(format, v...)
}

func (logger *ConsoleLogger) Fatal(format string, v ...any) {
	logger.FatalLog.Fatalf(format, v...)
}

func NewConsoleLogger() *ConsoleLogger {
	logger := new(ConsoleLogger)
	logger.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	logger.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	logger.FatalLog = log.New(os.Stderr, "FATAL\t", log.Ldate|log.Ltime)
	return logger
}
