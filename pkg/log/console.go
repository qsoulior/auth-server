package log

import (
	"log"
	"os"
)

type ConsoleLogger struct {
	*log.Logger
}

func (logger *ConsoleLogger) Debug(format string, v ...any) {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("DEBUG\t")
	logger.Printf(format, v...)
}

func (logger *ConsoleLogger) Info(format string, v ...any) {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("INFO\t")
	logger.Printf(format, v...)
}

func (logger *ConsoleLogger) Error(format string, v ...any) {
	logger.SetOutput(os.Stderr)
	logger.SetPrefix("ERROR\t")
	logger.Printf(format, v...)
}

func (logger *ConsoleLogger) Fatal(format string, v ...any) {
	logger.SetOutput(os.Stderr)
	logger.SetPrefix("FATAL\t")
	logger.Fatalf(format, v...)
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}
