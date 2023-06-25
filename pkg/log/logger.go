// Package log provides structures for logging.
package log

// Logger is interface implemented by types
// that can log at various levels.
type Logger interface {
	// Debug outputs variables in specified format with DEBUG prefix.
	Debug(format string, v ...any)

	// Info outputs variables in specified format with INFO prefix.
	Info(format string, v ...any)

	// Error outputs variables in specified format with ERROR prefix.
	Error(format string, v ...any)

	// Fatal outputs variables in specified format with FATAL prefix.
	// It also calls os.Exit(1).
	Fatal(format string, v ...any)
}
