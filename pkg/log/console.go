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

func (logger *ConsoleLogger) Info(v ...any) {
	logger.InfoLog.Print(v...)
}

func (logger *ConsoleLogger) Error(v ...any) {
	logger.ErrorLog.Print(v...)
}

func (logger *ConsoleLogger) Fatal(v ...any) {
	logger.FatalLog.Fatal(v...)
}

func NewConsoleLogger() Logger {
	logger := new(ConsoleLogger)
	logger.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	logger.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	logger.FatalLog = log.New(os.Stderr, "FATAL\t", log.Ldate|log.Ltime)
	return logger
}
