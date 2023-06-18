package log

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func catchStd(std **os.File, f func()) (string, error) {
	old := *std

	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	*std = w

	outCh := make(chan string)
	go func() {
		defer r.Close()
		buf := new(bytes.Buffer)
		io.Copy(buf, r)
		outCh <- buf.String()
	}()

	f()
	w.Close()
	*std = old
	return <-outCh, nil
}

func TestConsoleLogger_Debug(t *testing.T) {
	logger := NewConsoleLogger()

	msg := "debug message"
	got, err := catchStd(&os.Stdout, func() {
		logger.Debug(msg)
	})

	if err != nil {
		t.Errorf("ConsoleLogger.Debug() error = %v", err)
	}

	if wantPrefix := "DEBUG\t"; !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("ConsoleLogger.Debug() = %v, wantPrefix %v", got, wantPrefix)
	}

	if wantSuffix := msg + "\n"; !strings.HasSuffix(got, wantSuffix) {
		t.Errorf("ConsoleLogger.Debug() = %v, wantSuffix %v", got, wantSuffix)
	}
}

func TestConsoleLogger_Info(t *testing.T) {
	logger := NewConsoleLogger()

	msg := "info message"
	got, err := catchStd(&os.Stdout, func() {
		logger.Info(msg)
	})

	if err != nil {
		t.Errorf("ConsoleLogger.Info() error = %v", err)
	}

	if wantPrefix := "INFO\t"; !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("ConsoleLogger.Info() = %v, wantPrefix %v", got, wantPrefix)
	}

	if wantSuffix := msg + "\n"; !strings.HasSuffix(got, wantSuffix) {
		t.Errorf("ConsoleLogger.Info() = %v, wantSuffix %v", got, wantSuffix)
	}
}

func TestConsoleLogger_Error(t *testing.T) {
	logger := NewConsoleLogger()

	msg := "error message"
	got, err := catchStd(&os.Stderr, func() {
		logger.Error(msg)
	})

	if err != nil {
		t.Errorf("ConsoleLogger.Error() error = %v", err)
	}

	if wantPrefix := "ERROR\t"; !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("ConsoleLogger.Error() = %v, wantPrefix %v", got, wantPrefix)
	}

	if wantSuffix := msg + "\n"; !strings.HasSuffix(got, wantSuffix) {
		t.Errorf("ConsoleLogger.Error() = %v, wantSuffix %v", got, wantSuffix)
	}
}

func TestNewConsoleLogger(t *testing.T) {
	NewConsoleLogger()
}
