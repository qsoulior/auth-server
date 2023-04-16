package logger

import (
	"log"
	"os"
)

type Logger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (logger *Logger) Info(v ...any) {
	logger.InfoLog.Print(v...)
}

func (logger *Logger) Error(v ...any) {
	logger.ErrorLog.Print(v...)
}

func New() *Logger {
	logger := new(Logger)
	logger.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	logger.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
