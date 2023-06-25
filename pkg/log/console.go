package log

import (
	"log"
	"os"
)

// ConsoleLogger implements Logger interface.
// It represents logger that outputs to os.Stdout and os.Stderr.
type ConsoleLogger struct {
	*log.Logger
}

// Debug outputs variables in specified format with DEBUG prefix.
func (logger *ConsoleLogger) Debug(format string, v ...any) {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("DEBUG\t")
	logger.Printf(format, v...)
}

// Info outputs variables in specified format with INFO prefix.
func (logger *ConsoleLogger) Info(format string, v ...any) {
	logger.SetOutput(os.Stdout)
	logger.SetPrefix("INFO\t")
	logger.Printf(format, v...)
}

// Error outputs variables in specified format with ERROR prefix.
func (logger *ConsoleLogger) Error(format string, v ...any) {
	logger.SetOutput(os.Stderr)
	logger.SetPrefix("ERROR\t")
	logger.Printf(format, v...)
}

// Fatal outputs variables in specified format with FATAL prefix.
// It also calls os.Exit(1).
func (logger *ConsoleLogger) Fatal(format string, v ...any) {
	logger.SetOutput(os.Stderr)
	logger.SetPrefix("FATAL\t")
	logger.Fatalf(format, v...)
}

// NewConsoleLogger creates *log.Logger and embeds it into ConsoleLogger.
// It returns pointer to a ConsoleLogger instance.
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{
		log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}
